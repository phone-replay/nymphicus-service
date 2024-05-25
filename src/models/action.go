package models

type Action struct {
	Activities []Activity `json:"activities"`
	Device     Device     `json:"device"`
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
