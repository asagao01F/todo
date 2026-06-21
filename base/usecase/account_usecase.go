package usecase

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"todo/base/model"
)

// AccountRepository インターフェースの宣言
type AccountRepository interface {
	Create(ctx context.Context, account *model.Account) error
	FindByID(ctx context.Context, id uint) (*model.Account, error)
	FindByEmail(ctx context.Context, email string) (*model.Account, error)
}

type AccountUsecase struct {
	accountRepo AccountRepository
}

func NewAccountUsecase(accountRepo AccountRepository) *AccountUsecase {
	return &AccountUsecase{accountRepo: accountRepo}
}

// RegisterAccount: アカウントの新規登録
func (u *AccountUsecase) RegisterAccount(ctx context.Context, account *model.Account) error {
	// バリデーション
	if strings.TrimSpace(account.Email) == "" || strings.TrimSpace(account.Password) == "" {
		return errors.New("email and password cannot be empty")
	}

	// 【ビジネスロジック】Emailの重複チェック
	existing, err := u.accountRepo.FindByEmail(ctx, account.Email)
	if err != nil {
		return fmt.Errorf("failed to check email duplication: %w", err)
	}
	if existing != nil {
		return errors.New("email address is already registered")
	}

	// 保存を実行 (実態は PostgresAccountRepository.Create)
	if err := u.accountRepo.Create(ctx, account); err != nil {
		return fmt.Errorf("account registration failed: %w", err)
	}
	return nil
}

// GetAccountByID: アカウント情報の取得
func (u *AccountUsecase) GetAccountByID(ctx context.Context, id uint) (*model.Account, error) {
	account, err := u.accountRepo.FindByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch account: %w", err)
	}
	if account == nil {
		return nil, errors.New("account not found")
	}
	return account, nil
}
