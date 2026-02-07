package auth

import (
	"github.com/dgrijalva/jwt-go"
	"crowdfunding-api/config"
)

type Service interface {
	GenerateToken(userID int) (string, error)
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