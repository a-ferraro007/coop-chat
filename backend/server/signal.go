package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"golang.org/x/crypto/bcrypt"
)

var AllRooms RoomMap //Typing As RoomMap
var id string
var DB *sql.DB //Pointer to DB

type RoomResp struct {
	Id int `json:"roomId"`
	Uuid string `json:"roomUuid"`
	Name string `json:"roomName"`
}

type respMap struct {
	Rooms []RoomResp `json:"rooms"`
}

type UuidResp struct {
	Uuid string `json:"userUuid"`
}

type ErrorResp struct {
	ErrorMessage string `json:"error"`
}

type LoginBody struct {
	Username string
	Password string
}

////////////////////////////////
//// FIX API ERROR HANDLING ////
////////////////////////////////

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

func handleCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

//Create Coop and return CoopID
func CreateRoomRequestHandler(w http.ResponseWriter, r *http.Request)  {
	handleCors(&w)
	if r.Method == "OPTIONS" {
		return
 	}

	id = AllRooms.createRoom()
	fmt.Println(id)

	//json.NewEncoder(w).Encode(resp{RoomID: id})
}

//Join Room
func JoinRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	handleCors(&w)
	if r.Method == "OPTIONS" {
		return
 	}

	roomId, ok := r.URL.Query()["roomID"]

	if !ok {
		fmt.Println("Missing room Id")
		return
	}

	ws, _  := Upgrader(w, r)
	AllRooms.insertIntoRoom(strings.Join(roomId, " "), ws)
}

//Get Room
func GetRoomsRequestHandler(w http.ResponseWriter, r *http.Request)  {
	handleCors(&w)
	if r.Method == "OPTIONS" {
		return
 	}

	rooms, err := getRooms()

	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(rooms)
	json.NewEncoder(w).Encode(respMap{Rooms: rooms})
}

func RegisterUserRequestHandler(w http.ResponseWriter, r *http.Request) {
	handleCors(&w)
	if r.Method == "OPTIONS" {
		return
 	}

	if r.Method != http.MethodPost {
		log.Printf("Require POST Request")
		w.WriteHeader(405)
		json.NewEncoder(w).Encode(ErrorResp{ErrorMessage: "Method not allowed: Route requires POST Request"})
		return
	}

	var body LoginBody

	//Read Request Body and return error if unsuccessful
	bodyErr := json.NewDecoder(r.Body).Decode(&body)
	if bodyErr != nil {
		log.Printf("%v", bodyErr)
		w.WriteHeader(500)
		return
	}

	//Check if username is valid
	isValid := checkValidUsername(body.Username)
	fmt.Println(isValid)
	if !isValid {
		w.WriteHeader(500)
		return
	}

	//Hash password with bcrypt
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(500)
		return
	}
	fmt.Println(string(hash))

	//Convert hash byte slice into string, insert user into DB.
	//If unsuccessful return DB generated Uuid if not return error.
	//Convert hash byte slice into string.
	uuid, err := insertUser(body.Username, string(hash))
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(500)
	}

	json.NewEncoder(w).Encode(UuidResp{Uuid: uuid})
}

func LoginRequestHandler(w http.ResponseWriter, r *http.Request) {
	handleCors(&w)
	if r.Method == "OPTIONS" {
		return
 	}

	if r.Method != http.MethodPost {
		log.Printf("Require POST Request")
		w.WriteHeader(405)
		json.NewEncoder(w).Encode(ErrorResp{ErrorMessage: "Method not allowed: Route requires POST Request"})
		return
	}

	var body LoginBody

	//Read Request Body and return error if unsuccessful
	bodyErr := json.NewDecoder(r.Body).Decode(&body)
	if bodyErr != nil {
		log.Printf("%v", bodyErr)
		w.WriteHeader(400)
	}

	uuid, err := validateLogin(body.Username, body.Password)
	if err != nil {
		fmt.Println("error validating: ", err)
		w.WriteHeader(400)
		json.NewEncoder(w).Encode(ErrorResp{ErrorMessage: err.Error()})
		return
	}
	fmt.Println(uuid)
	json.NewEncoder(w).Encode(UuidResp{Uuid: uuid})
}


/////////////////////////////////////
//Move these to seperate Go module?//
////////////////////////////////////

//Check DB for username
func checkValidUsername(username string) (bool) {
	var uuid string

	query := fmt.Sprintf("SELECT user_uuid FROM users WHERE username = '%s'", username)
	err := DB.QueryRow(query).Scan(&uuid)

	if err == sql.ErrNoRows {
		return true
	} else {
		return false
	}
}


//Return Rooms from DB
func getRooms() ([]RoomResp, error)  {
	rows, err := DB.Query(`SELECT * from rooms`,)
	rooms := make([]RoomResp, 0)

	defer rows.Close()

	if err != nil {
		return nil, err
	}
		fmt.Println("rooms")
		fmt.Println(rows)

		for rows.Next() {
			room := RoomResp{}

			err := rows.Scan(&room.Id, &room.Uuid, &room.Name)
			if err != nil {
				return nil, err
			}

			rooms = append(rooms, room)
		}

		return rooms, nil
}

//Insert User into DB
//Usernames can be duplicated
func insertUser(username string, password string) (string, error){
	query := fmt.Sprintf("INSERT into users VALUES (nextval('users_user_id_seq'::regclass), default, '%s', '%s') RETURNING user_uuid;", username, password)

	var uuid string
	err := DB.QueryRow(query).Scan(&uuid)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return uuid, nil
}

//Query DB for username and compare hashed password
func validateLogin(username string, password string) (string, error) {
	var uuid, dbUsername, dbPassword string
	query := fmt.Sprintf("SELECT user_uuid, username, password FROM users where username = '%s'", username)

	err := DB.QueryRow(query).Scan(&uuid, &dbUsername, &dbPassword)
	if err != nil {
		fmt.Println(err)
		return "", err
	}

	bcryptErr := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
	if bcryptErr != nil {
		fmt.Println(bcryptErr)
		return "", bcryptErr
	}

	return uuid, nil
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

//INSERT USER
//insert into users values (nextval('users_user_id_seq'::regclass), default, username, password)
//INSERT USER RETURN UUID
// insert into users values (nextval('users_user_id_seq'::regclass), default, 'test', 'test') returning user_uuid;

//Select user from userID
//select user_uuid, username, password from users where username = '%s'