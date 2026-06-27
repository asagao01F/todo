package usecase

import (
	"context"
	"errors"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// GORMとsqlmockを初期化するヘルパー関数
func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	// PingContextを検証するため MonitorPingsOption を true にする
	db, mock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
	if err != nil {
		t.Fatalf("failed to open a stub database connection: %v", err)
	}

	// GORM初期化（gorm.Open）時に内部で走る自動Pingをあらかじめ期待値に入れて消費させる
	mock.ExpectPing()

	dialer := postgres.New(postgres.Config{
		Conn: db,
	})

	gormDB, err := gorm.Open(dialer, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm db: %v", err)
	}

	return gormDB, mock
}

func TestHealthUsecase_DiagnoseHealth(t *testing.T) {
	ctx := context.Background()

	t.Run("成功: DB接続が正常な場合、CONNECTED を返すこと", func(t *testing.T) {
		db, mock := setupMockDB(t)
		uc := NewHealthUsecase(db)

		// DiagnoseHealth 内で実行される sqlDB.PingContext(ctx) に対する期待値
		mock.ExpectPing()

		res := uc.DiagnoseHealth(ctx)

		// アサーション
		if !res.IsHealthy {
			t.Errorf("expected IsHealthy to be true, got false")
		}
		if res.DBStatus != "CONNECTED" {
			t.Errorf("expected DBStatus to be 'CONNECTED', got '%s'", res.DBStatus)
		}

		// すべての期待値が満たされたか検証
		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})

	t.Run("失敗: そもそもDBインスタンスがnilの場合、NOT_INITIALIZED を返すこと", func(t *testing.T) {
		// 意図的に nil を渡して Usecase を初期化
		uc := NewHealthUsecase(nil)

		res := uc.DiagnoseHealth(ctx)

		// アサーション
		if res.IsHealthy {
			t.Errorf("expected IsHealthy to be false, got true")
		}
		if res.DBStatus != "NOT_INITIALIZED" {
			t.Errorf("expected DBStatus to be 'NOT_INITIALIZED', got '%s'", res.DBStatus)
		}
	})

	t.Run("失敗: Pingがエラーを返す場合、DISCONNECTED を返すこと", func(t *testing.T) {
		db, mock := setupMockDB(t)
		uc := NewHealthUsecase(db)

		// PingContext がエラーを返すようにモックを設定
		mock.ExpectPing().WillReturnError(errors.New("network connection timeout"))

		res := uc.DiagnoseHealth(ctx)

		// アサーション
		if res.IsHealthy {
			t.Errorf("expected IsHealthy to be false, got true")
		}
		if res.DBStatus != "DISCONNECTED" {
			t.Errorf("expected DBStatus to be 'DISCONNECTED', got '%s'", res.DBStatus)
		}

		if err := mock.ExpectationsWereMet(); err != nil {
			t.Errorf("there were unfulfilled expectations: %s", err)
		}
	})
}
