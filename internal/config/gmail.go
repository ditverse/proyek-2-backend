package config

import (
	"os"
)

// GmailConfig holds Gmail API OAuth2 configuration
type GmailConfig struct {
	ClientID     string
	ClientSecret string
	RefreshToken string
	SenderEmail  string
	SenderName   string
	Enabled      bool
}

// GetGmailConfig returns Gmail API configuration from environment variables
func GetGmailConfig() GmailConfig {
	clientID := os.Getenv("GMAIL_CLIENT_ID")
	clientSecret := os.Getenv("GMAIL_CLIENT_SECRET")
	refreshToken := os.Getenv("GMAIL_REFRESH_TOKEN")
	senderEmail := os.Getenv("GMAIL_SENDER_EMAIL")
	senderName := os.Getenv("GMAIL_SENDER_NAME")

	if senderName == "" {
		senderName = "Sarpras System"
	}

	// Email is enabled only if all required credentials are present
	enabled := clientID != "" && clientSecret != "" && refreshToken != "" && senderEmail != ""

	return GmailConfig{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RefreshToken: refreshToken,
		SenderEmail:  senderEmail,
		SenderName:   senderName,
		Enabled:      enabled,
	}
}
