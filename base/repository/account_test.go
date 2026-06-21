package repository

import (
	"context"
	"fmt"
	"regexp"
	"testing"
	"time"

	"todo/base/model"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GORMとsqlmockを初期化するヘルパー関数
func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open a stub database connection: %v", err)
	}

	dialer := postgres.New(postgres.Config{
		Conn:                 db,
		PreferSimpleProtocol: true,
	})

	gormDB, err := gorm.Open(dialer, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	return gormDB, mock
}

// 1. Create のテスト
func TestPostgresAccountRepository_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("成功: アカウントが正常に作成されること", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewPostgresAccountRepository(db)

		account := &model.Account{
			Username: "testuser",
			Email:    "test@example.com",
			Password: "hashed_password_here",
		}

		// テーブル名が "account" になり、全カラムが対象になります
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "account" ("username","email","password","created_at","updated_at") VALUES ($1,$2,$3,$4,$5) RETURNING "id"`)).
			WithArgs(account.Username, account.Email, account.Password, sqlmock.AnyArg(), sqlmock.AnyArg()).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(uint(1)))
		mock.ExpectCommit()

		err := repo.Create(ctx, account)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if account.ID != 1 {
			t.Errorf("expected account ID to be 1, got %d", account.ID)
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("失敗: DBエラー発生時", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewPostgresAccountRepository(db)

		account := &model.Account{
			Username: "erroruser",
			Email:    "error@example.com",
			Password: "password",
		}

		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "account"`)).
			WillReturnError(fmt.Errorf("db error"))
		mock.ExpectRollback()

		err := repo.Create(ctx, account)

		if err == nil {
			t.Error("expected error, got nil")
		}
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}

// 2. FindByID のテスト
func TestPostgresAccountRepository_FindByID(t *testing.T) {
	ctx := context.Background()

	t.Run("成功: レコードが存在する場合", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewPostgresAccountRepository(db)

		targetID := uint(1)
		now := time.Now()

		// テーブル名が "account" に変更され、ORDER BY も "account"."id" になります
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "account" WHERE "account"."id" = $1 ORDER BY "account"."id" LIMIT $2`)).
			WithArgs(targetID, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}).
				AddRow(targetID, "testuser", "test@example.com", "hashed_password", now, now))

		res, err := repo.FindByID(ctx, targetID)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if res == nil || res.ID != targetID || res.Username != "testuser" {
			t.Errorf("unexpected result: %+v", res)
		}
	})

	t.Run("成功: レコードが存在しない場合 (nilを返す)", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewPostgresAccountRepository(db)

		targetID := uint(999)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "account"`)).
			WithArgs(targetID, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		res, err := repo.FindByID(ctx, targetID)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if res != nil {
			t.Errorf("expected nil, got %+v", res)
		}
	})
}

// 3. FindByEmail のテスト
func TestPostgresAccountRepository_FindByEmail(t *testing.T) {
	ctx := context.Background()

	t.Run("成功: メールアドレスに一致するレコードが存在する場合", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewPostgresAccountRepository(db)

		targetEmail := "user@example.com"
		now := time.Now()

		// WHERE句の条件と、First()による ORDER BY "account"."id" LIMIT 1
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "account" WHERE email = $1 ORDER BY "account"."id" LIMIT $2`)).
			WithArgs(targetEmail, 1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "username", "email", "password", "created_at", "updated_at"}).
				AddRow(uint(10), "user", targetEmail, "hashed_password", now, now))

		res, err := repo.FindByEmail(ctx, targetEmail)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if res == nil || res.Email != targetEmail || res.ID != 10 {
			t.Errorf("unexpected result: %+v", res)
		}
	})

	t.Run("成功: メールアドレスに一致するレコードがない場合 (nilを返す)", func(t *testing.T) {
		db, mock := setupMockDB(t)
		repo := NewPostgresAccountRepository(db)

		targetEmail := "notfound@example.com"

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "account" WHERE email = $1`)).
			WithArgs(targetEmail, 1).
			WillReturnError(gorm.ErrRecordNotFound)

		res, err := repo.FindByEmail(ctx, targetEmail)

		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if res != nil {
			t.Errorf("expected nil, got %+v", res)
		}
	})
}
