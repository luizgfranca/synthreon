package sessionmodule

import (
	"fmt"
	commonmodule "synthreon/modules/common"

	"github.com/golang-jwt/jwt/v5"
)

type TokenClaims struct {
	jwt.RegisteredClaims
	Iss   string `json:"iss"`
	Sub   string `json:"sub"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type SessionService struct {
	secret string
}

func (s *SessionService) CreateToken(session Session) (*string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"iss":   "synthreon",
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

func (s *SessionService) DecodeToken(token_str string) (*Session, error) {
	token, err := jwt.Parse(token_str, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.secret), nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, &commonmodule.GenericLogicError{Message: "unable to extract claims"}
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, &commonmodule.GenericLogicError{Message: "unable to extract claims"}
	}

	name, ok := claims["name"].(string)
	if !ok {
		return nil, &commonmodule.GenericLogicError{Message: "unable to extract claims"}
	}

	session := Session{
		Email: email,
		Name:  name,
	}

	return &session, nil
}

func NewSessionService(secret string) *SessionService {
	return &SessionService{
		secret: secret,
	}
}
