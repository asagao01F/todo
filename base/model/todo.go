package model

import "time"

// Todo はデータベースの todos テーブルに対応する構造体です
type Todo struct {
	ID          int64     `json:"id"`          // プライマリキー
	AccountId   int64     `json:"accountId"`   // accountId
	Title       string    `json:"title"`       // タスクのタイトル
	Description string    `json:"description"` // タスクの詳細説明
	Completed   bool      `json:"completed"`   // 完了フラグ
	CreatedAt   time.Time `json:"created_at"`  // 作成日時
	UpdatedAt   time.Time `json:"updated_at"`  // 更新日時
}