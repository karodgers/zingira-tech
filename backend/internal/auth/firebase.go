package auth

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	"google.golang.org/api/option"
)

// AuthService handles Firebase authentication
type AuthService struct {
	client *auth.Client
}

// NewAuthService initializes the AuthService
func NewAuthService() (*AuthService, error) {
	ctx := context.Background()

	// Get the directory where this file is located
	_, b, _, _ := runtime.Caller(0)
	projectRoot := filepath.Join(filepath.Dir(b), "../..")

	// Check if the credentials file path is provided via an environment variable
	credentialsFile := os.Getenv("FIREBASE_CREDENTIALS")
	if credentialsFile == "" {
		// If not, fall back to a default path within the project structure
		credentialsFile = filepath.Join(projectRoot, "config", "firebase", "serviceAccountKey.json")
	}

	// Initialize Firebase app with credentials file
	opt := option.WithCredentialsFile(credentialsFile)
	app, err := firebase.NewApp(ctx, nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %v", err)
	}

	// Create the Auth client
	authClient, err := app.Auth(ctx)
	if err != nil {
		return nil, fmt.Errorf("error getting Auth client: %v", err)
	}

	// Log successful initialization
	log.Printf("Firebase Auth initialized successfully using credentials from: %s", credentialsFile)

	// Return the AuthService instance
	return &AuthService{client: authClient}, nil
}

// VerifyIDToken verifies the Firebase ID token
func (as *AuthService) VerifyIDToken(ctx context.Context, idToken string) (*auth.Token, error) {
	// Verify the token with Firebase Auth
	token, err := as.client.VerifyIDToken(ctx, idToken)
	if err != nil {
		return nil, fmt.Errorf("error verifying ID token: %v", err)
	}
	return token, nil
}

// ExtractClaims extracts custom claims from a verified token
func (as *AuthService) ExtractClaims(token *auth.Token) map[string]interface{} {
	// Extract and return the claims from the token
	return token.Claims
}
