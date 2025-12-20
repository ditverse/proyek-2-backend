package services

import (
	"backend-sarpras/internal/config"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// WhatsappService handles WhatsApp messaging via Fonnte API
type WhatsappService struct {
	Token string
	URL   string
}

// NewWhatsappService creates a new WhatsappService from config
func NewWhatsappService(cfg *config.Config) *WhatsappService {
	return &WhatsappService{
		Token: cfg.FonnteToken,
		URL:   cfg.FonnteURL,
	}
}

// SendMessage sends a WhatsApp message asynchronously using goroutine
// This method does not block the main thread
func (s *WhatsappService) SendMessage(target, message string) {
	go func() {
		if err := s.sendMessageSync(target, message); err != nil {
			log.Printf("❌ WhatsApp failed to send to %s: %v", target, err)
		} else {
			log.Printf("✅ WhatsApp sent successfully to %s", target)
		}
	}()
}

// sendMessageSync sends WhatsApp message synchronously (internal use)
func (s *WhatsappService) sendMessageSync(target, message string) error {
	if s.Token == "" {
		return fmt.Errorf("Fonnte token not configured")
	}

	// Normalize phone number (08xx -> 628xx)
	normalizedTarget := s.normalizePhoneNumber(target)

	// Prepare request body
	payload := map[string]string{
		"target":  normalizedTarget,
		"message": message,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %v", err)
	}

	// Create HTTP request
	req, err := http.NewRequest("POST", s.URL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("Authorization", s.Token)
	req.Header.Set("Content-Type", "application/json")

	// Send request with timeout
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %v", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Fonnte API returned status %d", resp.StatusCode)
	}

	return nil
}

// normalizePhoneNumber converts phone number to international format (628xx)
func (s *WhatsappService) normalizePhoneNumber(phone string) string {
	// Remove spaces and special characters
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")
	phone = strings.ReplaceAll(phone, "+", "")

	// Convert 08xx to 628xx
	if strings.HasPrefix(phone, "08") {
		phone = "62" + phone[1:]
	}

	// If already starts with 62, keep as is
	if strings.HasPrefix(phone, "62") {
		return phone
	}

	// If starts with 8, add 62
	if strings.HasPrefix(phone, "8") {
		return "62" + phone
	}

	return phone
}

// IsConfigured checks if WhatsApp service is properly configured
func (s *WhatsappService) IsConfigured() bool {
	return s.Token != ""
}
