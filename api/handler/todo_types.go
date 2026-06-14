package handler

import "time"

// -----------------------------------------------------------------------------
// 1. Create (作成) 用の構造体
// -----------------------------------------------------------------------------

// CreateTodoRequest は、POST /todos で受け取るリクエストボディです
type CreateTodoRequest struct {
	Title       string `json:"title" validate:"required,max=100"` // 必須、最大100文字
	Description string `json:"description" validate:"max=500"`    // 任意、最大500文字
}

// -----------------------------------------------------------------------------
// 2. Read (取得) 用の構造体
// -----------------------------------------------------------------------------

// TodoResponse は、APIがクライアントに返す標準的なTodoのデータ形式です
// (内部的な型をそのまま出さず、API仕様に合わせたJSONタグや型を定義します)
type TodoResponse struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Completed   bool      `json:"completed"`
	CreatedAt   time.Time `json:"created_at"`
}

// TodoListResponse は、GET /todos などの一覧取得で返すレスポンスです
type TodoListResponse struct {
	Todos []*TodoResponse `json:"todos"`
	Total int             `json:"total"`
}

// -----------------------------------------------------------------------------
// 3. Update (更新) 用の構造体
// -----------------------------------------------------------------------------

// UpdateTodoRequest は、PUT /todos/:id で受け取るリクエストボディです
type UpdateTodoRequest struct {
	Title       string `json:"title" validate:"required,max=100"`
	Description string `json:"description" validate:"max=500"`
	Completed   bool   `json:"completed"` // 完了ステータスの変更
}

// -----------------------------------------------------------------------------
// 4. Delete (削除) や共通エラーの構造体
// -----------------------------------------------------------------------------

// ErrorResponse は、バリデーションエラーやシステムエラーが発生した際に一貫した形式で返すための構造体です
type ErrorResponse struct {
	Message string   `json:"message"`           // エラーの概要文
	Details []string `json:"details,omitempty"` // 詳細なエラー内容（バリデーションの項目など）
}