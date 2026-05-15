package service

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
	"unihub-workshop/internal/model"
	"unihub-workshop/internal/repository"
)

type BatchImportService struct {
	importRepo *repository.ImportRepo
	userRepo   *repository.UserRepo
	importDir  string
	archiveDir string
}

func NewBatchImportService(importRepo *repository.ImportRepo, userRepo *repository.UserRepo, importDir, archiveDir string) *BatchImportService {
	// Create directories if they don't exist
	os.MkdirAll(importDir, 0755)
	os.MkdirAll(archiveDir, 0755)

	return &BatchImportService{
		importRepo: importRepo,
		userRepo:   userRepo,
		importDir:  importDir,
		archiveDir: archiveDir,
	}
}

// ProcessCSV implements the Batch Sequential pipeline: Extract → Transform → Load → Cleanup
func (s *BatchImportService) ProcessCSV(ctx context.Context, filePath string) (*model.ImportJob, error) {
	fileName := filepath.Base(filePath)
	log.Printf("[BATCH_IMPORT] Starting import: %s", fileName)

	job := &model.ImportJob{
		FileName: fileName,
		Status:   model.ImportProcessing,
	}
	if err := s.importRepo.CreateJob(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to create import job: %w", err)
	}

	// Phase 1: Extract
	records, err := s.extract(filePath)
	if err != nil {
		job.Status = model.ImportFailed
		s.importRepo.UpdateJob(ctx, job)
		return job, fmt.Errorf("extract failed: %w", err)
	}

	job.TotalRecords = len(records)

	// Phase 2: Transform
	validRecords, importErrors := s.transform(records, job.ID)

	// Phase 3: Load (chunk-oriented, 1000 records per chunk)
	chunkSize := 1000
	successCount := 0
	for i := 0; i < len(validRecords); i += chunkSize {
		end := i + chunkSize
		if end > len(validRecords) {
			end = len(validRecords)
		}
		chunk := validRecords[i:end]

		for _, rec := range chunk {
			if err := s.userRepo.UpsertFromCSV(ctx, rec[0], rec[1], rec[2], rec[3], rec[4], rec[5]); err != nil {
				log.Printf("[BATCH_IMPORT] Failed to upsert student %s: %v", rec[0], err)
				importErrors = append(importErrors, model.ImportError{
					JobID:       job.ID,
					RowNumber:   i,
					RawData:     strings.Join(rec, ","),
					ErrorReason: err.Error(),
				})
			} else {
				successCount++
			}
		}
	}

	// Log errors to DB
	for _, e := range importErrors {
		_ = s.importRepo.CreateError(ctx, &e)
	}

	// Phase 4: Cleanup
	s.cleanup(filePath, fileName)

	job.SuccessCount = successCount
	job.ErrorCount = len(importErrors)
	job.Status = model.ImportCompleted
	s.importRepo.UpdateJob(ctx, job)

	log.Printf("[BATCH_IMPORT] Completed: %s | total=%d success=%d errors=%d",
		fileName, job.TotalRecords, job.SuccessCount, job.ErrorCount)

	return job, nil
}

func (s *BatchImportService) extract(filePath string) ([][]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("cannot open file: %w", err)
	}
	defer file.Close()

	reader := csv.NewReader(file)
	// Skip header
	if _, err := reader.Read(); err != nil {
		return nil, fmt.Errorf("cannot read header: %w", err)
	}

	var records [][]string
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("csv read error: %w", err)
		}
		records = append(records, record)
	}
	return records, nil
}

func (s *BatchImportService) transform(records [][]string, jobID string) ([][]string, []model.ImportError) {
	var valid [][]string
	var errors []model.ImportError
	seen := make(map[string]bool)

	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	studentIDRegex := regexp.MustCompile(`^\d{8,12}$`)

	for i, record := range records {
		rowNum := i + 2 // +2 for header and 0-index

		// Expect: student_id, password, full_name, email, phone, role
		if len(record) < 5 {
			errors = append(errors, model.ImportError{
				JobID:       jobID,
				RowNumber:   rowNum,
				RawData:     strings.Join(record, ","),
				ErrorReason: "insufficient columns",
			})
			continue
		}

		studentID := strings.TrimSpace(record[0])
		password := strings.TrimSpace(record[1])
		fullName := strings.TrimSpace(record[2])
		email := strings.TrimSpace(record[3])
		phone := strings.TrimSpace(record[4])
		role := "STUDENT"
		if len(record) >= 6 && strings.TrimSpace(record[5]) != "" {
			role = strings.ToUpper(strings.TrimSpace(record[5]))
		}

		// Validate student_id
		if !studentIDRegex.MatchString(studentID) {
			errors = append(errors, model.ImportError{
				JobID:       jobID,
				RowNumber:   rowNum,
				RawData:     strings.Join(record, ","),
				ErrorReason: fmt.Sprintf("invalid student_id format: %s", studentID),
			})
			continue
		}

		// Validate email
		if !emailRegex.MatchString(email) {
			errors = append(errors, model.ImportError{
				JobID:       jobID,
				RowNumber:   rowNum,
				RawData:     strings.Join(record, ","),
				ErrorReason: fmt.Sprintf("invalid email format: %s", email),
			})
			continue
		}

		// Deduplicate
		if seen[studentID] {
			errors = append(errors, model.ImportError{
				JobID:       jobID,
				RowNumber:   rowNum,
				RawData:     strings.Join(record, ","),
				ErrorReason: fmt.Sprintf("duplicate student_id in file: %s", studentID),
			})
			continue
		}
		seen[studentID] = true

		// Hash password
		hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			errors = append(errors, model.ImportError{
				JobID:       jobID,
				RowNumber:   rowNum,
				RawData:     strings.Join(record, ","),
				ErrorReason: "password hash failed",
			})
			continue
		}

		valid = append(valid, []string{studentID, string(hashedPw), fullName, email, phone, role})
	}

	return valid, errors
}

func (s *BatchImportService) cleanup(filePath, fileName string) {
	archivePath := filepath.Join(s.archiveDir, fmt.Sprintf("%s_%s", time.Now().Format("20060102_150405"), fileName))
	if err := os.Rename(filePath, archivePath); err != nil {
		log.Printf("[BATCH_IMPORT] Failed to archive file: %v", err)
	}
}

// ScanAndImport scans the import directory for CSV files and processes them
func (s *BatchImportService) ScanAndImport(ctx context.Context) {
	entries, err := os.ReadDir(s.importDir)
	if err != nil {
		log.Printf("[BATCH_IMPORT] Failed to scan directory: %v", err)
		return
	}

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".csv") {
			continue
		}
		filePath := filepath.Join(s.importDir, entry.Name())
		if _, err := s.ProcessCSV(ctx, filePath); err != nil {
			log.Printf("[BATCH_IMPORT] Failed to process %s: %v", entry.Name(), err)
		}
	}
}

func (s *BatchImportService) GetJobs(ctx context.Context) ([]model.ImportJob, error) {
	return s.importRepo.FindAllJobs(ctx)
}
