package main

type EventType int

const (
	joinRoom EventType = iota
	createRoom
	invalid
)

func getEventType(data map[string]interface{}) EventType {
	msgType, exists := data["msgType"]
	if !exists {
		return invalid
	}

	return EventType(msgType.(float64))
}
