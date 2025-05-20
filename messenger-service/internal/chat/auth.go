package chat

import (
	"context"
	"fmt"
	"time"

	ssov6 "github.com/Nipup1/SSO_Protos/gen/go/sso"
	"github.com/golang-jwt/jwt/v5"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthService struct {
	Secret string
	Conn   *grpc.ClientConn
	Client ssov6.AuhtClient
}

func NewAuthService(secret string, gRPCAddres string) *AuthService {
	conn, err := grpc.NewClient(
		gRPCAddres,
		grpc.WithTransportCredentials(insecure.NewCredentials()), // Используем insecure для примера
	)
	if err != nil {
		panic("sso is not connected")
	}

	client := ssov6.NewAuhtClient(conn)
	
	return &AuthService{
		Secret: secret,
		Client: client,
		Conn: conn,
	}
}

func (s *AuthService) ValidateToken(tokenString string) (int, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(s.Secret), nil
	})
	if err != nil {
		return 0, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return 0, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return 0, fmt.Errorf("invalid claims format")
	}

	uid, ok := claims["uid"].(float64)
	if !ok {
		return 0, fmt.Errorf("userID not found or invalid")
	}

	exp, ok := claims["exp"].(float64)
	if !ok || int64(exp) < time.Now().Unix() {
		return 0, fmt.Errorf("token expired")
	}

	return int(uid), nil
}

func (c *AuthService) Close() error {
	return c.Conn.Close()
}

func (c *AuthService) Register(ctx context.Context, email, password, fullName string) (int64, error) {
	const op = "authService.Register"

	req := &ssov6.RegisterRequest{
		Email:     email,
		Password:  password,
		FullName: fullName,
	}

	resp, err := c.Client.Register(ctx, req)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return resp.GetUserId(), nil
}

func (c *AuthService) Login(ctx context.Context, email, password string) (string, error) {
	const op = "authService.Login"

	req := &ssov6.LoginRequest{
		Email:   email,
		Password: password,
	}

	resp, err := c.Client.Login(ctx, req)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return resp.GetToken(), nil
}

func (c *AuthService) Users(ctx context.Context) ([]*ssov6.User, error) {
	const op = "authService.Users"

	req := &ssov6.UsersRequest{}

	resp, err := c.Client.Users(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resp.GetUsers(), nil
}

func (c *AuthService) User(ctx context.Context, userId int64) (*ssov6.UserResponse, error) {
	const op = "authService.User"

	req := &ssov6.UserRequest{
		UserId: userId,
	}

	resp, err := c.Client.User(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return resp, nil
}