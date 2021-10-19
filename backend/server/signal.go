package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	firebase "firebase.google.com/go"
	fa "firebase.google.com/go/auth"
	"google.golang.org/api/option"

	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
)

var AllRooms RoomMap //Typing As RoomMap
var DB *sql.DB //Pointer to DB

type RoomResp struct {
	Id int `json:"roomId"`
	Uuid string `json:"roomUuid"`
	Name string `json:"roomName"`
}

type respMap struct {
	Rooms []RoomResp `json:"rooms"`
}

type FirebaseAuthUser struct {
	Uid string `json:"uid"`
	//Username string
	//Email string
}

//Unused right now

//type User struct {
//	Uuid string `json:"userUuid"`
//	Username string `json:"username"`
//	Password string `json:password`
//}

type ErrorResp struct {
	ErrorMessage string `json:"error"`
}


//MOVE OUT OF HERE??
type Context struct {
	Writer http.ResponseWriter
	Req *http.Request
	Token *fa.Token
}




////////////////////////////////
//// FIX API ERROR HANDLING ////
////////////////////////////////


/////////////////////////////////////
//Move these to seperate Go module//
////////////////////////////////////

var upgrader = websocket.Upgrader{
	HandshakeTimeout: 0,
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	WriteBufferPool:  nil,
	Subprotocols:     []string{},
	Error: func(w http.ResponseWriter, r *http.Request, status int, reason error) {
		fmt.Println("~~~~~~~~~~~ERROR IN CHECK UPGRADE~~~~~~~~~")
		fmt.Println(reason)
	},
	CheckOrigin: func(r *http.Request) bool {
		fmt.Println("~~~~~~~  CHECK ORIGIN REQUEST ~~~~~~~~~")
		return true
	},
	EnableCompression: false,
}

func Upgrader( w http.ResponseWriter, r  *http.Request) (*websocket.Conn, error) {
	fmt.Println("~~~~~~~ CONNECTION UPGRADE ~~~~~~~")
	conn, err := upgrader.Upgrade(w, r, nil)
	fmt.Println("~~~~~~~ CONNECTION ~~~~~~~")
	fmt.Println(conn)
	fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~~~~")

	if err != nil {
		fmt.Fprintf(w, "%+V\n", err)
		log.Fatal("Upgrade error", err)
	}
	fmt.Println("UPGRADE RETURN")
	return conn, nil
}

func handleCors(w *http.ResponseWriter) {
	(*w).Header().Set("Access-Control-Allow-Origin", "http://localhost:8080")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	(*w).Header().Set("Access-Control-Allow-Credentials", "true")
}

//Create room and return RoomID
//rooms/createRoom.go

func CreateRoomRequestHandler(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("~~~~~~~CREATE~~~~~~~")
	 handleCors(&w)

	 if r.Method == "OPTIONS" {
		 return
		}

	 uuid, err := dbCreateRoom("test")

	 if err != nil {
		 fmt.Println(err)
		 return
	 }
	 AllRooms.createRoom(uuid)

	 json.NewEncoder(w).Encode(RoomResp{Uuid: uuid})
 }

//Join Room
//rooms/joinRoom.go
func JoinRoomRequestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("JOIN ROOM")
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
	fmt.Println("~~~~~~~WEBSOCKET~~~~~~~")
	fmt.Println( ws)
	fmt.Println("~~~~~~~~~~~~~~~~~~~~~~~")
	AllRooms.insertIntoRoom(strings.Join(roomId, " "), ws)
}

//Get Room
//rooms/getRooms.go
func GetRoomsRequestHandler(w http.ResponseWriter, r *http.Request)  {
	fmt.Println("~~~~~~ GET ROOMS ~~~~~")
	handleCors(&w)
	if r.Method == "OPTIONS" {
		return
 	}

	rooms, err := getRooms()

	if err != nil {
		fmt.Println(err)
		return
	}

	json.NewEncoder(w).Encode(respMap{Rooms: rooms})
}

