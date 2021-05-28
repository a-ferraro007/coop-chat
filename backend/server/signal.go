package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var AllRooms RoomMap //Typing As CoopMap
var id string

type resp struct {
	RoomID string `json:"roomID"`
}

type respMap struct {
	Rooms []string `json:"rooms"`
}

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,

	//Check origin of the connection and allow any connection for now
	CheckOrigin: func(r *http.Request) bool { return true },
}

func Upgrader( w http.ResponseWriter, r  *http.Request) (*websocket.Conn, error) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Fprintf(w, "%+V\n", err)
		log.Fatal("Upgrade error", err)
	}

	return conn, nil
}

//Create Coop and return CoopID
func CreateRoomRequestHandler(w http.ResponseWriter, r *http.Request)  {

	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))

	id = AllRooms.createRoom()
	fmt.Println(id)

	json.NewEncoder(w).Encode(resp{RoomID: id})
}

//Join Coop
func JoinRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	roomId, ok := r.URL.Query()["roomID"]

	if !ok {
		fmt.Println("Missing room Id")
		return
	}

	ws, _  := Upgrader(w, r)
	AllRooms.insertIntoRoom(strings.Join(roomId, " "), ws)
}

func GetRoomsRequestHandler(w http.ResponseWriter, r *http.Request)  {
	w.Header().Set("Access-Control-Allow-Origin", r.Header.Get("Origin"))
	resp := make([]string, 0)

	for i, _ := range AllRooms.Map {
		resp = append(resp, i)
	}

	fmt.Println(resp)
	json.NewEncoder(w).Encode(respMap{Rooms: resp})
}

