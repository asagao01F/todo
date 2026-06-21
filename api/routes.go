package main

import (
    "net/http"
    "todo/api/handler"
    "github.com/go-chi/chi/v5"
)

func NewRouter(todoHandler *handler.TodoHandler, healthHandler *handler.HealthHandler) chi.Router {
    // 【変更】Go標準の ServeMux ではなく、chi のルーターを生成
    r := chi.NewRouter()

    apiServer := &CombinedServer{
        TodoHandler:   todoHandler,
        HealthHandler: healthHandler,
    }

    handler.HandlerFromMux(apiServer, r)

    return r
}

// -----------------------------------------------------------------------------
// 補助構造体: 2つのハンドラーを1つにまとめる
// -----------------------------------------------------------------------------

// CombinedServer は、自動生成された `handler.ServerInterface` を満たすための複合構造体です。
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