//register/register.go
func RegisterUserRequestHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("~~~~~~~REGISTER~~~~~~~")



	if r.Method != http.MethodPost {
		log.Printf("Require POST Request")
		w.WriteHeader(405)
		json.NewEncoder(w).Encode(ErrorResp{ErrorMessage: "Method not allowed: Route requires POST Request"})
		return
	}

	var user FirebaseAuthUser

	//Read Request Body and return error if unsuccessful
	bodyErr := json.NewDecoder(r.Body).Decode(&user)
	if bodyErr != nil {
		log.Printf("%v", bodyErr)
		w.WriteHeader(500)
		return
	}

	err := insertUser(user.Uid)
	if err != nil {
		log.Printf("%v", err)
		w.WriteHeader(500)
		return
	}

	json.NewEncoder(w).Encode(user)
}


/////////////////////////////////////
//Move these to seperate Go module?//
////////////////////////////////////




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

		fmt.Println("~~~~~~ ROOMS ~~~~~~~")
		fmt.Println(rooms)
		fmt.Println("~~~~~~~~~~~~~~~~~~~~")

		return rooms, nil
}

//Insert User into DB
func insertUser(uid string) (error){
	query := fmt.Sprintf("INSERT into firebase_user (uid) VALUES ('%s') RETURNING uid;", uid)

	var uuid string
	err := DB.QueryRow(query).Scan(&uuid)

	if err != nil {
		fmt.Println(err)
		return  err
	}

	return  nil
}


func dbCreateRoom(roomName string ) (string, error) {
	query := fmt.Sprintf("INSERT into rooms VALUES (nextval('users_user_id_seq'::regclass), default, '%s') RETURNING room_uuid;", roomName)

	var uuid string
	err := DB.QueryRow(query).Scan(&uuid)

	if err != nil {
		fmt.Println(err)
		return "", err
	}

	return uuid, nil
}


func loadDotEnv(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println(err)
	}

	return os.Getenv(key)
}

func FirebaseAuthRoute (handler http.Handler) http.HandlerFunc {
	return func (w http.ResponseWriter, r *http.Request) {
		handleCors(&w)

		if r.Method == "OPTIONS" {
			return
		}

		rc := &Context{
			Writer: w,
			Req: r,
			Token: nil,
		}
		path := loadDotEnv("FIREBASE_SECRET_PATH")
		opt := option.WithCredentialsFile(path)
		app, err := firebase.NewApp(context.Background(), nil, opt)

		if err != nil {
			fmt.Println("ERROR IN FIREBASE APP")
			return
		}

		auth, err := app.Auth(context.Background())

		if err != nil {
			fmt.Println("ERROR IN FIREBASE AUTH")
			return
		}

		header := r.Header.Get("Authorization")
		idToken := strings.TrimSpace(strings.Replace(header, "Bearer", "", 1))

		token , err := auth.VerifyIDToken(context.Background(), idToken)
		if err != nil {
			fmt.Println("ERROR IN FIREBASE TOKEN")
			fmt.Println(idToken)
			fmt.Println(err)
			return
		}

		rc.Token = token
		handler.ServeHTTP(w,r)
	}
}




///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////
///////////////////LEGACY LOGIN STUFF//////////////////////
///////////////////////////////////////////////////////////
///////////////////////////////////////////////////////////


//type Authentication struct {
//	Username string
//	Password string
//}

//type Token struct {
//	Username string `json:"username"`
//	Uuid string `json:"userUuid"`
//	TokenString string `json:"token"`
//}

//type UuidResp struct {
//	Uuid string `json:"userUuid"`
//}

