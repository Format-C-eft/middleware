package token

import (
	"errors"

	"github.com/dgrijalva/jwt-go/v4"
	uuid "github.com/satori/go.uuid"
)

type token struct {
	UID string `json:"UID"`
	jwt.StandardClaims
}

// CreateToken - create new token
func CreateToken(pass string) (tokenString, sessionUID string) {

	sessionUID = uuid.NewV4().String()

	tk := token{}
	tk.UID = sessionUID
	tk.IssuedAt = jwt.Now()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, tk)
	tokenString, _ = token.SignedString([]byte(pass))

	return
}

// CheckValid - token validation
func CheckValid(tokenStr, pass string) (*token, error) {

	jwtToken, err := jwt.ParseWithClaims(tokenStr, &token{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(pass), nil
	})

	// Invalid JWT token
	if err != nil {
		return &token{}, err
	}

	if jwtToken.Method != jwt.SigningMethodHS256 {
		return &token{}, errors.New("signing method not allowed")
	}

	claims, ok := jwtToken.Claims.(*token)
	if !ok || !jwtToken.Valid {
		return &token{}, errors.New("not correct JWT token")
	}

	return claims, nil
}
