package unipile

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"chatsheet/config"
)

// API 響應結構

// CheckpointResponse 處理 Unipile 返回的 Checkpoint 結構
type CheckpointResponse struct {
	Object     string `json:"object"`
	AccountID  string `json:"account_id"`
	Checkpoint *struct {
		Type string `json:"type"` // 例如: "2FA", "OTP", "IN_APP_VALIDATION", "CAPTCHA"
	} `json:"checkpoint"`
	// 其他成功的欄位，例如 provider, status
}

// Unipile API 端點
const (
	AccountsEndpoint   = "/api/v1/accounts"
	CheckpointEndpoint = "/api/v1/accounts/checkpoint"
)

// PerformRequest 執行對 Unipile API 的 POST 請求
func PerformRequest(cfg config.UnipileConfig, endpoint string, data interface{}, target interface{}) (int, error) {
	url := cfg.APIBaseURL + endpoint

	jsonBody, err := json.Marshal(data)
	if err != nil {
		return 0, fmt.Errorf("無法序列化請求體: %w", err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return 0, fmt.Errorf("無法建立請求: %w", err)
	}

	req.Header.Set("X-API-KEY", cfg.APIKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("API 請求失敗: %w", err)
	}
	defer resp.Body.Close()

	// 讀取響應體
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return resp.StatusCode, fmt.Errorf("讀取響應體失敗: %w", err)
	}

	// 檢查狀態碼
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusAccepted {
		// 即使是錯誤，也嘗試解析為 CheckpointResponse 以獲取可能的 account_id
		var checkpoint CheckpointResponse
		if json.Unmarshal(bodyBytes, &checkpoint) == nil && checkpoint.Object == "Checkpoint" {
			// 成功解析為 Checkpoint，雖然不是 200/202，但我們將其視為 Checkpoint 流程
			// 實際情況應為 202，但為了穩健性，這裡將其視為成功的 Checkpoint 響應
			if target != nil {
				json.Unmarshal(bodyBytes, target)
			}
			return resp.StatusCode, nil
		}

		slog.Error("Unipile API 請求失敗", "status", resp.StatusCode, "body", string(bodyBytes))
		return resp.StatusCode, fmt.Errorf("Unipile API 錯誤: %s", string(bodyBytes))
	}

	// 成功或 202 (Accepted/Checkpoint)
	if target != nil {
		if err := json.Unmarshal(bodyBytes, target); err != nil {
			return resp.StatusCode, fmt.Errorf("解析響應失敗: %w", err)
		}
	}

	return resp.StatusCode, nil
}