//register/register.go
//func RegisterUserRequestHandler(w http.ResponseWriter, r *http.Request) {
	//fmt.Println("~~~~~~~REGISTER~~~~~~~")
	//handleCors(&w)
	//if r.Method == "OPTIONS" {
	//	return
 	//}

	//if r.Method != http.MethodPost {
	//	log.Printf("Require POST Request")
	//	w.WriteHeader(405)
	//	json.NewEncoder(w).Encode(ErrorResp{ErrorMessage: "Method not allowed: Route requires POST Request"})
	//	return
	//}

	//var authDetails Authentication

	////Read Request Body and return error if unsuccessful
	//bodyErr := json.NewDecoder(r.Body).Decode(&authDetails)
	//if bodyErr != nil {
	//	log.Printf("%v", bodyErr)
	//	w.WriteHeader(500)
	//	return
	//}

//	//Check if username is valid
//	isValid := checkValidUsername(authDetails.Username)
//	fmt.Println(isValid)
//	if !isValid {
//		w.WriteHeader(500)
//		return
//	}

//	//Hash password with bcrypt
//	hash, err := bcrypt.GenerateFromPassword([]byte(authDetails.Password), bcrypt.DefaultCost)
//	if err != nil {
//		log.Printf("%v", err)
//		w.WriteHeader(500)
//		return
//	}
//	fmt.Println(string(hash))

//	//Convert hash byte slice into string, insert user into DB.
//	//If successful return DB generated Uuid if not return error.
	//user, err := insertUser(authDetails.Username, string(hash))
	//if err != nil {
	//	log.Printf("%v", err)
	//	w.WriteHeader(500)
	//}

//	validToken, err := GenerateJWT(user.Username, user.Uuid)
//	if err != nil {
//		fmt.Println("error generating token: ", err)
//		w.WriteHeader(400)
//		json.NewEncoder(w).Encode(ErrorResp{ErrorMessage: err.Error()})
//		return
//	}

//	var token Token
//	token.Username = user.Username
//	token.Uuid = user.Uuid
//	token.TokenString = validToken
//	cookie := http.Cookie{
//		Name:       "auth-cookie",
//		Value:      token.TokenString,
//		HttpOnly:   true,
//	}
//	http.SetCookie(w, &cookie)

//	json.NewEncoder(w).Encode(token)
//}

////login/login.go
//func LoginRequestHandler(w http.ResponseWriter, r *http.Request) {
//	handleCors(&w)
//	if r.Method == "OPTIONS" {
//		return
// 	}

//	if r.Method != http.MethodPost {
//		log.Printf("Require POST Request")
//		w.WriteHeader(405)
//		json.NewEncoder(w).Encode(ErrorResp{ErrorMessage: "Method not allowed: Route requires POST Request"})
//		return
//	}

//	var authDetails Authentication

//	//Read Request Body and return error if unsuccessful
//	bodyErr := json.NewDecoder(r.Body).Decode(&authDetails)
//	if bodyErr != nil {
//		log.Printf("%v", bodyErr)
//		w.WriteHeader(400)
//	}

//	user, err := validateLogin(authDetails.Username, authDetails.Password)
//	if err != nil {
//		fmt.Println("error validating: ", err)
//		w.WriteHeader(400)
//		json.NewEncoder(w).Encode(ErrorResp{ErrorMessage: err.Error()})
//		return
//	}

//	validToken, err := GenerateJWT(user.Username, user.Uuid)
//	if err != nil {
//		fmt.Println("error generating token: ", err)
//		w.WriteHeader(400)
//		json.NewEncoder(w).Encode(ErrorResp{ErrorMessage: err.Error()})
//		return
//	}

//	var token Token
//	token.Username = user.Username
//	token.Uuid = user.Uuid
//	token.TokenString = validToken
//	expiration := time.Now().AddDate(0,1,0)
//	cookie := http.Cookie{
//		Name:       "session",
//		Value:      token.TokenString,
//		Path:       "/",
//		Domain:     "localhost",
//		Expires:    expiration,
//		RawExpires: "",
//		MaxAge:     0,
//		Secure:     true,
//		HttpOnly:   false,
//		SameSite:   http.SameSiteNoneMode,
//		Raw:        "",

