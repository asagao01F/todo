package repository

import (
    "log"

    "github.com/DATA-DOG/go-sqlmock"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

// SetupTestDB はテスト用のモックDBと、クエリを検証するためのmockオブジェクトを返します
func SetupTestDB() (*gorm.DB, sqlmock.Sqlmock) {
    // 1. sqlmockで疑似的な sql.DB を作成
    db, mock, err := sqlmock.New()
    if err != nil {
        log.Fatalf("failed to open a stub database connection: %v", err)
    }

    // 2. GORMに疑似 sql.DB を流し込んで初期化
    dialector := postgres.New(postgres.Config{
        Conn: db,
    })

    gormDB, err := gorm.Open(dialector, &gorm.Config{})
    if err != nil {
        log.Fatalf("failed to open gorm DB: %v", err)
    }

    return gormDB, mock
}
