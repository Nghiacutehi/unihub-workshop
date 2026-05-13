package service

import (
	"context"
	"fmt"
	"log"
	"net/smtp"

	"unihub-workshop/internal/model"
	"unihub-workshop/internal/repository"
)

// NotificationStrategy interface (Strategy Pattern)
type NotificationStrategy interface {
	Send(ctx context.Context, notif *model.Notification) error
	Channel() model.NotificationChannel
}

// NotificationService dispatches notifications using multiple strategies (Observer + Strategy Pattern)
type NotificationService struct {
	repo       *repository.NotificationRepo
	strategies []NotificationStrategy
}

func NewNotificationService(repo *repository.NotificationRepo, strategies ...NotificationStrategy) *NotificationService {
	return &NotificationService{repo: repo, strategies: strategies}
}

// Dispatch sends notification through all configured channels
func (s *NotificationService) Dispatch(ctx context.Context, event model.NotificationEvent) {
	for _, strategy := range strategies(s.strategies) {
		notif := &model.Notification{
			UserID:         event.UserID,
			RegistrationID: &event.RegistrationID,
			Channel:        strategy.Channel(),
			Title:          buildTitle(event.Type, event.WorkshopTitle),
			Content:        buildContent(event),
			Status:         model.NotifPending,
			EventID:        event.EventID,
		}

		// Idempotent insert (ON CONFLICT DO NOTHING)
		if err := s.repo.Create(ctx, notif); err != nil {
			log.Printf("[NOTIFICATION] Duplicate event %s for channel %s, skipping", event.EventID, strategy.Channel())
			continue
		}

		// Send through strategy
		go func(strat NotificationStrategy, n *model.Notification) {
			if err := strat.Send(ctx, n); err != nil {
				errMsg := err.Error()
				_ = s.repo.UpdateStatus(ctx, n.ID, model.NotifFailed, &errMsg)
				log.Printf("[NOTIFICATION] Failed to send via %s: %v", strat.Channel(), err)
			} else {
				_ = s.repo.UpdateStatus(ctx, n.ID, model.NotifSent, nil)
				log.Printf("[NOTIFICATION] Sent via %s for event %s", strat.Channel(), n.EventID)
			}
		}(strategy, notif)
	}
}

func (s *NotificationService) GetUserNotifications(ctx context.Context, userID string) ([]model.Notification, error) {
	return s.repo.FindByUser(ctx, userID)
}

func strategies(s []NotificationStrategy) []NotificationStrategy { return s }

func buildTitle(eventType, workshopTitle string) string {
	switch eventType {
	case "REGISTRATION_SUCCESS":
		return fmt.Sprintf("Đăng ký thành công: %s", workshopTitle)
	case "PAYMENT_SUCCESS":
		return fmt.Sprintf("Thanh toán thành công: %s", workshopTitle)
	default:
		return fmt.Sprintf("Thông báo: %s", workshopTitle)
	}
}

func buildContent(event model.NotificationEvent) string {
	switch event.Type {
	case "REGISTRATION_SUCCESS":
		return fmt.Sprintf("Bạn đã đăng ký thành công workshop \"%s\". Mã QR check-in đã được tạo.", event.WorkshopTitle)
	case "PAYMENT_SUCCESS":
		return fmt.Sprintf("Thanh toán cho workshop \"%s\" đã được xác nhận. Mã QR check-in đã được tạo.", event.WorkshopTitle)
	default:
		return fmt.Sprintf("Cập nhật về workshop \"%s\".", event.WorkshopTitle)
	}
}

// ==========================================
// EmailStrategy
// ==========================================

type EmailStrategy struct {
	host string
	port string
	from string
	userRepo *repository.UserRepo
}

func NewEmailStrategy(host, port, from string, userRepo *repository.UserRepo) *EmailStrategy {
	return &EmailStrategy{host: host, port: port, from: from, userRepo: userRepo}
}

func (e *EmailStrategy) Channel() model.NotificationChannel {
	return model.ChannelEmail
}

func (e *EmailStrategy) Send(ctx context.Context, notif *model.Notification) error {
	user, err := e.userRepo.FindByID(ctx, notif.UserID)
	if err != nil {
		return fmt.Errorf("user not found: %w", err)
	}

	msg := fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n<h2>%s</h2><p>%s</p>",
		e.from, user.Email, notif.Title, notif.Title, notif.Content)

	addr := fmt.Sprintf("%s:%s", e.host, e.port)
	return smtp.SendMail(addr, nil, e.from, []string{user.Email}, []byte(msg))
}

// ==========================================
// WebNotificationStrategy
// ==========================================

type WebNotificationStrategy struct{}

func NewWebNotificationStrategy() *WebNotificationStrategy {
	return &WebNotificationStrategy{}
}

func (w *WebNotificationStrategy) Channel() model.NotificationChannel {
	return model.ChannelWeb
}

func (w *WebNotificationStrategy) Send(ctx context.Context, notif *model.Notification) error {
	// Web notifications are stored in DB and fetched by frontend polling
	// The Create call already persisted it, so this is a no-op success
	log.Printf("[WEB_NOTIF] Stored web notification for user %s", notif.UserID)
	return nil
}
