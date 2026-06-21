package main

import (
    "fmt"
    "log"
    "net/http"
    "time"
    "todo/api/handler"
    "todo/base/repository"
    "todo/base/usecase"
    "todo/config"

    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

// 【変更】戻り値を *gorm.DB に変更
func ConnectDB(cfg *config.PostgresConfig) (*gorm.DB, error) {
    // 接続文字列（DSN）の作成
    dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
        cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
    )

    // 【変更】gorm.Open を使用して接続
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to connect database: %w", err)
    }

    // 【変更】コネクションプールの設定を行うため、内部の *sql.DB を取得
    sqlDB, err := db.DB()
    if err != nil {
        return nil, fmt.Errorf("failed to get sql.DB: %w", err)
    }

    // コネクションプールの設定（sqlDBに対して行う）
    sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
    sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
    sqlDB.SetConnMaxLifetime(5 * time.Minute)

    // 疎通確認
    if err := sqlDB.Ping(); err != nil {
        return nil, fmt.Errorf("database ping failed: %w", err)
    }

    return db, nil
}

func main() {
    // 1. 設定ファイルの読み込み
    cfg, err := config.LoadConfig()
    if err != nil {
        log.Fatalf("Failed to load config: %v", err)
    }

    // 2. PostgreSQLデータベース接続（戻り値が *gorm.DB になっています）
    db, err := ConnectDB(&cfg.Database.Postgres)
    if err != nil {
        log.Fatalf("Failed to connect to database: %v", err)
    }

    // 3. リポジトリ・ハンドラーの初期化
    todoRepo := repository.NewPostgresTodoRepository(db) // エラーが解消されます
    healthUsecase := usecase.NewHealthUsecase(db) // 修正: HealthUsecaseの初期化も追加

    todoHandler := handler.NewTodoHandler(todoRepo)
    healthHandler := handler.NewHealthHandler(healthUsecase) // 修正: healthHandler を作成

    // 4. ルーターの作成
    router := NewRouter(todoHandler, healthHandler)

    // 5. HTTPサーバーの起動設定
    serverAddr := fmt.Sprintf(":%d", cfg.API.Port)
    server := &http.Server{
        Addr:         serverAddr,
        Handler:      router,
        ReadTimeout:  time.Duration(cfg.API.TimeoutSeconds) * time.Second,
        WriteTimeout: time.Duration(cfg.API.TimeoutSeconds) * time.Second,
    }

    log.Printf("[%s] API Server starting on %s...", cfg.App.Env, serverAddr)
    if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
        log.Fatalf("Server failed to start: %v", err)
    }
}
