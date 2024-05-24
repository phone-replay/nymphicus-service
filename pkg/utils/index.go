package utils

import (
	"nymphicus-service/src/models"
	"sort"
)

type TimeLine struct {
	Coordinates string `json:"coordinates"`
	GestureType string `json:"gestureType"`
	TargetTime  string `json:"targetTime"`
}

func GetTimeLines(activities []models.Activity) []TimeLine {
	var timeLines []TimeLine
	for _, activity := range activities {
		for _, gesture := range activity.Gestures {
			timeLine := TimeLine{
				Coordinates: gesture.Coordinates,
				GestureType: gesture.GestureType,
				TargetTime:  gesture.TargetTime,
			}
			timeLines = append(timeLines, timeLine)
		}
	}
	sort.Slice(timeLines, func(i, j int) bool {
		return timeLines[i].TargetTime < timeLines[j].TargetTime
	})

	return timeLines
}
