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
	todo, err := h.todoUsecase.CreateTodo(r.Context(), req.Title, desc, req.AccountId, *req.DueDate)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "データの保存に失敗しました", []string{err.Error()})
		return
	}

	// Usecaseから戻ってきた model.Todo を OpenAPIのレスポンス型に変換
	res := TodoResponse{
		Id:          int64(todo.ID),
		Title:       todo.Title,
		DueDate:     todo.DueDate,
		Description: req.Description,
		Completed:   todo.IsCompleted,
		CreatedAt:   todo.CreatedAt,
	}

	h.respondWithJSON(w, http.StatusCreated, res)
}

// -----------------------------------------------------------------------------
// 2. Read (List): GET /todos  ★新規追加
// -----------------------------------------------------------------------------
// ※ 第3引数の params には oapi-codegen によりクエリパラメータ（GetTodoListParams）が自動パースされて入ります
func (h *TodoHandler) GetTodoList(w http.ResponseWriter, r *http.Request, params GetTodoListParams) {
	// クエリパラメータからオプショナルな accountId を安全に抽出
	var accountId *int64
	if params.AccountId != nil {
		accountId = params.AccountId
	}

	// 【紐付け】Usecase 側に FetchTodoList(ctx, accountId) のような一覧取得メソッドがある想定です
	todos, err := h.todoUsecase.FetchTodoList(r.Context(), accountId)
	if err != nil {
		h.respondWithError(w, http.StatusInternalServerError, "TODO一覧の取得に失敗しました", []string{err.Error()})
		return
	}

	// スライスが空だった場合に JSON で null ではなく [] を返却するための初期化
	res := make([]TodoResponse, 0, len(todos))

	// model.Todo のスライスを OpenAPI のレスポンス配列にマッピング
	for _, todo := range todos {
		// ループ内でのアドレス固定化と nil チェック対策
		todoItem := todo
		var descPtr *string
		if todoItem.Description != "" {
			descPtr = &todoItem.Description
		}

		res = append(res, TodoResponse{
			Id:          int64(todoItem.ID),
			Title:       todoItem.Title,
			Description: descPtr,
			Completed:   todoItem.IsCompleted,
			CreatedAt:   todoItem.CreatedAt,
			DueDate:     todoItem.DueDate,
		})
	}

	h.respondWithJSON(w, http.StatusOK, res)
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
