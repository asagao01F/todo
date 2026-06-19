package main

import (
	"net/http"
	"todo/api/handler"
)

// NewRouter は、すべてのAPIルートを定義したマルチプレクサ（ルーター）を返します
func NewRouter(todoHandler *handler.TodoHandler, healthHandler *handler.HealthHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// -----------------------------------------------------------------------------
	// ヘルスチェック
	// -----------------------------------------------------------------------------
	mux.HandleFunc("GET /health", healthHandler.Check)

	// -----------------------------------------------------------------------------
	// TODO API (CRUD)
	// -----------------------------------------------------------------------------
	// Go 1.22からは "METHOD /path" の形式でHTTPメソッドを制限できるようになりました
	mux.HandleFunc("POST /todos", todoHandler.Create)      // C: 作成
	mux.HandleFunc("GET /todos/{id}", todoHandler.GetByID) // R: 1件取得 ({id}でパスパラメータを受け取る)
	// mux.HandleFunc("GET /todos", todoHandler.GetAll)     // R: 一覧取得（必要に応じて実装）
	// mux.HandleFunc("PUT /todos/{id}", todoHandler.Update)  // U: 更新（必要に応じて実装）
	// mux.HandleFunc("DELETE /todos/{id}", todoHandler.Delete) // D: 削除（必要に応じて実装）

	return mux
}
