package handler

import (
	"encoding/json"
	"net/http"
	"time"
	"todo/base/usecase"
    "context"
)

// 1. Handlerが求めるUsecaseのインターフェースを定義
type HealthUsecaseInterface interface {
    DiagnoseHealth(ctx context.Context) usecase.HealthStatus
}

type HealthHandler struct {
    // 2. 構造体ではなく、インターフェースに依存させる
    usecase HealthUsecaseInterface 
}

// NewHealthHandler はHealthHandlerのインスタンスを生成します
func NewHealthHandler(u HealthUsecaseInterface) *HealthHandler {
	return &HealthHandler{usecase: u}
}

func (h *HealthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// 1. Usecaseを呼び出して、ビジネスロジック（診断）を実行
	health := h.usecase.DiagnoseHealth(r.Context())

	// 2. 診断結果をもとにレスポンスを設定
	res := HealthCheckResponse{
		Status:    "OK",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Database:  health.DBStatus,
	}

	// 健康状態が異常（false）な場合は503をセットして即座に返却
	if !health.IsHealthy {
		res.Status = "FAIL"
		w.WriteHeader(http.StatusServiceUnavailable)
		_ = json.NewEncoder(w).Encode(res)
		return
	}

	// 3. すべて正常なら 200 OK を返す
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}
