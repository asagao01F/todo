package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"todo/base/usecase"
)


// 1. HealthUsecaseの実態の代わりに使う「手動モック」を定義
type mockHealthUsecase struct {
	mockResult usecase.HealthStatus
}

// usecase.HealthUsecase と同じメソッドを実装してインターフェースを満たす
// ※ もしusecase側が構造体で定義されている場合は、usecase側をインターフェース化するか、
// このモックがusecaseと同じ振る舞いをするように設定します。
func (m *mockHealthUsecase) DiagnoseHealth(ctx context.Context) usecase.HealthStatus {
	return m.mockResult
}

func TestHealthHandler_CheckHealth(t *testing.T) {
	t.Run("成功: DB接続が正常な場合、200 OKを返すこと", func(t *testing.T) {
		// Usecaseのモックを作成し、成功時のステータスを設定
		mockUc := &mockHealthUsecase{
			mockResult: usecase.HealthStatus{
				IsHealthy: true,
				DBStatus:  "CONNECTED",
			},
		}
		// Usecaseの代わりにモックを注入（NewHealthHandlerの引数の型に合わせて調整してください）
		// ※ もしHandlerの引数が具象構造体の場合は、インターフェースを挟む形に修正すると綺麗にテストできます。
		// ここでは、NewHealthHandler が usecase のインターフェースまたはモックを受け取れる前提として記述します。
		handler := NewHealthHandler(mockUc)

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()

		handler.CheckHealth(rec, req)

		// アサーション
		if rec.Code != http.StatusOK {
			t.Errorf("expected status code %d, got %d", http.StatusOK, rec.Code)
		}

		var res HealthCheckResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if res.Status != "OK" || res.Database != "CONNECTED" {
			t.Errorf("unexpected response content: %+v", res)
		}
		if res.Timestamp == "" {
			t.Error("expected timestamp to be not empty")
		}
	})

	t.Run("失敗: DB接続が切れている場合、503 Service Unavailableを返すこと", func(t *testing.T) {
		// Usecaseのモックに失敗（切断）時のステータスを設定
		mockUc := &mockHealthUsecase{
			mockResult: usecase.HealthStatus{
				IsHealthy: false,
				DBStatus:  "DISCONNECTED",
			},
		}
		handler := NewHealthHandler(mockUc)

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()

		handler.CheckHealth(rec, req)

		// アサーション
		if rec.Code != http.StatusServiceUnavailable {
			t.Errorf("expected status code %d, got %d", http.StatusServiceUnavailable, rec.Code)
		}

		var res HealthCheckResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if res.Status != "FAIL" || res.Database != "DISCONNECTED" {
			t.Errorf("unexpected response content: %+v", res)
		}
	})

	t.Run("失敗: そもそもDBインスタンスがnilの場合、503 Service Unavailableを返すこと", func(t *testing.T) {
		// Usecaseのモックに未初期化時のステータスを設定
		mockUc := &mockHealthUsecase{
			mockResult: usecase.HealthStatus{
				IsHealthy: false,
				DBStatus:  "NOT_INITIALIZED",
			},
		}
		handler := NewHealthHandler(mockUc)

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rec := httptest.NewRecorder()

		handler.CheckHealth(rec, req)

		// アサーション
		if rec.Code != http.StatusServiceUnavailable {
			t.Errorf("expected status code %d, got %d", http.StatusServiceUnavailable, rec.Code)
		}

		var res HealthCheckResponse
		if err := json.Unmarshal(rec.Body.Bytes(), &res); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if res.Status != "FAIL" || res.Database != "NOT_INITIALIZED" {
			t.Errorf("unexpected response content: %+v", res)
		}
	})
}
