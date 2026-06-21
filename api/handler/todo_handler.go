package handler

import (
	"encoding/json"
	"net/http"
	"todo/base/usecase"
)

type TodoHandler struct {
	todoUsecase usecase.TodoUsecaseInterface
}

func NewTodoHandler(u usecase.TodoUsecaseInterface) *TodoHandler {
	return &TodoHandler{todoUsecase: u}
}

// -----------------------------------------------------------------------------
// 1. Create: POST /todos
// -----------------------------------------------------------------------------
func (h *TodoHandler) CreateTodo(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "不正なリクエスト形式です", nil)
		return
	}

	var desc string
	if req.Description != nil {
		desc = *req.Description
	}

	// 【変更】ビジネスロジックとDB保存はUsecaseに丸投げする
	todo, err := h.todoUsecase.CreateTodo(r.Context(), req.Title, desc, req.AccountId)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "データの保存に失敗しました", []string{err.Error()})
		return
	}

	// Usecaseから戻ってきた model.Todo を OpenAPIのレスポンス型に変換
	res := TodoResponse{
		Id:          int64(todo.ID),
		Title:       todo.Title,
		Description: req.Description,
		Completed:   todo.IsCompleted,
		CreatedAt:   todo.CreatedAt,
	}

	h.respondWithJSON(w, http.StatusCreated, res)
}

// -----------------------------------------------------------------------------
// 2. Read (Single): GET /todos/{id}
// -----------------------------------------------------------------------------
func (h *TodoHandler) GetTodoById(w http.ResponseWriter, r *http.Request, id int64) {
	// 【変更】スタブを排し、Usecaseから本物のデータを取得する
	todo, err := h.todoUsecase.GetTodoByID(r.Context(), uint(id))
	if err != nil {
		// 簡易的に404エラー。必要に応じてエラー内容でステータスコードを分岐してください
		h.respondWithError(w, http.StatusNotFound, "指定されたタスクが見つかりません", []string{err.Error()})
		return
	}

	// OpenAPIレスポンス型へのマッピング
	res := TodoResponse{
		Id:          int64(todo.ID),
		Title:       todo.Title,
		Description: &todo.Description,
		Completed:   todo.IsCompleted,
		CreatedAt:   todo.CreatedAt,
	}

	h.respondWithJSON(w, http.StatusOK, res)
}

// -----------------------------------------------------------------------------
// 共通ヘルパー関数
// -----------------------------------------------------------------------------

func (h *TodoHandler) respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *TodoHandler) respondWithError(w http.ResponseWriter, status int, message string, details []string) {
	res := ErrorResponse{
		Message: message,
	}
	if details != nil {
		res.Details = &details
	}
	h.respondWithJSON(w, status, res)
}
