package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
)

// Konfigurasi OAuth - sesuaikan dengan credentials Anda
var oauthConfig = &oauth2.Config{
	ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
	ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
	RedirectURL:  "http://localhost:8080/callback",
	Scopes:       []string{gmail.GmailSendScope},
	Endpoint:     google.Endpoint,
}

var tokenChan = make(chan *oauth2.Token)

func main() {
	separator := strings.Repeat("=", 60)

	// Start local server to handle callback
	http.HandleFunc("/callback", handleCallback)
	go func() {
		log.Println("Starting local server on http://localhost:8080")
		if err := http.ListenAndServe(":8080", nil); err != nil {
			log.Fatal(err)
		}
	}()

	// Generate authorization URL
	authURL := oauthConfig.AuthCodeURL("state", oauth2.AccessTypeOffline, oauth2.SetAuthURLParam("prompt", "consent"))

	fmt.Println("\n" + separator)
	fmt.Println("GMAIL OAUTH TOKEN GENERATOR")
	fmt.Println(separator)
	fmt.Println("\n1. Buka URL berikut di browser Anda:")
	fmt.Println(authURL)
	fmt.Println("\n2. Login dengan akun Gmail yang akan digunakan untuk mengirim email")
	fmt.Println("3. Izinkan akses aplikasi")
	fmt.Println("4. Tunggu redirect ke localhost...")
	fmt.Println("\n" + separator)

	// Wait for token from callback
	token := <-tokenChan

	fmt.Println("\nTOKEN BERHASIL DI-GENERATE!")
	fmt.Println(separator)
	fmt.Println("\nSalin Refresh Token berikut ke file .env Anda:")
	fmt.Println("\nGMAIL_REFRESH_TOKEN=" + token.RefreshToken)
	fmt.Println("\n" + separator)

	// Save token to file as backup
	saveToken("gmail_token.json", token)
	fmt.Println("\nToken juga disimpan di: gmail_token.json")
	fmt.Println("\nSelesai! Restart server backend Anda setelah update .env")
}

func handleCallback(w http.ResponseWriter, r *http.Request) {
	code := r.URL.Query().Get("code")
	if code == "" {
		http.Error(w, "No code in callback", http.StatusBadRequest)
		return
	}

	// Exchange code for token
	token, err := oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		http.Error(w, "Failed to exchange token: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Send success response to browser
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprintf(w, `
		<html>
		<head><title>OAuth Success</title></head>
		<body style="font-family: Arial, sans-serif; text-align: center; padding: 50px;">
			<h1 style="color: green;">Authorization Successful!</h1>
			<p>Token berhasil di-generate.</p>
			<p>Anda bisa menutup tab ini dan kembali ke terminal.</p>
		</body>
		</html>
	`)

	// Send token to main goroutine
	tokenChan <- token
}

func saveToken(path string, token *oauth2.Token) {
	f, err := os.Create(path)
	if err != nil {
		log.Printf("Unable to save token: %v", err)
		return
	}
	defer f.Close()
	json.NewEncoder(f).Encode(token)
}
