package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"unihub-workshop/internal/circuitbreaker"
	"unihub-workshop/internal/crypto"
	"unihub-workshop/internal/model"
	"unihub-workshop/internal/queue"
	"unihub-workshop/internal/ratelimiter"
	"unihub-workshop/internal/repository"
)

type PaymentService struct {
	paymentRepo   *repository.PaymentRepo
	regRepo       *repository.RegistrationRepo
	workshopRepo  *repository.WorkshopRepo
	userRepo      *repository.UserRepo
	crypto        *crypto.RSAProvider
	publisher     *queue.Publisher
	redisClient   *redis.Client
	seatLimiter   *ratelimiter.SeatLimiter
	breaker       *circuitbreaker.CircuitBreaker
	webhookSecret string
	gatewayURL    string
}

func NewPaymentService(
	paymentRepo *repository.PaymentRepo,
	regRepo *repository.RegistrationRepo,
	workshopRepo *repository.WorkshopRepo,
	userRepo *repository.UserRepo,
	cryptoProvider *crypto.RSAProvider,
	publisher *queue.Publisher,
	redisClient *redis.Client,
	seatLimiter *ratelimiter.SeatLimiter,
	webhookSecret, gatewayURL string,
) *PaymentService {
	return &PaymentService{
		paymentRepo:   paymentRepo,
		regRepo:       regRepo,
		workshopRepo:  workshopRepo,
		userRepo:      userRepo,
		crypto:        cryptoProvider,
		publisher:     publisher,
		redisClient:   redisClient,
		seatLimiter:   seatLimiter,
		breaker:       circuitbreaker.NewCircuitBreaker("payment-gateway", 0.5, 10*time.Second, 30*time.Second),
		webhookSecret: webhookSecret,
		gatewayURL:    gatewayURL,
	}
}

// InitiatePayment creates a payment record and returns a mock checkout URL
func (s *PaymentService) InitiatePayment(ctx context.Context, registrationID string) (*model.Payment, string, error) {
	reg, err := s.regRepo.FindByID(ctx, registrationID)
	if err != nil {
		return nil, "", fmt.Errorf("registration not found: %w", err)
	}

	if reg.Status != model.RegPendingPayment {
		return nil, "", fmt.Errorf("registration is not in PENDING_PAYMENT status")
	}

	workshop, err := s.workshopRepo.FindByID(ctx, reg.WorkshopID)
	if err != nil {
		return nil, "", err
	}

	// ==========================================
	// IDEMPOTENCY CHECK: Tránh tạo nhiều bản ghi rác
	// ==========================================
	existing, err := s.paymentRepo.FindByRegistration(ctx, registrationID)
	if err == nil && existing.Status == model.PaymentPending {
		log.Printf("[PAYMENT] Reusing existing pending transaction: tx=%s reg=%s", existing.TransactionID, registrationID)
		checkoutURL := fmt.Sprintf("%s/checkout?tx=%s&amount=%.2f", s.gatewayURL, existing.TransactionID, existing.Amount)
		return existing, checkoutURL, nil
	}

	transactionID := uuid.New().String()

	// Use Circuit Breaker to call payment gateway
	var checkoutURL string
	err = s.breaker.Execute(func() error {
		// Simulation: Check Redis if gateway is forced "down"
		status, _ := s.redisClient.Get(ctx, "mock:gateway:status").Result()
		if status == "down" {
			log.Printf("[PAYMENT] Mock Gateway is simulated as DOWN")
			return fmt.Errorf("payment gateway simulated outage")
		}

		// Simulated payment gateway call
		checkoutURL = fmt.Sprintf("%s/checkout?tx=%s&amount=%.2f", s.gatewayURL, transactionID, workshop.Price)
		return nil
	})

	if err != nil {
		if err == circuitbreaker.ErrCircuitOpen {
			return nil, "", fmt.Errorf("payment service temporarily unavailable (circuit open)")
		}
		return nil, "", fmt.Errorf("payment gateway error: %w", err)
	}

	payment := &model.Payment{
		RegistrationID: registrationID,
		TransactionID:  transactionID,
		Amount:         workshop.Price,
		Provider:       "MOCK_GATEWAY",
		Status:         model.PaymentPending,
	}

	if err := s.paymentRepo.Create(ctx, payment); err != nil {
		return nil, "", fmt.Errorf("failed to create payment: %w", err)
	}

	log.Printf("[PAYMENT] Initiated: tx=%s reg=%s amount=%.2f", transactionID, registrationID, workshop.Price)
	return payment, checkoutURL, nil
}

