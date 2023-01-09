package jmanager

import (
	"testing"
)

func TestJWT(t *testing.T) {

	// run tests
	t.Run("Test generatation/validation of JWT token", func(t *testing.T) {
		key := "verysecretkey"
		var userID int64 = 12

		token, err := Generate(userID, key)
		if err != nil {
			t.Fatalf("cant generate token: %s", err.Error())
		}

		validatedToken, err := Validate(token.Token, key)
		if err != nil {
			t.Fatalf("Token validation was vailed: %s", err.Error())
		}

		if validatedToken.Claims.UserID != userID {
			t.Errorf("UserId mismatch: have %d, want %d",
				validatedToken.Claims.UserID, userID)
		}

	})
}
