package services

import (
	"backend-sarpras/internal/config"
	"fmt"
	"log"
	"strconv"

	"gopkg.in/gomail.v2"
)

// EmailService handles email sending via SMTP
type EmailService struct {
	Host     string
	Port     int
	Email    string
	Password string
}

// NewEmailService creates a new EmailService from config
func NewEmailService(cfg *config.Config) *EmailService {
	port, err := strconv.Atoi(cfg.SMTPPort)
	if err != nil {
		port = 587 // default SMTP port
	}

	return &EmailService{
		Host:     cfg.SMTPHost,
		Port:     port,
		Email:    cfg.SMTPEmail,
		Password: cfg.SMTPPassword,
	}
}

// SendEmail sends an email asynchronously using goroutine
// This method does not block the main thread
func (s *EmailService) SendEmail(to, subject, body string) {
	go func() {
		if err := s.sendEmailSync(to, subject, body); err != nil {
			log.Printf("❌ Email failed to send to %s: %v", to, err)
		} else {
			log.Printf("✅ Email sent successfully to %s", to)
		}
	}()
}

// SendEmailHTML sends an HTML email asynchronously
func (s *EmailService) SendEmailHTML(to, subject, htmlBody string) {
	go func() {
		if err := s.sendEmailHTMLSync(to, subject, htmlBody); err != nil {
			log.Printf("❌ Email HTML failed to send to %s: %v", to, err)
		} else {
			log.Printf("✅ Email HTML sent successfully to %s", to)
		}
	}()
}

// sendEmailSync sends email synchronously (internal use)
func (s *EmailService) sendEmailSync(to, subject, body string) error {
	if s.Email == "" || s.Password == "" {
		return fmt.Errorf("SMTP credentials not configured")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.Email)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/plain", body)

	d := gomail.NewDialer(s.Host, s.Port, s.Email, s.Password)

	return d.DialAndSend(m)
}

// sendEmailHTMLSync sends HTML email synchronously (internal use)
func (s *EmailService) sendEmailHTMLSync(to, subject, htmlBody string) error {
	if s.Email == "" || s.Password == "" {
		return fmt.Errorf("SMTP credentials not configured")
	}

	m := gomail.NewMessage()
	m.SetHeader("From", s.Email)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", htmlBody)

	d := gomail.NewDialer(s.Host, s.Port, s.Email, s.Password)

	return d.DialAndSend(m)
}

// IsConfigured checks if email service is properly configured
func (s *EmailService) IsConfigured() bool {
	return s.Email != "" && s.Password != ""
}
