package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	"todo/base/model"
	"todo/base/usecase" // 今後作成するusecaseパッケージ
)

// TodoHandler は HTTPリクエストを受け付け、Usecaseを呼び出す構造体です
type TodoHandler struct {
	todoUsecase usecase.TodoRepository // 本来はusecaseのインターフェースを指定
}

// NewTodoHandler はハンドラーのインスタンスを生成します（DI用）
func NewTodoHandler(u usecase.TodoRepository) *TodoHandler {
	return &TodoHandler{todoUsecase: u}
}

// -----------------------------------------------------------------------------
// 1. Create: POST /todos
// -----------------------------------------------------------------------------
func (h *TodoHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateTodoRequest
	// リクエストボディのJSONを構造体にパース
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "不正なリクエスト形式です", nil)
		return
	}

	// 本来はここでバリデーション（req.Titleの空チェックなど）を行う

	// DTOからビジネスロジック用のModel（内部ドメイン）へ変換
	todo := &model.Todo{
		Title:       req.Title,
		Description: req.Description,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Usecase（Transaction層）の呼び出し
	if err := h.todoUsecase.Create(r.Context(), todo); err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "データの保存に失敗しました", []string{err.Error()})
		return
	}

	// レスポンス用のDTOに変換して返却
	res := TodoResponse{
		ID:          todo.ID,
		Title:       todo.Title,
		Description: todo.Description,
		Completed:   todo.Completed,
		CreatedAt:   todo.CreatedAt,
	}

	h.respondWithJSON(w, http.StatusCreated, res)
}

// -----------------------------------------------------------------------------
// 2. Read (Single): GET /todos/{id}
// -----------------------------------------------------------------------------
func (h *TodoHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	// ルーターで定義した {id} の値を文字列で取得
	idStr := r.PathValue("id") 
	
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		h.respondWithError(w, http.StatusBadRequest, "無効なID形式です", nil)
		return
	}

	// 本来はここで usecase からデータを取得
	// todo, err := h.todoUsecase.GetByID(r.Context(), id)

	// 一旦スタブ（仮）レスポンスを返却
	res := TodoResponse{
		ID:          id,
		Title:       "仮のタスクタイトル",
		Description: "これはパスパラメータから取得したID: " + idStr + " のタスクです",
		Completed:   false,
		CreatedAt:   time.Now(),
	}

	h.respondWithJSON(w, http.StatusOK, res)
}

// -----------------------------------------------------------------------------
// 共通ヘルパー関数（JSON返却用）
// -----------------------------------------------------------------------------

// respondWithJSON は一貫したフォーマットで正常系JSONを返却するためのヘルパーです
func (h *TodoHandler) respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

// respondWithError は一貫したフォーマットで異常系JSON（エラー）を返却するためのヘルパーです
func (h *TodoHandler) respondWithError(w http.ResponseWriter, status int, message string, details []string) {
	res := ErrorResponse{
		Message: message,
		Details: details,
	}
	h.respondWithJSON(w, status, res)
}