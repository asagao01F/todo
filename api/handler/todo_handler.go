package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"todo/base/model"
	"todo/base/usecase"
)

type TodoHandler struct {
	todoUsecase usecase.TodoRepository
}

func NewTodoHandler(u usecase.TodoRepository) *TodoHandler {
	return &TodoHandler{todoUsecase: u}
}

// -----------------------------------------------------------------------------
// 1. Create: POST /todos
// -----------------------------------------------------------------------------
// CreateTodo は自動生成された ServerInterface のメンバです
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	// CreateTodoRequest 型は openapi.gen.go で自動生成されたものを使用
	var req CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "不正なリクエスト形式です", nil)
		return
	}

	// ポインタ型（オプショナル）の Description を安全に扱うための処理
	var desc string
	if req.Description != nil {
		desc = *req.Description
	}

	todo := &model.Todo{
		Title:       req.Title,
		Description: desc,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := h.todoUsecase.Create(r.Context(), todo); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "データの保存に失敗しました", []string{err.Error()})
		return
	}

	// TodoResponse 型は openapi.gen.go で自動生成されたものを使用
	res := TodoResponse{
		Id:          &todo.ID, // OpenAPIの定義に合わせ、ポインタやフィールド名のケースが自動調整されます
		Title:       todo.Title,
		Description: req.Description, // ポインタをそのまま渡せます
		Completed:   todo.Completed,
		CreatedAt:   todo.CreatedAt,
	}

	h.respondWithJSON(w, http.StatusCreated, res)
}

// -----------------------------------------------------------------------------
// 2. Read (Single): GET /todos/{id}
// -----------------------------------------------------------------------------
// GetTodoById は自動生成された ServerInterface のメンバです
// ※ 第3引数に最初からパース済みの id (int64) が入ってくるため、strconvの処理が不要になります！
func (h *TodoHandler) GetTodoById(w http.ResponseWriter, r *http.Request, id int64) {
	idStr := strconv.FormatInt(id, 10)

	// 一旦スタブ（仮）レスポンスを返却
	desc := "これはパスパラメータから取得したID: " + idStr + " のタスクです"
	res := TodoResponse{
		Id:          &id,
		Title:       "仮のタスクタイトル",
		Description: &desc,
		Completed:   false,
		CreatedAt:   time.Now(),
	}

	h.respondWithJSON(w, http.StatusOK, res)
}

// -----------------------------------------------------------------------------
// 共通ヘルパー関数（JSON返却用）
// -----------------------------------------------------------------------------

func (h *TodoHandler) respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *TodoHandler) respondWithError(w http.ResponseWriter, status int, message string, details []string) {
	// ErrorResponse も自動生成されたものを使用
	res := ErrorResponse{
		Message: message,
	}
	if details != nil {
		res.Details = &details
	}
	h.respondWithJSON(w, status, res)
}
