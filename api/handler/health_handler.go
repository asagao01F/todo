package handler

import (
    "encoding/json"
    "net/http"
    "time"
    "gorm.io/gorm"
)

type HealthHandler struct {
    // 【変更】*sql.DB から *gorm.DB に変更
    db *gorm.DB
}

// NewHealthHandler はHealthHandlerのインスタンスを生成します
// 【変更】引数の型を *gorm.DB に変更
func NewHealthHandler(db *gorm.DB) *HealthHandler {
    return &HealthHandler{db: db}
}

func (h *HealthHandler) CheckHealth(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "application/json")

    // HealthCheckResponse は openapi.gen.go で自動生成された型をそのまま使用
    res := HealthCheckResponse{
        Status:    "OK",
        Timestamp: time.Now().UTC().Format(time.RFC3339),
        Database:  "CONNECTED",
    }

    // 1. PostgreSQLへの疎通確認 (Readiness Check)
    if h.db != nil {
        // 【変更】GORMから内部の *sql.DB を取得する
        sqlDB, err := h.db.DB()
        if err != nil {
            // *sql.DBの取得自体に失敗した場合
            w.WriteHeader(http.StatusServiceUnavailable)
            res.Status = "FAIL"
            res.Database = "ERROR"
            _ = json.NewEncoder(w).Encode(res)
            return
        }

        // 取得した sqlDB を使って Ping を打つ
        if err := sqlDB.PingContext(r.Context()); err != nil {
            // DBとの接続が切れている場合
            w.WriteHeader(http.StatusServiceUnavailable) // 503を返す
            res.Status = "FAIL"
            res.Database = "DISCONNECTED"
            _ = json.NewEncoder(w).Encode(res)
            return
        }
    } else {
        // そもそもDB初期化に失敗している場合
        w.WriteHeader(http.StatusServiceUnavailable) // 503を返す
        res.Status = "FAIL"
        res.Database = "NOT_INITIALIZED"
        _ = json.NewEncoder(w).Encode(res)
        return
    }

    // 2. すべて正常なら 200 OK を返す
    w.WriteHeader(http.StatusOK)
    _ = json.NewEncoder(w).Encode(res)
}
