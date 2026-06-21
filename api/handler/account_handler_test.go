package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"todo/base/model"

	openapi_types "github.com/oapi-codegen/runtime/types"
)

// 1. Usecaseのインターフェースを満たすテスト用のモック構造体を定義
type mockAccountUsecase struct {
	fakeRegisterAccountFn func(ctx context.Context, account *model.Account) error
	fakeGetAccountByIDFn  func(ctx context.Context, id uint) (*model.Account, error)
}

func (m *mockAccountUsecase) RegisterAccount(ctx context.Context, account *model.Account) error {
	return m.fakeRegisterAccountFn(ctx, account)
}

func (m *mockAccountUsecase) GetAccountByID(ctx context.Context, id uint) (*model.Account, error) {
	return m.fakeGetAccountByIDFn(ctx, id)
}

// -----------------------------------------------------------------------------
// RegisterAccount (POST /accounts) のテスト
// -----------------------------------------------------------------------------
func TestAccountHandler_RegisterAccount(t *testing.T) {
	t.Run("成功: 正しいリクエストで201 Createdと作成データが返ること", func(t *testing.T) {
		mockUc := &mockAccountUsecase{
			fakeRegisterAccountFn: func(ctx context.Context, account *model.Account) error {
				// DB保存成功をシミュレートしてIDを付与
				account.ID = 1
				return nil
			},
		}
		handler := NewAccountHandler(mockUc)

		reqBody := RegisterAccountRequest{
			Email:    "test@example.com",
			Password: "password123",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(bodyBytes))
		rec := httptest.NewRecorder()

		handler.RegisterAccount(rec, req)

		if rec.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", rec.Code)
		}

		var res AccountResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if res.Id != 1 || res.Email != openapi_types.Email("test@example.com") {
			t.Errorf("unexpected response: %+v", res)
		}
	})

	t.Run("失敗: 重複するEmailの場合に400 Bad Requestが返ること", func(t *testing.T) {
		mockUc := &mockAccountUsecase{
			fakeRegisterAccountFn: func(ctx context.Context, account *model.Account) error {
				return errors.New("email address is already registered")
			},
		}
		handler := NewAccountHandler(mockUc)

		reqBody := RegisterAccountRequest{
			Email:    "duplicate@example.com",
			Password: "password123",
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewBuffer(bodyBytes))
		rec := httptest.NewRecorder()

		handler.RegisterAccount(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}

		var res ErrorResponse
		_ = json.Unmarshal(rec.Body.Bytes(), &res)
		if res.Message != "email address is already registered" {
			t.Errorf("unexpected error message: %s", res.Message)
		}
	})

	t.Run("失敗: 不正なJSON形式の場合に400 Bad Requestが返ること", func(t *testing.T) {
		mockUc := &mockAccountUsecase{}
		handler := NewAccountHandler(mockUc)

		req := httptest.NewRequest(http.MethodPost, "/accounts", bytes.NewBufferString("{invalid-json"))
		rec := httptest.NewRecorder()

		handler.RegisterAccount(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}
	})
}

// -----------------------------------------------------------------------------
// GetAccountById (GET /accounts/{id}) のテスト
// -----------------------------------------------------------------------------
func TestAccountHandler_GetAccountById(t *testing.T) {
	t.Run("成功: 存在するアカウントIDで200 OKと該当データが返ること", func(t *testing.T) {
		targetID := uint(10)
		mockUc := &mockAccountUsecase{
			fakeGetAccountByIDFn: func(ctx context.Context, id uint) (*model.Account, error) {
				return &model.Account{
					ID:    id,
					Email: "find@example.com",
				}, nil
			},
		}
		handler := NewAccountHandler(mockUc)

		req := httptest.NewRequest(http.MethodGet, "/accounts/10", nil)
		rec := httptest.NewRecorder()

		handler.GetAccountById(rec, req, int64(targetID))

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rec.Code)
		}

		var res AccountResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if res.Id != int64(targetID) || res.Email != openapi_types.Email("find@example.com") {
			t.Errorf("unexpected response: %+v", res)
		}
	})

	t.Run("失敗: アカウントが存在しない場合に404 Not Foundが返ること", func(t *testing.T) {
		mockUc := &mockAccountUsecase{
			fakeGetAccountByIDFn: func(ctx context.Context, id uint) (*model.Account, error) {
				return nil, errors.New("account not found")
			},
		}
		handler := NewAccountHandler(mockUc)

		req := httptest.NewRequest(http.MethodGet, "/accounts/999", nil)
		rec := httptest.NewRecorder()

		handler.GetAccountById(rec, req, 999)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}

		var res ErrorResponse
		_ = json.Unmarshal(rec.Body.Bytes(), &res)
		if res.Message != "指定されたアカウントが見つかりません" {
			t.Errorf("unexpected error message: %s", res.Message)
		}
	})
}
