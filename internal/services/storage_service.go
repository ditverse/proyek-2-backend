package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"backend-sarpras/internal/config"
)

func UploadPDFToSupabase(objectPath string, fileBytes []byte) error {
	cfg := config.GetSupabaseConfig()
	if cfg.URL == "" || cfg.ServiceKey == "" || cfg.Bucket == "" {
		return fmt.Errorf("supabase config incomplete")
	}

	url := fmt.Sprintf("%s/storage/v1/object/%s/%s", cfg.URL, cfg.Bucket, objectPath)

	// Use PUT instead of POST to allow upsert (create or replace)
	req, err := http.NewRequest("PUT", url, bytes.NewReader(fileBytes))
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+cfg.ServiceKey)
	req.Header.Set("Content-Type", "application/pdf")
	req.Header.Set("x-upsert", "true") // Allow overwrite if file exists

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("upload gagal: %s", string(body))
	}

	return nil
}

func GenerateSignedURL(objectPath string) (string, error) {
	cfg := config.GetSupabaseConfig()
	if cfg.URL == "" || cfg.ServiceKey == "" || cfg.Bucket == "" {
		return "", fmt.Errorf("supabase config incomplete")
	}

	url := fmt.Sprintf("%s/storage/v1/object/sign/%s/%s", cfg.URL, cfg.Bucket, objectPath)

	payload := fmt.Sprintf(`{"expiresIn": %d}`, cfg.Expire)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(payload)))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+cfg.ServiceKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("gagal generate signed URL: %s", string(body))
	}

	body, _ := io.ReadAll(resp.Body)
	var result struct {
		SignedURL string `json:"signedURL"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("invalid signed url response: %s", string(body))
	}
	if result.SignedURL == "" {
		return "", fmt.Errorf("signed url empty")
	}

	if strings.HasPrefix(result.SignedURL, "http") {
		return result.SignedURL, nil
	}

	signedPath := result.SignedURL
	if !strings.HasPrefix(signedPath, "/") {
		signedPath = "/" + signedPath
	}
	if !strings.HasPrefix(signedPath, "/storage/v1/") {
		signedPath = "/storage/v1" + signedPath
	}

	return fmt.Sprintf("%s%s", cfg.URL, signedPath), nil
}

// MoveFile moves a file from oldPath to newPath in Supabase Storage
// This is done by: 1) Copy file to new location, 2) Delete old file
func MoveFile(oldPath, newPath string) error {
	cfg := config.GetSupabaseConfig()
	if cfg.URL == "" || cfg.ServiceKey == "" || cfg.Bucket == "" {
		return fmt.Errorf("supabase config incomplete")
	}

	// Step 1: Download file from old path
	downloadURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", cfg.URL, cfg.Bucket, oldPath)

	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create download request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+cfg.ServiceKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to download file (status %d): %s", resp.StatusCode, string(body))
	}

	// Read file content
	fileBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read file content: %w", err)
	}

	// Step 2: Upload file to new path
	if err := UploadPDFToSupabase(newPath, fileBytes); err != nil {
		return fmt.Errorf("failed to upload to new path: %w", err)
	}

	// Step 3: Delete old file
	deleteURL := fmt.Sprintf("%s/storage/v1/object/%s/%s", cfg.URL, cfg.Bucket, oldPath)

	delReq, err := http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		// File uploaded to new path, but old file not deleted - not critical
		return nil
	}
	delReq.Header.Set("Authorization", "Bearer "+cfg.ServiceKey)

	delResp, err := client.Do(delReq)
	if err != nil {
		// File uploaded to new path, but old file not deleted - not critical
		return nil
	}
	defer delResp.Body.Close()

	return nil
}
