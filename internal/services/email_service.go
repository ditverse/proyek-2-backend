package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"strings"

	"backend-sarpras/internal/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"
)

// EmailService handles sending emails via Gmail API
type EmailService struct {
	gmailService *gmail.Service
	senderEmail  string
	senderName   string
	enabled      bool
}

// NewEmailService creates a new email service with OAuth2 Gmail API
func NewEmailService() *EmailService {
	cfg := config.GetGmailConfig()

	if !cfg.Enabled {
		log.Println("⚠️  Email service disabled: Gmail credentials not configured")
		return &EmailService{enabled: false}
	}

	// Create OAuth2 config
	oauth2Config := &oauth2.Config{
		ClientID:     cfg.ClientID,
		ClientSecret: cfg.ClientSecret,
		Endpoint:     google.Endpoint,
		Scopes:       []string{gmail.GmailSendScope},
	}

	// Create token from refresh token
	token := &oauth2.Token{
		RefreshToken: cfg.RefreshToken,
	}

	// Create OAuth2 client
	ctx := context.Background()
	client := oauth2Config.Client(ctx, token)

	// Create Gmail service
	gmailSvc, err := gmail.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		log.Printf("❌ Failed to create Gmail service: %v", err)
		return &EmailService{enabled: false}
	}

	log.Println("✅ Email service initialized successfully")
	return &EmailService{
		gmailService: gmailSvc,
		senderEmail:  cfg.SenderEmail,
		senderName:   cfg.SenderName,
		enabled:      true,
	}
}

// IsEnabled returns whether email service is enabled
func (s *EmailService) IsEnabled() bool {
	return s.enabled
}

// SendEmail sends an email using Gmail API
func (s *EmailService) SendEmail(to, subject, htmlBody string) error {
	if !s.enabled {
		return fmt.Errorf("email service is not enabled")
	}

	// Build the email message
	fromHeader := fmt.Sprintf("From: %s <%s>\r\n", s.senderName, s.senderEmail)
	toHeader := fmt.Sprintf("To: %s\r\n", to)
	subjectHeader := fmt.Sprintf("Subject: %s\r\n", subject)
	mimeHeader := "MIME-Version: 1.0\r\n"
	contentType := "Content-Type: text/html; charset=UTF-8\r\n"

	message := fromHeader + toHeader + subjectHeader + mimeHeader + contentType + "\r\n" + htmlBody

	// Encode to base64
	encodedMessage := base64.URLEncoding.EncodeToString([]byte(message))

	// Create Gmail message
	gmailMessage := &gmail.Message{
		Raw: encodedMessage,
	}

	// Send email
	_, err := s.gmailService.Users.Messages.Send("me", gmailMessage).Do()
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	log.Printf("📧 Email sent successfully to: %s", to)
	return nil
}

// SendEmailToMultiple sends email to multiple recipients
func (s *EmailService) SendEmailToMultiple(recipients []string, subject, htmlBody string) []error {
	var errors []error
	for _, to := range recipients {
		if err := s.SendEmail(to, subject, htmlBody); err != nil {
			errors = append(errors, fmt.Errorf("failed to send to %s: %w", to, err))
		}
	}
	return errors
}

// FormatDate formats time for email display
func FormatDate(t interface{ Format(string) string }) string {
	return t.Format("02 Januari 2006, 15:04 WIB")
}

// FormatDateShort formats date only
func FormatDateShort(t interface{ Format(string) string }) string {
	return t.Format("02 Januari 2006")
}

// EscapeHTML escapes HTML special characters
func EscapeHTML(s string) string {
	s = strings.ReplaceAll(s, "&", "&amp;")
	s = strings.ReplaceAll(s, "<", "&lt;")
	s = strings.ReplaceAll(s, ">", "&gt;")
	s = strings.ReplaceAll(s, "\"", "&quot;")
	s = strings.ReplaceAll(s, "'", "&#39;")
	return s
}
