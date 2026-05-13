package model

import (
	"time"
)

// ==========================================
// User
// ==========================================

type Role string

const (
	RoleStudent   Role = "STUDENT"
	RoleStaff     Role = "STAFF"
	RoleOrganizer Role = "ORGANIZER"
)

type User struct {
	ID           string    `json:"id"`
	StudentID    string    `json:"student_id"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"full_name"`
	Email        string    `json:"email"`
	Phone        string    `json:"phone,omitempty"`
	Role         Role      `json:"role"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// ==========================================
// Workshop
// ==========================================

type WorkshopStatus string

const (
	WorkshopPublished WorkshopStatus = "PUBLISHED"
	WorkshopCancelled WorkshopStatus = "CANCELLED"
	WorkshopDraft     WorkshopStatus = "DRAFT"
)

type Workshop struct {
	ID             string         `json:"id"`
	Title          string         `json:"title"`
	Description    string         `json:"description,omitempty"`
	Speaker        string         `json:"speaker"`
	Room           string         `json:"room"`
	StartTime      time.Time      `json:"start_time"`
	EndTime        time.Time      `json:"end_time"`
	Capacity       int            `json:"capacity"`
	AvailableSeats int            `json:"available_seats"`
	Price          float64        `json:"price"`
	Summary        *string        `json:"summary,omitempty"`
	Status         WorkshopStatus `json:"status"`
	CreatedAt      time.Time      `json:"created_at"`
}

// ==========================================
// Registration
// ==========================================

type RegistrationStatus string

const (
	RegPendingPayment RegistrationStatus = "PENDING_PAYMENT"
	RegSuccess        RegistrationStatus = "SUCCESS"
	RegFailed         RegistrationStatus = "FAILED"
	RegCancelled      RegistrationStatus = "CANCELLED"
	RegRejected       RegistrationStatus = "REJECTED"
)

type Registration struct {
	ID          string             `json:"id"`
	UserID      string             `json:"user_id"`
	WorkshopID  string             `json:"workshop_id"`
	Status      RegistrationStatus `json:"status"`
	QRCode      *string            `json:"qr_code,omitempty"`
	IsCheckedIn bool               `json:"is_checked_in"`
	ScannedAt   *time.Time         `json:"scanned_at,omitempty"`
	CreatedAt   time.Time          `json:"created_at"`
}

// ==========================================
// Payment
// ==========================================

type PaymentStatus string

const (
	PaymentPending   PaymentStatus = "PENDING"
	PaymentSuccess   PaymentStatus = "SUCCESS"
	PaymentFailed    PaymentStatus = "FAILED"
	PaymentCancelled PaymentStatus = "CANCELLED"
)

type Payment struct {
	ID             string        `json:"id"`
	RegistrationID string        `json:"registration_id"`
	TransactionID  string        `json:"transaction_id"`
	Amount         float64       `json:"amount"`
	Provider       string        `json:"provider"`
	Status         PaymentStatus `json:"status"`
	CreatedAt      time.Time     `json:"created_at"`
}

// ==========================================
// Notification
// ==========================================

type NotificationChannel string

const (
	ChannelEmail NotificationChannel = "EMAIL"
	ChannelWeb   NotificationChannel = "WEB"
)

type NotificationStatus string

const (
	NotifPending NotificationStatus = "PENDING"
	NotifSent    NotificationStatus = "SENT"
	NotifFailed  NotificationStatus = "FAILED"
)

type Notification struct {
	ID             string              `json:"id"`
	UserID         string              `json:"user_id"`
	RegistrationID *string             `json:"registration_id,omitempty"`
	Channel        NotificationChannel `json:"channel"`
	Title          string              `json:"title"`
	Content        string              `json:"content"`
	Status         NotificationStatus  `json:"status"`
	EventID        string              `json:"event_id"`
	ErrorMessage   *string             `json:"error_message,omitempty"`
	CreatedAt      time.Time           `json:"created_at"`
	SentAt         *time.Time          `json:"sent_at,omitempty"`
}

// ==========================================
// Import Job
// ==========================================

type ImportJobStatus string

const (
	ImportProcessing ImportJobStatus = "PROCESSING"
	ImportCompleted  ImportJobStatus = "COMPLETED"
	ImportFailed     ImportJobStatus = "FAILED"
)

type ImportJob struct {
	ID           string          `json:"id"`
	FileName     string          `json:"file_name"`
	Status       ImportJobStatus `json:"status"`
	TotalRecords int             `json:"total_records"`
	SuccessCount int             `json:"success_count"`
	ErrorCount   int             `json:"error_count"`
	StartedAt    time.Time       `json:"started_at"`
	CompletedAt  *time.Time      `json:"completed_at,omitempty"`
}

type ImportError struct {
	ID          string    `json:"id"`
	JobID       string    `json:"job_id"`
	RowNumber   int       `json:"row_number"`
	RawData     string    `json:"raw_data"`
	ErrorReason string    `json:"error_reason"`
	CreatedAt   time.Time `json:"created_at"`
}

// ==========================================
// API Request/Response DTOs
// ==========================================

type LoginRequest struct {
	StudentID string `json:"student_id"`
	Password  string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type RegisterWorkshopRequest struct {
	WorkshopID string `json:"workshop_id"`
}

type RegistrationStatusResponse struct {
	CorrelationID string             `json:"correlation_id"`
	Status        RegistrationStatus `json:"status"`
	Registration  *Registration      `json:"registration,omitempty"`
	Message       string             `json:"message,omitempty"`
}

type CreateWorkshopRequest struct {
	Title       string  `json:"title"`
	Description string  `json:"description"`
	Speaker     string  `json:"speaker"`
	Room        string  `json:"room"`
	StartTime   string  `json:"start_time"`
	EndTime     string  `json:"end_time"`
	Capacity    int     `json:"capacity"`
	Price       float64 `json:"price"`
}

type UpdateWorkshopRequest struct {
	Title       *string  `json:"title,omitempty"`
	Description *string  `json:"description,omitempty"`
	Speaker     *string  `json:"speaker,omitempty"`
	Room        *string  `json:"room,omitempty"`
	StartTime   *string  `json:"start_time,omitempty"`
	EndTime     *string  `json:"end_time,omitempty"`
	Capacity    *int     `json:"capacity,omitempty"`
	Price       *float64 `json:"price,omitempty"`
	Status      *string  `json:"status,omitempty"`
}

type CheckinRequest struct {
	RegistrationID string `json:"registration_id"`
	WorkshopID     string `json:"workshop_id"`
	StudentID      string `json:"student_id"`
}

type BulkCheckinRequest struct {
	Records []OfflineCheckinRecord `json:"records"`
}

type OfflineCheckinRecord struct {
	ID         string `json:"id"`
	StudentID  string `json:"student_id"`
	WorkshopID string `json:"workshop_id"`
	ScannedAt  int64  `json:"scanned_at"`
}

type PaymentWebhookRequest struct {
	TransactionID string `json:"transaction_id"`
	Status        string `json:"status"`
	Signature     string `json:"signature"`
}

type QueueMessage struct {
	CorrelationID string `json:"correlation_id"`
	UserID        string `json:"user_id"`
	WorkshopID    string `json:"workshop_id"`
	Action        string `json:"action"`
}

type NotificationEvent struct {
	EventID        string `json:"event_id"`
	UserID         string `json:"user_id"`
	RegistrationID string `json:"registration_id"`
	Type           string `json:"type"`
	WorkshopTitle  string `json:"workshop_title"`
	QRCode         string `json:"qr_code,omitempty"`
}

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
	Error   string      `json:"error,omitempty"`
}
