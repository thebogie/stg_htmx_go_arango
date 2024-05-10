package auth

import (
	"testing"
)

func TestCheckJWTCreation(t *testing.T) {
	var jwtheader JwtHeader

	t.Setenv("JWT_ACCESS_KEY", "2c46f7651a176169c8d2f7aee60cce3da874501d")
	emailToTest := "john@test.com"

	token, err := jwtheader.GenerateToken(emailToTest)
	if err != nil {
		t.Errorf("Error generating token %v", err)
		return
	}

	result := jwtheader.CheckToken(token)
	if result != emailToTest {
		t.Errorf("Expected %s, got %s", emailToTest, result)
	}

}
