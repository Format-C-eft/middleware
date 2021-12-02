package token

import (
	"testing"

	"github.com/dgrijalva/jwt-go/v4"
	uuid "github.com/satori/go.uuid"
)

func TestValidateToken(t *testing.T) {

	tokenPass := "AsDfGhJkL123"

	tokenJWT, _ := CreateToken(tokenPass)

	if _, err := CheckValid(tokenJWT, tokenPass); err != nil {
		t.Errorf("Func CheckValid could not decrypt the correct token - %v", err)
	}

	if _, err := CheckValid(tokenJWT, "asd"); err == nil {
		t.Errorf("Func CheckValid decrypted the token with an incorrect password")
	}

	sessionUID := uuid.NewV4().String()

	tk := token{}
	tk.UID = sessionUID
	tk.IssuedAt = jwt.Now()

	tokenNull := jwt.NewWithClaims(jwt.SigningMethodNone, tk)
	tokenString, _ := tokenNull.SigningString()

	if _, err := CheckValid(tokenString, tokenPass); err == nil {
		t.Errorf("Func CheckValid vulnerable to null encryption")
	}

}
