package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"
	"todo/api/handler"
	"todo/base/repository"
	"todo/config"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func ConnectDB(cfg *config.PostgresConfig) (*sql.DB, error) {
	// 接続文字列（DSN）の作成
	dsn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName, cfg.SSLMode,
	)

	// DBオープン（この時点ではまだ接続テストは走らない）
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	// コネクションプールの設定
	db.SetMaxOpenConns(cfg.MaxOpenConns)
	db.SetMaxIdleConns(cfg.MaxIdleConns)
	db.SetConnMaxLifetime(5 * time.Minute)

	// 実際に Ping を打って疎通確認
	if err := db.Ping(); err != nil {
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

	// 2. PostgreSQLデータベース接続
	db, err := ConnectDB(&cfg.Database.Postgres)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// 3. リポジトリ・ハンドラーの初期化 (依存関係の注入: DI)
	todoRepo := repository.NewPostgresTodoRepository(db)

	todoHandler := handler.NewTodoHandler(todoRepo) // 本来はここにusecaseが入ります
	healthHandler := handler.NewHealthHandler(db)

	// 4. ルーターの作成（routes.goから呼び出し）
	router := NewRouter(todoHandler, healthHandler)

	// 5. HTTPサーバーの起動設定
	serverAddr := fmt.Sprintf(":%d", cfg.API.Port)
	server := &http.Server{
		Addr:         serverAddr,
		Handler:      router, // 作成したルーターをセット
		ReadTimeout:  time.Duration(cfg.API.TimeoutSeconds) * time.Second,
		WriteTimeout: time.Duration(cfg.API.TimeoutSeconds) * time.Second,
	}

	log.Printf("[%s] API Server starting on %s...", cfg.App.Env, serverAddr)
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}
}
