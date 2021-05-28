package main

import (
	"backend/server"
	"log"
	"net/http"
)

func main() {
	server.AllRooms.Init()

	http.HandleFunc("/create", server.CreateRoomRequestHandler)
	http.HandleFunc("/join", server.JoinRoomRequestHandler)
	http.HandleFunc("/get", server.GetRoomsRequestHandler)

	log.Println("Starting Server on :8000")
	http.ListenAndServe(":8000", nil)
}