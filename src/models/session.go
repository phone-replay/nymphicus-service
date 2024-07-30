package models

import (
	"nymphicus-service/src"
	"time"
)

type Session struct {
	ID         string                  `json:"id"`
	Activities src.ActivityGestureLogs `json:"activities"`
	Device     Device                  `json:"device"`
	VideoUrl   *string                 `json:"videoUrl"`
	Status     string                  `json:"status"`
	CreatedAt  time.Time               `json:"createdAt"`
	Key        string                  `json:"key"`
	Duration   int64                   `json:"duration"`
}
