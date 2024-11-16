package service

import (
	"fmt"
	"platformlab/controlpanel/model"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Session struct {
	secret string
}

func (s *Session) CreateToken(session model.Session) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"iss":   "platformlab",
		"sub":   session.Email,
		"email": session.Email,
		"name":  session.Name,
	})

	token_str, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return nil, err
	}

	return &token_str, nil
}

func (s *Session) DecodeToken(token_str string) (*model.Session, error) {
	token, err := jwt.Parse(token_str, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return s.secret, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return nil, &model.GenericLogicError{Message: "unable to extract claims"}
	}

	session := model.Session{
		Email: claims.Email,
		Name:  claims.Name,
	}

	return &session, nil
}

func NewSessionService(secret string) *Session {
	return &Session{
		secret: secret,
	}
}
