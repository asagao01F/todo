package main

import (
	"net/http"
	"todo/api/handler"
)

// NewRouter は、すべてのAPIルートを定義したマルチプレクサ（ルーター）を返します
func NewRouter(todoHandler *handler.TodoHandler, healthHandler *handler.HealthHandler) *http.ServeMux {
	mux := http.NewServeMux()

	// -----------------------------------------------------------------------------
	// OpenAPIが生成したルーティングを一括登録する
	// -----------------------------------------------------------------------------
	
	// 各ハンドラーを結合した、OpenAPI用の総合サーバーを作成
	// (後述の複合構造体をここで使います)
	apiServer := &CombinedServer{
		TodoHandler:   todoHandler,
		HealthHandler: healthHandler,
	}

	// 自動生成された HandlerFromMux 関数に mux と apiServer を渡すだけで、
	// YAMLに書かれたすべてのルート (GET /health, POST /todos, GET /todos/{id}) が
	// 自動的に Go 1.22 形式で mux に登録されます。
	handler.HandlerFromMux(apiServer, mux)

	return mux
}

// -----------------------------------------------------------------------------
// 補助構造体: 2つのハンドラーを1つにまとめる
// -----------------------------------------------------------------------------

// CombinedServer は、自動生成された `handler.ServerInterface` を満たすための複合構造体です。
// 分割して実装した各ハンドラーへ処理を委譲（Proxy）します。
type CombinedServer struct {
	*handler.TodoHandler
	*handler.HealthHandler
}

// 1. ヘルスチェックの委譲
func (c *CombinedServer) CheckHealth(w http.ResponseWriter, r *http.Request) {
	c.HealthHandler.CheckHealth(w, r)
}

// 2. TODO作成の委譲
func (c *CombinedServer) CreateTodo(w http.ResponseWriter, r *http.Request) {
	c.TodoHandler.CreateTodo(w, r)
}

// 3. TODO取得の委譲
func (c *CombinedServer) GetTodoById(w http.ResponseWriter, r *http.Request, id int64) {
	c.TodoHandler.GetTodoById(w, r, id)
}
