package repository

import (
	"context"
	"fmt"
	"todo/base/model"
	"gorm.io/gorm"
)

type PostgresAccountRepository struct {
	db *gorm.DB
}

func NewPostgresAccountRepository(db *gorm.DB) *PostgresAccountRepository {
	return &PostgresAccountRepository{db: db}
}

// 1. Create: アカウントの新規作成
// TODOと同様、実行後に自動採番されたIDや作成日時が account 構造体に書き戻されます
func (r *PostgresAccountRepository) Create(ctx context.Context, account *model.Account) error {
	result := r.db.WithContext(ctx).Create(account)
	if result.Error != nil {
		return fmt.Errorf("failed to create account: %w", result.Error)
	}
	return nil
}

// 2. FindByID: IDによるアカウントの1件取得
func (r *PostgresAccountRepository) FindByID(ctx context.Context, id uint) (*model.Account, error) {
	var account model.Account
	result := r.db.WithContext(ctx).First(&account, id)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // 見つからない場合は nil を返す
		}
		return nil, fmt.Errorf("failed to find account by id: %w", result.Error)
	}
	
	return &account, nil
}

// 3. FindByEmail: メールアドレスによるアカウントの1件取得
// ログイン処理の際、「入力されたEmailを持つユーザーがデータベースにいるか」を調べるために使用します
func (r *PostgresAccountRepository) FindByEmail(ctx context.Context, email string) (*model.Account, error) {
	var account model.Account
	// .Where("email = ?", email) で条件を指定し、.First で最初の1件を取得します
	result := r.db.WithContext(ctx).Where("email = ?", email).First(&account)
	
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, nil // 該当するEmailのユーザーがいない場合は nil を返す
		}
		return nil, fmt.Errorf("failed to find account by email: %w", result.Error)
	}
	
	return &account, nil
}
