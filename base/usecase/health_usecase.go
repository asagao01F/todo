package usecase

import (
	"context"
	"gorm.io/gorm"
)

// HealthStatus はUsecaseが判定した結果を格納する構造体
type HealthStatus struct {
	IsHealthy bool
	DBStatus  string
}

type HealthUsecase struct {
	db *gorm.DB
}

func NewHealthUsecase(db *gorm.DB) *HealthUsecase {
	return &HealthUsecase{db: db}
}

// DiagnoseHealth はデータベースの接続状態をチェックして結果を返します
func (u *HealthUsecase) DiagnoseHealth(ctx context.Context) HealthStatus {
	// 1. そもそもDB初期化に失敗している場合
	if u.db == nil {
		return HealthStatus{
			IsHealthy: false,
			DBStatus:  "NOT_INITIALIZED",
		}
	}

	// 2. GORMから内部の *sql.DB を取得
	sqlDB, err := u.db.DB()
	if err != nil {
		return HealthStatus{
			IsHealthy: false,
			DBStatus:  "ERROR",
		}
	}

	// 3. Pingを打って疎通確認
	if err := sqlDB.PingContext(ctx); err != nil {
		return HealthStatus{
			IsHealthy: false,
			DBStatus:  "DISCONNECTED",
		}
	}

	// すべて正常
	return HealthStatus{
		IsHealthy: true,
		DBStatus:  "CONNECTED",
	}
}
