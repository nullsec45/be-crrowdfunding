package auth

import (
	"github.com/dgrijalva/jwt-go"
	"crowdfunding-api/config"
	"errors"
)

type Service interface {
	GenerateToken(userID int) (string, error)
	ValidateToken(token string) (*jwt.Token, error)
}

type jwtService struct {
	cfg *config.Config
}

func NewService(cfg *config.Config) *jwtService {
	return &jwtService{
		cfg:cfg,
	}
}

func (s *jwtService) GenerateToken(userID int) (string, error) {
	var SECRET_KEY=[]byte(s.cfg.App.JwtSecretKey)
    
	claim := jwt.MapClaims{}
	claim["user_id"] = userID
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claim)

	signedToken, err := token.SignedString(SECRET_KEY)
	if err != nil {
		return signedToken, err
	}


	return signedToken, nil
}

func (s *jwtService) ValidateToken(encodedToken string) (*jwt.Token, error) {
	var SECRET_KEY=[]byte(s.cfg.App.JwtSecretKey)

	token, err := jwt.Parse(encodedToken, func(token *jwt.Token) (interface{}, error) {
			_, ok := token.Method.(*jwt.SigningMethodHMAC)

			if !ok {
				return nil, errors.New("Invalid token")
			}

			return []byte(SECRET_KEY), nil
	})

	if err != nil {
		return token err
	}

	return token, nil
}