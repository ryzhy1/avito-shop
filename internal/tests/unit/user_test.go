package unit

import (
	"avito-shop/internal/domain/dto"
	"avito-shop/internal/services"
	"avito-shop/internal/tests/mocks"
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"log/slog"
)

func TestUserService_GetUserInfo_Success(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	repo := new(mocks.UserRepositoryMock)
	repo.On("GetUserById", ctx, userID).
		Return(dto.UserDTO{Coins: 5000}, nil)
	repo.On("GetUserPurchases", ctx, userID).
		Return([]dto.PurchaseDTO{
			{Merch: "pen", Amount: 2},
		}, nil)
	repo.On("GetCoinTransactions", ctx, userID).
		Return(dto.TransactionDTO{}, nil)

	logger := slog.Default()
	svc := services.NewUserService(logger, repo)

	info, err := svc.GetUserInfo(ctx, userID)

	assert.NoError(t, err)
	assert.Equal(t, 5000, info.Coins)
	assert.Len(t, info.Inventory, 1)
	assert.Equal(t, "pen", info.Inventory[0].Merch)

	repo.AssertExpectations(t)
}

func TestUserService_GetUserInfo_RepoError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	repo := new(mocks.UserRepositoryMock)
	repo.On("GetUserById", ctx, userID).
		Return(dto.UserDTO{}, errors.New("db error"))

	logger := slog.Default()
	svc := services.NewUserService(logger, repo)

	info, err := svc.GetUserInfo(ctx, userID)

	assert.Error(t, err)
	assert.Empty(t, info.Inventory)
	repo.AssertExpectations(t)
}

func TestUserService_TransferCoins(t *testing.T) {
	ctx := context.Background()
	fromID := uuid.New()
	toID := uuid.New()

	repo := new(mocks.UserRepositoryMock)
	repo.On("TransferCoins", ctx, fromID, toID, 100).Return(nil)

	logger := slog.Default()
	svc := services.NewUserService(logger, repo)

	err := svc.TransferCoins(ctx, fromID, toID, 100)
	assert.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestUserService_BuyItem(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	repo := new(mocks.UserRepositoryMock)
	repo.On("BuyItem", ctx, userID, "t-shirt").Return(nil)

	logger := slog.Default()
	svc := services.NewUserService(logger, repo)

	err := svc.BuyItem(ctx, userID, "t-shirt")
	assert.NoError(t, err)

	repo.AssertExpectations(t)
}

func TestUserService_BuyItem_RepoError(t *testing.T) {
	ctx := context.Background()
	userID := uuid.New()

	repo := new(mocks.UserRepositoryMock)
	repo.On("BuyItem", ctx, userID, "t-shirt").Return(errors.New("not enough coins"))

	logger := slog.Default()
	svc := services.NewUserService(logger, repo)

	err := svc.BuyItem(ctx, userID, "t-shirt")
	assert.Error(t, err)

	repo.AssertExpectations(t)
}
