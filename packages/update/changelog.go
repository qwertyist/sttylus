package main

import "log"

type LogType string

const (
	MinorChange  LogType = "minor"
	MajorChange  LogType = "update"
	BugFixChange LogType = "fix"
)

type LogEntry struct {
	Type        LogType
	Description string
}

func initChangeLog(m *Manifest) {
	if m.Log == nil {
		log.Println("No log")
		m.Log = []*LogEntry{}
	}
}
