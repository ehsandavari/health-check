package services

import (
	"context"
	"health-check/presentation/grpc/proto/aaa"
)

type AuthService struct {
}

func NewAuthService() *AuthService {
	return &AuthService{}
}

func (g *AuthService) ValidateToken(ctx context.Context, request *aaa.ValidateTokenRequest) (*aaa.ValidateTokenResponse, error) {
	//TODO implement me
	panic("implement me")
}