//	}
//	//cookie.Domain =
//	http.SetCookie(w, &cookie)
//	//fmt.Println(uuid)
//	w.Header().Set("Content-Type", "application/json")
//	json.NewEncoder(w).Encode(token)
//}







//Query DB for username and compare hashed password
//func validateLogin(username string, password string) (User, error) {
//	var uuid, dbUsername, dbPassword string
//	var user User
//	query := fmt.Sprintf("SELECT user_uuid, username, password FROM users where username = '%s'", username)

//	err := DB.QueryRow(query).Scan(&uuid, &dbUsername, &dbPassword )
//	if err != nil {
//		fmt.Println(err)
//		return user, err
//	}

//	bcryptErr := bcrypt.CompareHashAndPassword([]byte(dbPassword), []byte(password))
//	if bcryptErr != nil {
//		fmt.Println(bcryptErr)
//		return user, bcryptErr
//	}
//	user = User{
//		uuid,
//		username,
//		dbPassword,
//	}
//	return user, nil
//}

//Check DB for username
//func checkValidUsername(username string) (bool) {
//	var uuid string

//	query := fmt.Sprintf("SELECT user_uuid FROM users WHERE username = '%s'", username)
//	err := DB.QueryRow(query).Scan(&uuid)

//	if err == sql.ErrNoRows {
//		return true
//	} else {
//		return false
//	}
//}



//func AuthorizedRoute (handler http.Handler) http.HandlerFunc {
//	return func (w http.ResponseWriter, r *http.Request) {
//		fmt.Printf("~~~~~~~~~ HANDLER ~~~~~~")
//		fmt.Println(r.Header["Token"])
//		if r.Header["Token"] == nil {
//			fmt.Printf("Token Not Found")
//			w.WriteHeader(401)
//			//json.NewEncoder(w).Encode(err)
//			json.NewEncoder(w)//.Encode(ErrorResp{ErrorMessage: err.Error()})
//			return
//		}
//		fmt.Printf("~~~~~~~~~ HANDLER Token ~~~~~~")

//		tokenString := r.Header["Token"][0]
//		signingKey := []byte(loadDotEnv("SECRET_KEY"))

//		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
//			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
//				return nil, fmt.Errorf("unexptected signing method: %v", token.Header["alg"])
//			}

//			return signingKey, nil
//		})

//		if err != nil {
//			fmt.Printf("Token Expired: %v", err.Error())
//			json.NewEncoder(w).Encode(ErrorResp{ErrorMessage: err.Error()})
//			return
//		}


//		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
//				if claims["authorized"] == true {
//					handler.ServeHTTP(w,r)
//				}
//		}
//	}
//}


//func GenerateJWT (username, uuid string) (string, error) {
//	signingKey := []byte(loadDotEnv("SECRET_KEY"))
//	token := jwt.New(jwt.SigningMethodHS256)
//	claims := token.Claims.(jwt.MapClaims)

//	claims["authorized"] = true
//	claims["username"] = username
//	claims["uuid"] = uuid
//	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

//	tokenString, err := token.SignedString(signingKey)

//	if err != nil {
//		fmt.Errorf("Something went Wrong: %s", err.Error())
//		return "", err
//	}

//	return tokenString, nil
//}

//func insertUser(username string, password string) (User, error){
//	query := fmt.Sprintf("INSERT into users VALUES (nextval('users_user_id_seq'::regclass), default, '%s', '%s') RETURNING user_uuid;", username, password)

//	var uuid string
//	var user User
//	err := DB.QueryRow(query).Scan(&uuid)

//	if err != nil {
//		fmt.Println(err)
//		return user, err
//	}

//	user = User {
//		uuid,
//		username,
//		password,
//	}

//	return user, nil
//}


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