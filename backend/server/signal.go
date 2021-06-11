package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
)

var AllRooms RoomMap //Typing As CoopMap
var id string
var DB *sql.DB

type Resp struct {
	Id int `json:"id"`
	Uuid string `json:"uuid"`
	Name string `json:"name"`
}

type respMap struct {
	Rooms []Resp `json:"rooms"`
}

//type Room struct {
//	id int
//	uuid string
//	name string
//}

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

	//json.NewEncoder(w).Encode(resp{RoomID: id})
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
	//fmt.Println(DB)
		rooms, err := getRooms()

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println(rooms)
		json.NewEncoder(w).Encode(respMap{Rooms: rooms})
	//resp := make([]string, 0)

	//for i, _ := range AllRooms.Map {
	//	resp = append(resp, i)
	//}

	//fmt.Println(resp)

}

func getRooms() ([]Resp, error)  {
	rows, err := DB.Query(`SELECT * from rooms`,)
	rooms := make([]Resp, 0)

	defer rows.Close()

	if err != nil {
		return nil, err
	}
		fmt.Println("rooms")
		fmt.Println(rows)

		for rows.Next() {
			room := Resp{}

			err := rows.Scan(&room.Id, &room.Uuid, &room.Name)
			if err != nil {
				return nil, err
			}

			rooms = append(rooms, room)
		}

			return rooms, nil
}




//JOIN ROOM
//insert into active_rooms (user_id, room_id, user_uuid, room_uuid)
//select user_id, room_id, user_uuid, room_uuid from users
//right join rooms on users.user_id = rooms.room_id
//where user_uuid = '93c2f9e2-f608-4ad5-a2ac-706856aec5ba';

//USERS IN ROOM
//select user_id, user_uuid
//from active_rooms
//where room_uuid = 'e6ba496a-7134-43e0-ac67-b7d76c565502'