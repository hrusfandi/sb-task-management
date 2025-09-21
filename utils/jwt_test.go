package utils

import (
	"os"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hrusfandi/sb-task-management/config"
)

func setupJWTTest() {
	// Set up test configuration
	config.AppConfig = &config.Config{
		JWTSecret: "test-secret-key-for-testing",
	}
}

func TestGenerateToken(t *testing.T) {
	setupJWTTest()

	tests := []struct {
		name    string
		userID  uint
		email   string
		wantErr bool
	}{
		{
			name:    "Valid token generation",
			userID:  1,
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "Zero user ID",
			userID:  0,
			email:   "test@example.com",
			wantErr: false,
		},
		{
			name:    "Empty email",
			userID:  1,
			email:   "",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := GenerateToken(tt.userID, tt.email)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && token == "" {
				t.Error("GenerateToken() returned empty token")
			}
		})
	}
}

func TestValidateToken(t *testing.T) {
	setupJWTTest()

	// Generate a valid token for testing
	validToken, _ := GenerateToken(1, "test@example.com")

	// Create an expired token
	expiredClaims := JWTClaims{
		UserID: 1,
		Email:  "test@example.com",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(-time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now().Add(-2 * time.Hour)),
		},
	}
	expiredTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, expiredClaims)
	expiredToken, _ := expiredTokenObj.SignedString([]byte(config.AppConfig.JWTSecret))

	// Create token with wrong signing method
	wrongMethodClaims := JWTClaims{
		UserID: 1,
		Email:  "test@example.com",
	}
	wrongMethodToken := jwt.NewWithClaims(jwt.SigningMethodNone, wrongMethodClaims)
	unsignedToken, _ := wrongMethodToken.SignedString(jwt.UnsafeAllowNoneSignatureType)

	tests := []struct {
		name      string
		token     string
		wantErr   bool
		wantEmail string
	}{
		{
			name:      "Valid token",
			token:     validToken,
			wantErr:   false,
			wantEmail: "test@example.com",
		},
		{
			name:    "Invalid token format",
			token:   "invalid.token.format",
			wantErr: true,
		},
		{
			name:    "Empty token",
			token:   "",
			wantErr: true,
		},
		{
			name:    "Expired token",
			token:   expiredToken,
			wantErr: true,
		},
		{
			name:    "Token with wrong signature",
			token:   unsignedToken,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			claims, err := ValidateToken(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ValidateToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && claims.Email != tt.wantEmail {
				t.Errorf("ValidateToken() email = %v, want %v", claims.Email, tt.wantEmail)
			}
		})
	}
}

func TestExtractUserID(t *testing.T) {
	setupJWTTest()

	// Generate tokens for testing
	validToken, _ := GenerateToken(42, "test@example.com")

	tests := []struct {
		name    string
		token   string
		want    uint
		wantErr bool
	}{
		{
			name:    "Valid token",
			token:   validToken,
			want:    42,
			wantErr: false,
		},
		{
			name:    "Invalid token",
			token:   "invalid.token",
			want:    0,
			wantErr: true,
		},
		{
			name:    "Empty token",
			token:   "",
			want:    0,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ExtractUserID(tt.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExtractUserID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("ExtractUserID() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestJWTClaimsExpiration(t *testing.T) {
	setupJWTTest()

	// Generate a token
	token, err := GenerateToken(1, "test@example.com")
	if err != nil {
		t.Fatalf("Failed to generate token: %v", err)
	}

	// Validate and check expiration
	claims, err := ValidateToken(token)
	if err != nil {
		t.Fatalf("Failed to validate token: %v", err)
	}

	// Check if expiration is set correctly (should be 24 hours from now)
	expectedExpiry := time.Now().Add(24 * time.Hour)
	actualExpiry := claims.ExpiresAt.Time

	// Allow 1 minute difference for test execution time
	diff := expectedExpiry.Sub(actualExpiry)
	if diff > time.Minute || diff < -time.Minute {
		t.Errorf("Token expiration not set correctly. Expected around %v, got %v", expectedExpiry, actualExpiry)
	}
}

func TestMain(m *testing.M) {
	// Setup
	setupJWTTest()

	// Run tests
	code := m.Run()

	// Exit
	os.Exit(code)
}