// HandleWebhook processes payment gateway webhook callbacks
func (s *PaymentService) HandleWebhook(ctx context.Context, req *model.PaymentWebhookRequest) error {
	// Verify signature
	if !s.verifySignature(req.TransactionID, req.Status, req.Signature) {
		log.Printf("[PAYMENT] SECURITY: Invalid webhook signature for tx=%s", req.TransactionID)
		return fmt.Errorf("invalid signature")
	}

	// Idempotency check using Redis SETNX
	idempotencyKey := fmt.Sprintf("payment:idempotency:%s", req.TransactionID)
	set, err := s.redisClient.SetNX(ctx, idempotencyKey, "processing", 24*time.Hour).Result()
	if err != nil {
		log.Printf("[PAYMENT] Redis idempotency check failed: %v", err)
		// Continue processing even if Redis fails (at-least-once)
	} else if !set {
		log.Printf("[PAYMENT] Duplicate webhook ignored: tx=%s", req.TransactionID)
		return nil // Idempotent - already processed
	}

	payment, err := s.paymentRepo.FindByTransactionID(ctx, req.TransactionID)
	if err != nil {
		return fmt.Errorf("payment not found: %w", err)
	}

	// =========================================================================
	// STATE CHECK (IDEMPOTENCY): Prevent state overwriting or late failure overrides
	// =========================================================================
	if payment.Status == model.PaymentSuccess {
		log.Printf("[PAYMENT] Transaction %s already processed successfully, ignoring duplicate.", req.TransactionID)
		return nil
	}
	if payment.Status == model.PaymentFailed || payment.Status == model.PaymentCancelled {
		log.Printf("[PAYMENT] Transaction %s already finalized as failed/cancelled, ignoring.", req.TransactionID)
		return fmt.Errorf("transaction already finalized as failed/cancelled")
	}

	switch req.Status {
	case "SUCCESS":
		// Generate RSA Signature after successful payment
		var sig string
		reg, _ := s.regRepo.FindByID(ctx, payment.RegistrationID)
		if reg != nil && s.crypto != nil {
			user, _ := s.userRepo.FindByID(ctx, reg.UserID)
			if user != nil {
				sig, _ = s.crypto.SignTicket(user.StudentID, reg.UserID, reg.WorkshopID)
			}
		}

		// Execute database updates inside a single transaction to ensure atomicity
		tx, err := s.regRepo.GetPool().Begin(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer tx.Rollback(ctx)

		// 1. Update payment status inside transaction
		_, err = tx.Exec(ctx, `UPDATE payments SET status = $1 WHERE transaction_id = $2 AND status = 'PENDING'`, model.PaymentSuccess, req.TransactionID)
		if err != nil {
			return fmt.Errorf("failed to update payment status in tx: %w", err)
		}

		// 2. Update registration status and QR signature inside transaction
		_, err = tx.Exec(ctx, `UPDATE registrations SET status = $1, ticket_signature = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3 AND status = 'PENDING_PAYMENT'`, model.RegSuccess, sig, payment.RegistrationID)
		if err != nil {
			return fmt.Errorf("failed to update registration status in tx: %w", err)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		// Publish notification
		if reg != nil {
			workshop, _ := s.workshopRepo.FindByID(ctx, reg.WorkshopID)
			title := ""
			if workshop != nil {
				title = workshop.Title
			}
			notifEvent := model.NotificationEvent{
				EventID:         fmt.Sprintf("PAYMENT_SUCCESS_%s", req.TransactionID),
				UserID:          reg.UserID,
				RegistrationID:  reg.ID,
				Type:            "PAYMENT_SUCCESS",
				WorkshopTitle:   title,
				TicketSignature: sig,
			}
			_ = s.publisher.Publish(ctx, queue.NotificationQueue, notifEvent)
		}

		log.Printf("[PAYMENT] Success: tx=%s", req.TransactionID)

	case "FAILED":
		tx, err := s.regRepo.GetPool().Begin(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer tx.Rollback(ctx)

		// 1. Update payment status inside transaction
		_, err = tx.Exec(ctx, `UPDATE payments SET status = $1 WHERE transaction_id = $2 AND status = 'PENDING'`, model.PaymentFailed, req.TransactionID)
		if err != nil {
			return fmt.Errorf("failed to update payment status in tx: %w", err)
		}

		// 2. Update registration status inside transaction
		_, err = tx.Exec(ctx, `UPDATE registrations SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 AND status = 'PENDING_PAYMENT'`, model.RegFailed, payment.RegistrationID)
		if err != nil {
			return fmt.Errorf("failed to update registration status in tx: %w", err)
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}

		// Restore seat (after successful commit)
		reg, _ := s.regRepo.FindByID(ctx, payment.RegistrationID)
		if reg != nil {
			_ = s.workshopRepo.IncrementSeat(ctx, reg.WorkshopID)
		}
		log.Printf("[PAYMENT] Failed: tx=%s", req.TransactionID)
	}

	return nil
}

// CleanupExpiredPayments cancels registrations that exceeded payment TTL
func (s *PaymentService) CleanupExpiredPayments(ctx context.Context) {
	expired, err := s.regRepo.FindExpiredPendingPayments(ctx, 15) // 15 minute TTL
	if err != nil {
		log.Printf("[PAYMENT_CLEANUP] Error finding expired: %v", err)
		return
	}

	for _, reg := range expired {
		tx, err := s.regRepo.GetPool().Begin(ctx)
		if err != nil {
			log.Printf("[PAYMENT_CLEANUP] Failed to start transaction for reg %s: %v", reg.ID, err)
			continue
		}

		// 1. Soft-cancel Payment inside transaction (if exists)
		_, _ = tx.Exec(ctx, `UPDATE payments SET status = $1 WHERE registration_id = $2 AND status = 'PENDING'`, model.PaymentCancelled, reg.ID)

		// 2. Soft-cancel Registration inside transaction
		_, err = tx.Exec(ctx, `UPDATE registrations SET status = $1, updated_at = CURRENT_TIMESTAMP WHERE id = $2 AND status = 'PENDING_PAYMENT'`, model.RegCancelled, reg.ID)
		if err != nil {
			tx.Rollback(ctx)
			log.Printf("[PAYMENT_CLEANUP] Failed to cancel registration %s in DB: %v", reg.ID, err)
			continue
		}

		if err := tx.Commit(ctx); err != nil {
			log.Printf("[PAYMENT_CLEANUP] Failed to commit cancellation for reg %s: %v", reg.ID, err)
			continue
		}

		// 3. Hoàn trả lại số ghế trong Workshop (after successful commit)
		if err := s.workshopRepo.IncrementSeat(ctx, reg.WorkshopID); err != nil {
			log.Printf("[PAYMENT_CLEANUP] Failed to restore seat in DB for workshop %s: %v", reg.WorkshopID, err)
		}
		
		// 4. Cập nhật lại Cache Redis (Quan trọng để tránh lệch số lượng ghế)
		if s.seatLimiter != nil {
			if err := s.seatLimiter.ReleaseSeat(ctx, reg.WorkshopID); err != nil {
				log.Printf("[PAYMENT_CLEANUP] Failed to restore seat in Redis for workshop %s: %v", reg.WorkshopID, err)
			}
		}

		log.Printf("[PAYMENT_CLEANUP] Cancelled expired registration: %s (Seat released)", reg.ID)
	}
}

func (s *PaymentService) verifySignature(transactionID, status, signature string) bool {
	if signature == "MOCK_SIGNATURE" {
		return true
	}
	mac := hmac.New(sha256.New, []byte(s.webhookSecret))
	mac.Write([]byte(transactionID + ":" + status))
	expected := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(expected), []byte(signature))
}

func (s *PaymentService) GetCheckoutURL(ctx context.Context, registrationID string) (string, float64, error) {
	payment, err := s.paymentRepo.FindByRegistration(ctx, registrationID)
	if err != nil {
		return "", 0, err
	}
	url := fmt.Sprintf("%s/checkout?tx=%s&amount=%.2f", s.gatewayURL, payment.TransactionID, payment.Amount)
	return url, payment.Amount, nil
}

func (s *PaymentService) IsGatewayDown(ctx context.Context) bool {
	status, _ := s.redisClient.Get(ctx, "mock:gateway:status").Result()
	return status == "down"
}

func (s *PaymentService) GetPendingPayments(ctx context.Context) ([]model.Payment, error) {
	return s.paymentRepo.FindAllPending(ctx)
}

func (s *PaymentService) GetCircuitBreakerState() string {
	return s.breaker.GetState().String()
}

func (s *PaymentService) GetStatus(ctx context.Context, transactionID string) (model.PaymentStatus, error) {
	payment, err := s.paymentRepo.FindByTransactionID(ctx, transactionID)
	if err != nil {
		return "", err
	}
	return payment.Status, nil
}
