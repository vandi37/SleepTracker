package tokens

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type TokenService struct {
	secret []byte
	exr    time.Duration
}

func New(secret string, exr time.Duration) *TokenService {
	return &TokenService{
		[]byte(secret),
		exr,
	}
}

func (s *TokenService) GetSecret() []byte {
	return s.secret
}

func (s *TokenService) GetExpiration() time.Duration {
	return s.exr
}

type Data struct {
	Key   string
	Value any
}

func (s *TokenService) Generate(sub string, data ...Data) (string, error) {
	now := time.Now()
	claims := jwt.MapClaims{
		"sub": sub,
		"iat": now.Unix(),
		"exp": now.Add(s.exr).Unix(),
	}
	for _, d := range data {
		claims[d.Key] = d.Value
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(s.secret)
}

func (s *TokenService) Parse(token string) (*jwt.Token, error) {
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return s.secret, nil
	})
	if err != nil {
		return nil, err
	}
	return parsedToken, nil
}
