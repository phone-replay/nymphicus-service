package src

type Action struct {
	Action      string `json:"action"`
	TargetTime  string `json:"targetTime"`
	Coordinates string `json:"coordinates"`
}

type Gesture struct {
	Actions []Action `json:"actions"`
}

type ActivityGesture struct {
	ActivityName string    `json:"activityName"`
	Gestures     []Gesture `json:"gestures"`
}

type ActivityGestureLogs struct {
	Activities []ActivityGesture `json:"activities"`
}
