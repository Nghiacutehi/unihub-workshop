package service

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"sync"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/skip2/go-qrcode"
	"unihub-workshop/internal/model"
	"unihub-workshop/internal/queue"
	"unihub-workshop/internal/ratelimiter"
	"unihub-workshop/internal/repository"
)

type RegistrationService struct {
	regRepo      *repository.RegistrationRepo
	workshopRepo *repository.WorkshopRepo
	publisher    *queue.Publisher
	redis        *redis.Client
	waitingRoom  *ratelimiter.WaitingRoom
	mu           sync.RWMutex
	statuses     map[string]*model.RegistrationStatusResponse
}

func NewRegistrationService(
	regRepo *repository.RegistrationRepo,
	workshopRepo *repository.WorkshopRepo,
	publisher *queue.Publisher,
	redisClient *redis.Client,
	waitingRoom *ratelimiter.WaitingRoom,
) *RegistrationService {
	return &RegistrationService{
		regRepo:      regRepo,
		workshopRepo: workshopRepo,
		publisher:    publisher,
		redis:        redisClient,
		waitingRoom:  waitingRoom,
		statuses:     make(map[string]*model.RegistrationStatusResponse),
	}
}

// CheckWaitingRoom checks a user's status in the virtual waiting room
func (s *RegistrationService) CheckWaitingRoom(ctx context.Context, workshopID, userID string) (*ratelimiter.WaitingRoomResult, error) {
	return s.waitingRoom.Enter(ctx, workshopID, userID)
}

// EnqueueRegistration pushes registration request to RabbitMQ and returns a correlation ID
func (s *RegistrationService) EnqueueRegistration(ctx context.Context, userID, workshopID string) (string, error) {
	// Check if already registered
	existing, _ := s.regRepo.FindByUserAndWorkshop(ctx, userID, workshopID)
	if existing != nil && (existing.Status == model.RegSuccess || existing.Status == model.RegPendingPayment) {
		return "", fmt.Errorf("already registered for this workshop")
	}

	correlationID := uuid.New().String()

	msg := model.QueueMessage{
		CorrelationID: correlationID,
		UserID:        userID,
		WorkshopID:    workshopID,
		Action:        "REGISTER",
	}

	// Set initial status
	s.SetStatus(correlationID, &model.RegistrationStatusResponse{
		CorrelationID: correlationID,
		Status:        "PROCESSING",
		Message:       "Your registration is being processed",
	})

	if err := s.publisher.Publish(ctx, queue.RegistrationQueue, msg); err != nil {
		s.SetStatus(correlationID, &model.RegistrationStatusResponse{
			CorrelationID: correlationID,
			Status:        model.RegFailed,
			Message:       "System is temporarily unavailable",
		})
		return "", fmt.Errorf("failed to enqueue registration: %w", err)
	}

	log.Printf("[REGISTRATION] Enqueued: correlation=%s user=%s workshop=%s", correlationID, userID, workshopID)
	return correlationID, nil
}

// ProcessRegistration is called by the background worker to process a registration message
func (s *RegistrationService) ProcessRegistration(ctx context.Context, msg model.QueueMessage) error {
	log.Printf("[WORKER] Processing registration: correlation=%s", msg.CorrelationID)

	// Get workshop info to determine if it's free or paid
	workshop, err := s.workshopRepo.FindByID(ctx, msg.WorkshopID)
	if err != nil {
		s.SetStatus(msg.CorrelationID, &model.RegistrationStatusResponse{
			CorrelationID: msg.CorrelationID,
			Status:        model.RegFailed,
			Message:       "Workshop not found",
		})
		return err
	}

	// Begin transaction with Pessimistic Locking
	tx, err := s.workshopRepo.GetPool().Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	// SELECT FOR UPDATE - Pessimistic Lock
	remainingSeats, err := s.workshopRepo.DecrementSeatWithLock(ctx, tx, msg.WorkshopID)
	if err != nil {
		s.SetStatus(msg.CorrelationID, &model.RegistrationStatusResponse{
			CorrelationID: msg.CorrelationID,
			Status:        model.RegRejected,
			Message:       "No available seats",
		})
		return err
	}

	// Determine status based on price
	var regStatus model.RegistrationStatus
	if workshop.Price > 0 {
		regStatus = model.RegPendingPayment
	} else {
		regStatus = model.RegSuccess
	}

	// Generate QR code for free workshops
	var qrCodeStr *string
	if regStatus == model.RegSuccess {
		qr, err := generateQRCode(msg.UserID, msg.WorkshopID)
		if err != nil {
			log.Printf("[WORKER] QR generation failed: %v", err)
		} else {
			qrCodeStr = &qr
		}
	}

	reg := &model.Registration{
		UserID:     msg.UserID,
		WorkshopID: msg.WorkshopID,
		Status:     regStatus,
		QRCode:     qrCodeStr,
	}

	if err := s.regRepo.Create(ctx, tx, reg); err != nil {
		return fmt.Errorf("failed to create registration: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	log.Printf("[WORKER] Registration created: id=%s seats_remaining=%d status=%s", reg.ID, remainingSeats, regStatus)

	s.SetStatus(msg.CorrelationID, &model.RegistrationStatusResponse{
		CorrelationID: msg.CorrelationID,
		Status:        regStatus,
		Registration:  reg,
		Message:       fmt.Sprintf("Registration %s", regStatus),
	})

	// If free workshop, publish notification event
	if regStatus == model.RegSuccess {
		notifEvent := model.NotificationEvent{
			EventID:        fmt.Sprintf("REG_SUCCESS_%s", reg.ID),
			UserID:         msg.UserID,
			RegistrationID: reg.ID,
			Type:           "REGISTRATION_SUCCESS",
			WorkshopTitle:  workshop.Title,
		}
		if qrCodeStr != nil {
			notifEvent.QRCode = *qrCodeStr
		}
		_ = s.publisher.Publish(ctx, queue.NotificationQueue, notifEvent)
	}

	return nil
}

func (s *RegistrationService) GetStatus(correlationID string) *model.RegistrationStatusResponse {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.statuses[correlationID]
}

func (s *RegistrationService) SetStatus(correlationID string, status *model.RegistrationStatusResponse) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.statuses[correlationID] = status
}

func (s *RegistrationService) GetUserRegistrations(ctx context.Context, userID string) ([]model.Registration, error) {
	return s.regRepo.FindByUser(ctx, userID)
}

func generateQRCode(userID, workshopID string) (string, error) {
	data := fmt.Sprintf(`{"user_id":"%s","workshop_id":"%s"}`, userID, workshopID)
	png, err := qrcode.Encode(data, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(png), nil
}
