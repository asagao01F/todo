package main

import (
    "net/http"
    "todo/api/handler"
    "github.com/go-chi/chi/v5"
)

func NewRouter(todoHandler *handler.TodoHandler, healthHandler *handler.HealthHandler, accountHandler *handler.AccountHandler) chi.Router {
    // 【変更】Go標準の ServeMux ではなく、chi のルーターを生成
    r := chi.NewRouter()

    apiServer := &CombinedServer{
        TodoHandler:   todoHandler,
        HealthHandler: healthHandler,
        AccountHandler: accountHandler,
    }

    handler.HandlerFromMux(apiServer, r)

    return r
}

// -----------------------------------------------------------------------------
// 補助構造体: 2つのハンドラーを1つにまとめる
// -----------------------------------------------------------------------------

// CombinedServer は、自動生成された `handler.ServerInterface` を満たすための複合構造体です。
type CombinedServer struct {
    *handler.AccountHandler
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

// -----------------------------------------------------------------------------
// 4. アカウント登録の委譲
func (c *CombinedServer) RegisterAccount(w http.ResponseWriter, r *http.Request) {
    c.AccountHandler.RegisterAccount(w, r)
}

// 5. アカウント取得の委譲
func (c *CombinedServer) GetAccountById(w http.ResponseWriter, r *http.Request, id int64) {
    c.AccountHandler.GetAccountById(w, r, id)
}
