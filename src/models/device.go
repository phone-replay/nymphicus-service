package models

type Device struct {
	BatteryLevel     float64 `json:"batteryLevel"`
	Brand            string  `json:"brand"`
	CurrentNetwork   string  `json:"currentNetwork"`
	Device           string  `json:"device"`
	InstallID        string  `json:"installID"`
	Language         string  `json:"language"`
	Manufacturer     string  `json:"manufacturer"`
	Model            string  `json:"model"`
	OsVersion        string  `json:"osVersion"`
	Platform         string  `json:"platform"`
	ScreenResolution string  `json:"screenResolution"`
	SdkVersion       int     `json:"sdkVersion"`
	SessionId        string  `json:"sessionId"`
	TotalRAM         string  `json:"totalRAM"`
	TotalStorage     string  `json:"totalStorage"`
}
