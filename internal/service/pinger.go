package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt"
	log "github.com/sirupsen/logrus"
	"github.com/spanwalla/docker-monitoring-backend/internal/entity"
	"github.com/spanwalla/docker-monitoring-backend/internal/repository"
	"github.com/spanwalla/docker-monitoring-backend/pkg/hasher"
	"time"
)

// TokenClaims -.
type TokenClaims struct {
	jwt.StandardClaims
	PingerId int
}

// PingerService -.
type PingerService struct {
	pingerRepo     repository.Pinger
	passwordHasher hasher.PasswordHasher
	signKey        string
	tokenTTL       time.Duration
}

// NewPingerService -.
func NewPingerService(pingerRepo repository.Pinger, passwordHasher hasher.PasswordHasher, signKey string, tokenTTL time.Duration) *PingerService {
	return &PingerService{
		pingerRepo,
		passwordHasher,
		signKey,
		tokenTTL,
	}
}

// CreatePinger -.
func (s *PingerService) CreatePinger(ctx context.Context, input PingerCreateInput) (int, error) {
	pinger := entity.Pinger{
		Name:     input.Name,
		Password: s.passwordHasher.Hash(input.Password),
	}

	pingerId, err := s.pingerRepo.CreatePinger(ctx, pinger)
	if err != nil {
		if errors.Is(err, repository.ErrAlreadyExists) {
			return 0, ErrPingerAlreadyExists
		}
		log.Errorf("PingerService.CreatePinger - s.pingerRepo.CreatePinger: %v", err)
		return 0, ErrCannotCreatePinger
	}
	return pingerId, nil
}

// GenerateToken -.
func (s *PingerService) GenerateToken(ctx context.Context, input PingerGenerateTokenInput) (string, error) {
	pinger, err := s.pingerRepo.GetPingerByNameAndPassword(ctx, input.Name, s.passwordHasher.Hash(input.Password))
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return "", ErrPingerNotFound
		}
		log.Errorf("PingerService.GenerateToken - s.pingerRepo.GetPingerByNameAndPassword: %v", err)
		return "", ErrCannotGetPinger
	}

	// Generate JWT
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &TokenClaims{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(s.tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		PingerId: pinger.Id,
	})

	// Sign token
	tokenString, err := token.SignedString([]byte(s.signKey))
	if err != nil {
		log.Errorf("PingerService.GenerateToken - token.SignedString: %v", err)
		return "", ErrCannotSignToken
	}

	return tokenString, nil
}

// ParseToken -.
func (s *PingerService) ParseToken(accessToken string) (int, error) {
	token, err := jwt.ParseWithClaims(accessToken, &TokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(s.signKey), nil
	})

	if err != nil {
		return 0, ErrCannotParseToken
	}

	claims, ok := token.Claims.(*TokenClaims)
	if !ok {
		return 0, ErrCannotParseToken
	}

	return claims.PingerId, nil
}
