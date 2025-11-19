package utils

import (
	"crypto/rand"
	"encoding/base32"
	"fmt"
	"strings"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

// GenerateMFASecret generates a new TOTP secret for MFA
func GenerateMFASecret(email string) (string, string, error) {
	key, err := totp.Generate(totp.GenerateOpts{
		Issuer:      "SessionManagement",
		AccountName: email,
		Period:      30,
		Digits:      otp.DigitsSix,
		Algorithm:   otp.AlgorithmSHA1,
	})
	if err != nil {
		return "", "", fmt.Errorf("failed to generate TOTP key: %w", err)
	}

	// Return the secret and QR code URL
	return key.Secret(), key.URL(), nil
}

// ValidateMFACode validates a TOTP code against a secret
func ValidateMFACode(code, secret string) bool {
	// Remove spaces and convert to uppercase
	code = strings.TrimSpace(code)
	code = strings.ToUpper(code)

	// Validate the code
	valid := totp.Validate(code, secret)
	return valid
}

// GenerateBackupCodes generates 10 random backup codes
func GenerateBackupCodes() ([]string, error) {
	codes := make([]string, 10)
	
	for i := 0; i < 10; i++ {
		// Generate 10 random bytes
		randomBytes := make([]byte, 10)
		_, err := rand.Read(randomBytes)
		if err != nil {
			return nil, fmt.Errorf("failed to generate random bytes: %w", err)
		}
		
		// Encode to base32 and take first 10 characters
		code := base32.StdEncoding.EncodeToString(randomBytes)
		code = code[:10]
		
		// Format as XXXXX-XXXXX
		codes[i] = fmt.Sprintf("%s-%s", code[:5], code[5:10])
	}
	
	return codes, nil
}