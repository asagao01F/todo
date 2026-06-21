package handler

import (
	"encoding/json"
	"net/http"
	"todo/base/model"
	"context"
	"github.com/oapi-codegen/runtime/types"
)

// AccountUsecaseInterface は Handlerが求めるUsecaseのインターフェース定義
type AccountUsecaseInterface interface {
	RegisterAccount(ctx context.Context, account *model.Account) error
	GetAccountByID(ctx context.Context, id uint) (*model.Account, error)
}

type AccountHandler struct {
	accountUsecase AccountUsecaseInterface
}

func NewAccountHandler(u AccountUsecaseInterface) *AccountHandler {
	return &AccountHandler{accountUsecase: u}
}

// -----------------------------------------------------------------------------
// 1. Register: POST /accounts
// -----------------------------------------------------------------------------
func (h *AccountHandler) RegisterAccount(w http.ResponseWriter, r *http.Request) {
	var req RegisterAccountRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.respondWithError(w, http.StatusBadRequest, "不正なリクエスト形式です", nil)
		return
	}

	// Usecaseの引数に合わせて model.Account を組み立てる
	account := &model.Account{
		Username :     req.Username,
		Email:    types.Email(req.Email),
		Password: req.Password, // パスワードのハッシュ化ロジック等をUsecase側に持たせる場合はこのままでOK
	}

	// ビジネスロジックの実行
	if err := h.accountUsecase.RegisterAccount(r.Context(), account); err != nil {
		// すでに登録済みのエラーなどの場合は 400 Bad Request、それ以外は 500 を返す
		if err.Error() == "email address is already registered" || err.Error() == "email and password cannot be empty" {
			h.respondWithError(w, http.StatusBadRequest, err.Error(), nil)
		} else {
			h.respondWithError(w, http.StatusInternalServerError, "アカウントの登録に失敗しました", []string{err.Error()})
		}
		return
	}

	// 登録成功時は、作成されたデータをレスポンス型にマッピング
	// (GORMでの保存成功時に account.ID が自動で採番されている前提)
	res := AccountResponse{
		Id:    int64(account.ID),
		Email: types.Email(account.Email),
		Username: account.Username,
	}

	h.respondWithJSON(w, http.StatusCreated, res)
}

// -----------------------------------------------------------------------------
// 2. Read (Single): GET /accounts/{id}
// -----------------------------------------------------------------------------
// ※ 第3引数には openapi-generator の仕様通りパース済みの id (int64) が入る想定です
func (h *AccountHandler) GetAccountById(w http.ResponseWriter, r *http.Request, id int64) {
	// ビジネスロジックの実行
	account, err := h.accountUsecase.GetAccountByID(r.Context(), uint(id))
	if err != nil {
		if err.Error() == "account not found" {
			h.respondWithError(w, http.StatusNotFound, "指定されたアカウントが見つかりません", nil)
		} else {
			h.respondWithError(w, http.StatusInternalServerError, "アカウント情報の取得に失敗しました", []string{err.Error()})
		}
		return
	}

	// レスポンス型へのマッピング
	res := AccountResponse{
		Id:    int64(account.ID),
		Email: types.Email(account.Email),
		Username: account.Username,
	}

	h.respondWithJSON(w, http.StatusOK, res)
}

// -----------------------------------------------------------------------------
// 共通ヘルパー関数（JSON返却用）
// -----------------------------------------------------------------------------

func (h *AccountHandler) respondWithJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func (h *AccountHandler) respondWithError(w http.ResponseWriter, status int, message string, details []string) {
	res := ErrorResponse{
		Message: message,
	}
	if details != nil {
		res.Details = &details
	}
	h.respondWithJSON(w, status, res)
}
