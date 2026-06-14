package handler

// HealthCheckResponse はヘルスチェックAPIのレスポンス形式です
type HealthCheckResponse struct {
	Status    string `json:"status"`              // "OK" または "FAIL"
	Timestamp string `json:"timestamp"`           // 現在時刻
	Database  string `json:"database,omitempty"`  // DBの接続状態（オプション）
}