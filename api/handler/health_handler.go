package handler

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type HealthHandler struct {
	db *sql.DB
}

// NewHealthHandler はHealthHandlerのインスタンスを生成します
func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{db: db}
}

// Check はヘルスチェックを行うエンドポイントです (GET /health)
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	res := HealthCheckResponse{
		Status:    "OK",
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Database:  "CONNECTED",
	}

	// 1. PostgreSQLへの疎通確認 (Readiness Check)
	// DBインスタンスが存在し、かつPingが通るか確認
	if h.db != nil {
		if err := h.db.PingContext(r.Context()); err != nil {
			// DBとの接続が切れている場合
			w.WriteHeader(http.StatusServiceUnavailable) // 503を返す
			res.Status = "FAIL"
			res.Database = "DISCONNECTED"
			_ = json.NewEncoder(w).Encode(res)
			return
		}
	} else {
		// そもそもDB初期化に失敗している場合
		w.WriteHeader(http.StatusServiceUnavailable)
		res.Status = "FAIL"
		res.Database = "NOT_INITIALIZED"
		_ = json.NewEncoder(w).Encode(res)
		return
	}

	// 2. すべて正常なら 200 OK を返す
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(res)
}