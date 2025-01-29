package domain

import (
	"time"
)

type ResponseMetrics struct {
	Method     string        `json:"method"`
	Path       string        `json:"path"`
	StatusCode int           `json:"status_code"`
	Time       time.Duration `json:"elapsed_time"`
	Size       int           `json:"response_size"`
}
