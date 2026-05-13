package repository

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"unihub-workshop/internal/model"
)

type RegistrationRepo struct {
	pool *pgxpool.Pool
}

func NewRegistrationRepo(pool *pgxpool.Pool) *RegistrationRepo {
	return &RegistrationRepo{pool: pool}
}

func (r *RegistrationRepo) Create(ctx context.Context, tx pgx.Tx, reg *model.Registration) error {
	return tx.QueryRow(ctx,
		`INSERT INTO registrations (user_id, workshop_id, status, qr_code)
		 VALUES ($1, $2, $3, $4) RETURNING id, created_at`,
		reg.UserID, reg.WorkshopID, reg.Status, reg.QRCode,
	).Scan(&reg.ID, &reg.CreatedAt)
}

func (r *RegistrationRepo) CreateDirect(ctx context.Context, reg *model.Registration) error {
	return r.pool.QueryRow(ctx,
		`INSERT INTO registrations (user_id, workshop_id, status, qr_code)
		 VALUES ($1, $2, $3, $4) RETURNING id, created_at`,
		reg.UserID, reg.WorkshopID, reg.Status, reg.QRCode,
	).Scan(&reg.ID, &reg.CreatedAt)
}

func (r *RegistrationRepo) FindByID(ctx context.Context, id string) (*model.Registration, error) {
	var reg model.Registration
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, workshop_id, status, qr_code, is_checked_in, scanned_at, created_at
		 FROM registrations WHERE id = $1`, id,
	).Scan(&reg.ID, &reg.UserID, &reg.WorkshopID, &reg.Status, &reg.QRCode,
		&reg.IsCheckedIn, &reg.ScannedAt, &reg.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("registration not found: %w", err)
	}
	return &reg, nil
}

func (r *RegistrationRepo) FindByUserAndWorkshop(ctx context.Context, userID, workshopID string) (*model.Registration, error) {
	var reg model.Registration
	err := r.pool.QueryRow(ctx,
		`SELECT id, user_id, workshop_id, status, qr_code, is_checked_in, scanned_at, created_at
		 FROM registrations WHERE user_id = $1 AND workshop_id = $2`, userID, workshopID,
	).Scan(&reg.ID, &reg.UserID, &reg.WorkshopID, &reg.Status, &reg.QRCode,
		&reg.IsCheckedIn, &reg.ScannedAt, &reg.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &reg, nil
}

func (r *RegistrationRepo) FindByUser(ctx context.Context, userID string) ([]model.Registration, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, workshop_id, status, qr_code, is_checked_in, scanned_at, created_at
		 FROM registrations WHERE user_id = $1 ORDER BY created_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var regs []model.Registration
	for rows.Next() {
		var reg model.Registration
		if err := rows.Scan(&reg.ID, &reg.UserID, &reg.WorkshopID, &reg.Status, &reg.QRCode,
			&reg.IsCheckedIn, &reg.ScannedAt, &reg.CreatedAt); err != nil {
			return nil, err
		}
		regs = append(regs, reg)
	}
	return regs, nil
}

func (r *RegistrationRepo) UpdateStatus(ctx context.Context, id string, status model.RegistrationStatus) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE registrations SET status = $1 WHERE id = $2`, status, id)
	return err
}

func (r *RegistrationRepo) UpdateStatusAndQR(ctx context.Context, id string, status model.RegistrationStatus, qrCode string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE registrations SET status = $1, qr_code = $2 WHERE id = $3`, status, qrCode, id)
	return err
}

func (r *RegistrationRepo) CheckIn(ctx context.Context, id string) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE registrations SET is_checked_in = TRUE, scanned_at = CURRENT_TIMESTAMP WHERE id = $1`, id)
	return err
}

func (r *RegistrationRepo) CheckInWithTime(ctx context.Context, id string, scannedAt int64) error {
	_, err := r.pool.Exec(ctx,
		`UPDATE registrations SET is_checked_in = TRUE, scanned_at = to_timestamp($1) WHERE id = $2
		 AND (is_checked_in = FALSE OR scanned_at > to_timestamp($1))`,
		scannedAt, id)
	return err
}

func (r *RegistrationRepo) FindByStudentAndWorkshop(ctx context.Context, studentID, workshopID string) (*model.Registration, error) {
	var reg model.Registration
	err := r.pool.QueryRow(ctx,
		`SELECT r.id, r.user_id, r.workshop_id, r.status, r.qr_code, r.is_checked_in, r.scanned_at, r.created_at
		 FROM registrations r
		 JOIN users u ON r.user_id = u.id
		 WHERE u.student_id = $1 AND r.workshop_id = $2 AND r.status IN ('SUCCESS', 'PENDING_PAYMENT')`,
		studentID, workshopID,
	).Scan(&reg.ID, &reg.UserID, &reg.WorkshopID, &reg.Status, &reg.QRCode,
		&reg.IsCheckedIn, &reg.ScannedAt, &reg.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &reg, nil
}

// FindExpiredPendingPayments finds registrations that have been PENDING_PAYMENT for more than the given minutes
func (r *RegistrationRepo) FindExpiredPendingPayments(ctx context.Context, minutes int) ([]model.Registration, error) {
	rows, err := r.pool.Query(ctx,
		`SELECT id, user_id, workshop_id, status, qr_code, is_checked_in, scanned_at, created_at
		 FROM registrations
		 WHERE status = 'PENDING_PAYMENT'
		   AND created_at < NOW() - INTERVAL '1 minute' * $1`, minutes,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var regs []model.Registration
	for rows.Next() {
		var reg model.Registration
		if err := rows.Scan(&reg.ID, &reg.UserID, &reg.WorkshopID, &reg.Status, &reg.QRCode,
			&reg.IsCheckedIn, &reg.ScannedAt, &reg.CreatedAt); err != nil {
			return nil, err
		}
		regs = append(regs, reg)
	}
	return regs, nil
}

func (r *RegistrationRepo) GetPool() *pgxpool.Pool {
	return r.pool
}
