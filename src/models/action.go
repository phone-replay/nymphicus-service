package models

type Action struct {
	ID         string     `json:"id"`
	Activities []Activity `json:"activities"`
	Device     Device     `json:"device"`
	VideoUrl   *string    `json:"videoUrl"`
	Key        string     `json:"key"`
}

type Activity struct {
	ActivityName string    `json:"activityName"`
	Gestures     []Gesture `json:"gestures"`
	ID           string    `json:"id"`
}

type Gesture struct {
	ActivityId  string `json:"activityId"`
	Coordinates string `json:"coordinates"`
	CreatedAt   string `json:"createdAt"`
	GestureType string `json:"gestureType"`
	TargetTime  string `json:"targetTime"`
}
