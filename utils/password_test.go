package utils

import (
	"testing"
)

func TestHashPassword(t *testing.T) {
	tests := []struct {
		name     string
		password string
		wantErr  bool
	}{
		{
			name:     "Valid password",
			password: "password123",
			wantErr:  false,
		},
		{
			name:     "Empty password",
			password: "",
			wantErr:  false, // bcrypt accepts empty strings
		},
		{
			name:     "Long password",
			password: "verylongpasswordverylongpasswordverylongpasswordverylongpassword",
			wantErr:  false,
		},
		{
			name:     "Special characters",
			password: "p@$$w0rd!#",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			hash, err := HashPassword(tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("HashPassword() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if hash == "" {
					t.Error("HashPassword() returned empty hash")
				}
				if hash == tt.password {
					t.Error("HashPassword() returned plain password, not hashed")
				}
			}
		})
	}
}

func TestComparePassword(t *testing.T) {
	password := "testpassword123"
	hashedPassword, _ := HashPassword(password)

	tests := []struct {
		name           string
		hashedPassword string
		password       string
		wantErr        bool
	}{
		{
			name:           "Correct password",
			hashedPassword: hashedPassword,
			password:       password,
			wantErr:        false,
		},
		{
			name:           "Incorrect password",
			hashedPassword: hashedPassword,
			password:       "wrongpassword",
			wantErr:        true,
		},
		{
			name:           "Empty password",
			hashedPassword: hashedPassword,
			password:       "",
			wantErr:        true,
		},
		{
			name:           "Invalid hash",
			hashedPassword: "invalid-hash",
			password:       password,
			wantErr:        true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ComparePassword(tt.hashedPassword, tt.password)
			if (err != nil) != tt.wantErr {
				t.Errorf("ComparePassword() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestPasswordHashingConsistency(t *testing.T) {
	password := "consistentPassword123"

	// Hash the same password twice
	hash1, err1 := HashPassword(password)
	if err1 != nil {
		t.Fatalf("Failed to hash password first time: %v", err1)
	}

	hash2, err2 := HashPassword(password)
	if err2 != nil {
		t.Fatalf("Failed to hash password second time: %v", err2)
	}

	// Hashes should be different (bcrypt adds salt)
	if hash1 == hash2 {
		t.Error("Two hashes of the same password should be different due to salt")
	}

	// Both hashes should validate against the original password
	if err := ComparePassword(hash1, password); err != nil {
		t.Errorf("First hash failed to validate: %v", err)
	}

	if err := ComparePassword(hash2, password); err != nil {
		t.Errorf("Second hash failed to validate: %v", err)
	}
}