package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/airdanapi/API_Integrator_gateway/backend/internal/model"
	"github.com/golang-jwt/jwt/v5"
)

type Claims struct {
	Username string     `json:"username"`
	Role     model.Role `json:"role"`
	AppName  string     `json:"app_name"`
	jwt.RegisteredClaims
}

type JWTService struct {
	secret []byte
	issuer string
	ttl    time.Duration
	now    func() time.Time
}

func NewJWTService(
	secret string,
	issuer string,
	ttl time.Duration,
	now func() time.Time,
) *JWTService {
	if now == nil {
		now = time.Now
	}
	return &JWTService{
		secret: []byte(secret),
		issuer: issuer,
		ttl:    ttl,
		now:    now,
	}
}

func (service *JWTService) Generate(user model.User) (string, int64, error) {
	now := service.now().UTC()
	expiresAt := now.Add(service.ttl)
	claims := Claims{
		Username: user.Username,
		Role:     user.Role,
		AppName:  user.AppName,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   strconv.FormatInt(user.ID, 10),
			Issuer:    service.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString(service.secret)
	if err != nil {
		return "", 0, fmt.Errorf("sign token: %w", err)
	}
	return signedToken, int64(service.ttl / time.Second), nil
}

func (service *JWTService) Validate(tokenString string) (Claims, error) {
	claims := Claims{}
	token, err := jwt.ParseWithClaims(
		tokenString,
		&claims,
		func(token *jwt.Token) (any, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return service.secret, nil
		},
		jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}),
		jwt.WithIssuer(service.issuer),
		jwt.WithExpirationRequired(),
		jwt.WithIssuedAt(),
		jwt.WithTimeFunc(service.now),
	)
	if err != nil {
		return Claims{}, fmt.Errorf("validate token: %w", err)
	}
	if !token.Valid {
		return Claims{}, fmt.Errorf("validate token: token is invalid")
	}
	if claims.Subject == "" ||
		claims.Username == "" ||
		claims.Role == "" ||
		claims.AppName == "" {
		return Claims{}, fmt.Errorf("validate token: required claims are missing")
	}
	return claims, nil
}
