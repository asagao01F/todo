package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"todo/base/model"
)

// 1. Usecaseのインターフェースを満たすテスト用のモック構造体を定義
type mockTodoUsecase struct {
	// テストケースごとに期待する戻り値を外からコントロールするためのフィールド
	fakeCreateTodoFn func(ctx context.Context, title string, description string, accountId *int64, dueDate time.Time) (*model.Todo, error)
	fakeGetTodoByIDFn func(ctx context.Context, id uint) (*model.Todo, error)
}

func (m *mockTodoUsecase) CreateTodo(ctx context.Context, title string, description string, accountId *int64, dueDate time.Time) (*model.Todo, error) {
	return m.fakeCreateTodoFn(ctx, title, description, accountId, dueDate)
}

func (m *mockTodoUsecase) GetTodoByID(ctx context.Context, id uint) (*model.Todo, error) {
	return m.fakeGetTodoByIDFn(ctx, id)
}

// -----------------------------------------------------------------------------
// CreateTodo (POST /todos) のテスト
// -----------------------------------------------------------------------------
func TestTodoHandler_CreateTodo(t *testing.T) {
	t.Run("成功: 正しいリクエストが送られた場合、201 Createdと作成データを返すこと", func(t *testing.T) {
		// モックの振る舞いを定義
		mockUc := &mockTodoUsecase{
			fakeCreateTodoFn: func(ctx context.Context, title string, description string, accountId *int64, dueDate time.Time) (*model.Todo, error) {
				return &model.Todo{
					ID:          123,
					AccountId:   accountId,
					Title:       title,
					Description: description,
					IsCompleted: false,
					CreatedAt:   time.Now(),
					DueDate:     dueDate,
				}, nil
			},
		}
		handler := NewTodoHandler(mockUc)

		// 擬似リクエストボディの作成
		desc := "テスト詳細"
		now := time.Now()
		reqBody := CreateTodoRequest{
			Title:       "テストタスク",
			Description: &desc,
			AccountId:   &[]int64{1}[0],
			DueDate: &now,
		}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(bodyBytes))
		rec := httptest.NewRecorder()

		handler.CreateTodo(rec, req)

		// ステータスコードの検証
		if rec.Code != http.StatusCreated {
			t.Errorf("expected status 201, got %d", rec.Code)
		}

		// レスポンスボディの検証
		var res TodoResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if res.Id != 123 || res.Title != "テストタスク" || *res.Description != "テスト詳細" || res.Completed != false {
			t.Errorf("unexpected response: %+v", res)
		}
	})

	t.Run("失敗: 不正なJSON形式データが送られた場合、400 Bad Requestを返すこと", func(t *testing.T) {
		mockUc := &mockTodoUsecase{} // 内部の関数は呼ばれないので空でOK
		handler := NewTodoHandler(mockUc)

		// 壊れたJSONデータを送信
		req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBufferString("{invalid-json"))
		rec := httptest.NewRecorder()

		handler.CreateTodo(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Errorf("expected status 400, got %d", rec.Code)
		}

		var res ErrorResponse
		_ = json.Unmarshal(rec.Body.Bytes(), &res)
		if res.Message != "不正なリクエスト形式です" {
			t.Errorf("unexpected error message: %s", res.Message)
		}
	})

	t.Run("失敗: Usecase側でエラー（保存失敗など）が起きた場合、500 Internal Server Errorを返すこと", func(t *testing.T) {
		mockUc := &mockTodoUsecase{
			fakeCreateTodoFn: func(ctx context.Context, title string, description string, accountId *int64, dueDate time.Time) (*model.Todo, error) {
				return nil, errors.New("database breakdown")
			},
		}
		handler := NewTodoHandler(mockUc)

		now := time.Now()

		reqBody := CreateTodoRequest{Title: "エラータスク", AccountId: &[]int64{1}[0], DueDate: &now}
		bodyBytes, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/todos", bytes.NewBuffer(bodyBytes))
		rec := httptest.NewRecorder()

		handler.CreateTodo(rec, req)

		if rec.Code != http.StatusInternalServerError {
			t.Errorf("expected status 500, got %d", rec.Code)
		}

		var res ErrorResponse
		_ = json.Unmarshal(rec.Body.Bytes(), &res)
		if res.Message != "データの保存に失敗しました" || res.Details == nil || (*res.Details)[0] != "database breakdown" {
			t.Errorf("unexpected error response: %+v", res)
		}
	})
}

// -----------------------------------------------------------------------------
// GetTodoById (GET /todos/{id}) のテスト
// -----------------------------------------------------------------------------
func TestTodoHandler_GetTodoById(t *testing.T) {
	t.Run("成功: 存在するIDが指定された場合、200 OKと該当データを返すこと", func(t *testing.T) {
		targetID := uint(456)
		mockUc := &mockTodoUsecase{
			fakeGetTodoByIDFn: func(ctx context.Context, id uint) (*model.Todo, error) {
				if id != targetID {
					return nil, errors.New("wrong id passed to mock")
				}
				return &model.Todo{
					ID:          id,
					Title:       "取得できたタスク",
					Description: "説明書き",
					IsCompleted: true,
					CreatedAt:   time.Now(),
				}, nil
			},
		}
		handler := NewTodoHandler(mockUc)

		req := httptest.NewRequest(http.MethodGet, "/todos/456", nil)
		rec := httptest.NewRecorder()

		// 第3引数に最初からパース済みのIDを渡す(openapi-generatorの定義に準拠)
		handler.GetTodoById(rec, req, int64(targetID))

		if rec.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", rec.Code)
		}

		var res TodoResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if res.Id != int64(targetID) || res.Title != "取得できたタスク" || res.Completed != true {
			t.Errorf("unexpected response: %+v", res)
		}
	})

	t.Run("失敗: Usecaseがエラー（データ不在）を返した場合、404 Not Foundを返すこと", func(t *testing.T) {
		mockUc := &mockTodoUsecase{
			fakeGetTodoByIDFn: func(ctx context.Context, id uint) (*model.Todo, error) {
				return nil, errors.New("not found in db")
			},
		}
		handler := NewTodoHandler(mockUc)

		req := httptest.NewRequest(http.MethodGet, "/todos/999", nil)
		rec := httptest.NewRecorder()

		handler.GetTodoById(rec, req, 999)

		if rec.Code != http.StatusNotFound {
			t.Errorf("expected status 404, got %d", rec.Code)
		}

		var res ErrorResponse
		_ = json.Unmarshal(rec.Body.Bytes(), &res)
		if res.Message != "指定されたタスクが見つかりません" {
			t.Errorf("unexpected error message: %s", res.Message)
		}
	})
}
