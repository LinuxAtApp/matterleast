package datalog

import "sync"

type ChatHistory struct {
	Name string
	History []string
	Authors []string
	Mux sync.Mutex
}

var TownSquare = ChatHistory{ Name: "Town Square" }
var OffTopic = ChatHistory{ Name: "Off-Topic" }
var TermDev = ChatHistory{ Name: "Terminal UI Development" }
