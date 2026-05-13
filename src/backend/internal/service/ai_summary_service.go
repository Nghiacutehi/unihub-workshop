package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"regexp"
	"strings"
	"time"

	"unihub-workshop/internal/circuitbreaker"
	"unihub-workshop/internal/repository"
)

// AISummaryService implements Pipe-and-Filter architecture for PDF summarization
type AISummaryService struct {
	workshopRepo *repository.WorkshopRepo
	breaker      *circuitbreaker.CircuitBreaker
	apiURL       string
	apiKey       string
	httpClient   *http.Client
}

func NewAISummaryService(workshopRepo *repository.WorkshopRepo, apiURL, apiKey string) *AISummaryService {
	return &AISummaryService{
		workshopRepo: workshopRepo,
		breaker:      circuitbreaker.NewCircuitBreaker("ai-service", 0.5, 30*time.Second, 60*time.Second),
		apiURL:       apiURL,
		apiKey:       apiKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// ProcessPDF implements the Pipe-and-Filter pipeline
func (s *AISummaryService) ProcessPDF(ctx context.Context, workshopID string, pdfContent []byte) (string, error) {
	log.Printf("[AI_SUMMARY] Starting pipeline for workshop %s", workshopID)

	// Filter 1: Extract text from PDF
	rawText, err := s.extractText(pdfContent)
	if err != nil {
		return "", fmt.Errorf("extraction failed: %w", err)
	}

	// Filter 2: Clean text
	cleanedText := s.cleanText(rawText)

	// Filter 3: Build prompt
	prompt := s.buildPrompt(cleanedText)

	// Filter 4: Call AI model
	var summary string
	err = s.breaker.Execute(func() error {
		var aiErr error
		summary, aiErr = s.callAI(ctx, prompt)
		return aiErr
	})
	if err != nil {
		if err == circuitbreaker.ErrCircuitOpen {
			return "", fmt.Errorf("AI service temporarily unavailable")
		}
		return "", fmt.Errorf("AI call failed: %w", err)
	}

	// Sink: Persist to database
	if err := s.workshopRepo.UpdateSummary(ctx, workshopID, summary); err != nil {
		return "", fmt.Errorf("failed to persist summary: %w", err)
	}

	log.Printf("[AI_SUMMARY] Pipeline completed for workshop %s", workshopID)
	return summary, nil
}

// Filter 1: Extract text from PDF bytes
func (s *AISummaryService) extractText(pdfContent []byte) (string, error) {
	// Simple text extraction - look for text between stream markers
	content := string(pdfContent)
	var textParts []string

	// Extract readable text using regex patterns
	re := regexp.MustCompile(`\(([^)]+)\)`)
	matches := re.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) > 1 {
			text := match[1]
			if len(text) > 2 && isPrintable(text) {
				textParts = append(textParts, text)
			}
		}
	}

	if len(textParts) == 0 {
		// Fallback: just use the raw bytes as text if it looks like text
		cleaned := strings.Map(func(r rune) rune {
			if r >= 32 && r < 127 || r == '\n' || r == '\r' || r == '\t' {
				return r
			}
			return -1
		}, content)
		if len(cleaned) > 100 {
			return cleaned, nil
		}
		return "", fmt.Errorf("no text content found in PDF")
	}

	return strings.Join(textParts, " "), nil
}

func isPrintable(s string) bool {
	for _, r := range s {
		if r < 32 || r > 126 {
			return false
		}
	}
	return true
}

// Filter 2: Clean and normalize text
func (s *AISummaryService) cleanText(text string) string {
	// Remove extra whitespace
	spaceRe := regexp.MustCompile(`\s+`)
	text = spaceRe.ReplaceAllString(text, " ")

	// Trim
	text = strings.TrimSpace(text)

	// Token limit (approximately 4000 tokens ~ 16000 chars)
	maxChars := 16000
	if len(text) > maxChars {
		text = text[:maxChars]
	}

	return text
}

// Filter 3: Build AI prompt
func (s *AISummaryService) buildPrompt(text string) string {
	return fmt.Sprintf(`Bạn là trợ lý AI chuyên tóm tắt nội dung workshop. 
Hãy tóm tắt nội dung sau trong 5 gạch đầu dòng bằng Tiếng Việt, tập trung vào:
- Chủ đề chính
- Mục tiêu học tập
- Đối tượng phù hợp
- Kỹ năng sẽ học được
- Điểm nổi bật

Nội dung:
%s`, text)
}

// Filter 4: Call AI API (OpenAI-compatible)
func (s *AISummaryService) callAI(ctx context.Context, prompt string) (string, error) {
	if s.apiKey == "" {
		// Fallback: generate a mock summary when no API key
		return "• Workshop này giới thiệu các khái niệm quan trọng trong lĩnh vực liên quan\n• Phù hợp cho sinh viên muốn mở rộng kiến thức chuyên ngành\n• Người tham dự sẽ được thực hành trực tiếp\n• Diễn giả là chuyên gia có kinh nghiệm trong ngành\n• Cơ hội networking với các bạn cùng chí hướng", nil
	}

	reqBody := map[string]interface{}{
		"model": "gpt-3.5-turbo",
		"messages": []map[string]string{
			{"role": "system", "content": "Bạn là trợ lý AI chuyên tóm tắt nội dung workshop đại học."},
			{"role": "user", "content": prompt},
		},
		"temperature": 0.3,
		"max_tokens":  500,
	}

	jsonBody, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, "POST", s.apiURL, bytes.NewReader(jsonBody))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.apiKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != 200 {
		return "", fmt.Errorf("AI API error %d: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", err
	}
	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from AI")
	}

	return result.Choices[0].Message.Content, nil
}
