package server

import (
	"pinx/game_server/move"
	"pinx/game_server/offline"
	online "pinx/game_server/playonline"
	"pinx/game_server/talk"
	"pinx/pinterface"
	"pinx/pnet"
)

var Serve pinterface.IServer

func init() {
	Serve = pnet.NewServer("Game PK!")
	Serve.SetOnConnStart(online.Online)
	Serve.SetOnConnStop(offline.Offline)
	Serve.AddRouter(2, &talk.WorldChatApi{})
	Serve.AddRouter(3, &move.Move{})
}

func Run() {
	Serve.Serve()
}
