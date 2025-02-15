package unit_test

import (
	"avito-shop/internal/lib/jwt"
	"avito-shop/internal/repository"
	"avito-shop/internal/services"
	"avito-shop/internal/tests/mocks"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
	"log/slog"
)

func TestAuthService_Login_NewUser(t *testing.T) {
	ctx := context.Background()

	authRepo := new(mocks.AuthRepositoryMock)
	redisMock := new(mocks.RedisClientMock)
	jwtGen := jwt.NewGenerator("secret", 0, 0) // TTL=0 для наглядности

	logger := slog.Default()
	authService := services.NewAuthService(logger, authRepo, redisMock, jwtGen)

	authRepo.On("LoginUser", ctx, "username", "newuser").
		Return("", []byte(nil), repository.ErrUserNotFound).Once()

	authRepo.On("SaveUser", ctx, "newuser", mock.AnythingOfType("[]uint8")).
		Return(nil).Once()

	hashedPass, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)
	authRepo.On("LoginUser", ctx, "username", "newuser").
		Return("someUserID", hashedPass, nil).Once()

	redisMock.On("StoreRefreshToken", "someUserID", mock.Anything).
		Return(nil).Once()

	access, refresh, err := authService.Login(ctx, "newuser", "password")

	assert.NoError(t, err)
	assert.NotEmpty(t, access)
	assert.NotEmpty(t, refresh)

	authRepo.AssertExpectations(t)
	redisMock.AssertExpectations(t)
}

func TestAuthService_Login_ExistingUser_WrongPass(t *testing.T) {
	ctx := context.Background()

	authRepo := new(mocks.AuthRepositoryMock)
	redisMock := new(mocks.RedisClientMock)
	jwtGen := jwt.NewGenerator("secret", 0, 0)

	logger := slog.Default()
	authService := services.NewAuthService(logger, authRepo, redisMock, jwtGen)

	hashedPass, _ := bcrypt.GenerateFromPassword([]byte("correct"), bcrypt.DefaultCost)
	authRepo.On("LoginUser", ctx, "username", "john123").
		Return("someUserID", hashedPass, nil).Once()

	access, refresh, err := authService.Login(ctx, "john123", "wrongpass")

	assert.Error(t, err, "ожидаем ошибку, так как пароль неверен")
	assert.Empty(t, access)
	assert.Empty(t, refresh)

	authRepo.AssertExpectations(t)
	redisMock.AssertExpectations(t)
}

func TestAuthService_Login_RepoError(t *testing.T) {
	ctx := context.Background()

	authRepo := new(mocks.AuthRepositoryMock)
	redisMock := new(mocks.RedisClientMock)
	jwtGen := jwt.NewGenerator("secret", 0, 0)

	logger := slog.Default()
	authService := services.NewAuthService(logger, authRepo, redisMock, jwtGen)

	authRepo.On("LoginUser", ctx, "username", "userWithError").
		Return("", []byte(nil), errors.New("db error")).Once()

	access, refresh, err := authService.Login(ctx, "userWithError", "password")
	assert.Error(t, err, "Ожидаем ошибку, если репозиторий возвращает ошибку")
	assert.Empty(t, access)
	assert.Empty(t, refresh)

	authRepo.AssertExpectations(t)
	redisMock.AssertExpectations(t)
}
