package main

import (
	"backend/server"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)



func main() {
	dbUrl := loadDotEnv("dbURL")
	db, err := connect(dbUrl)

	fmt.Println(db)
	fmt.Println(err)

	server.AllRooms.Init()
	server.DB = db //I need a reference to the DB connection in the signal go module

	//API Routes
	http.HandleFunc("/create", server.CreateRoomRequestHandler)
	http.HandleFunc("/join", server.JoinRoomRequestHandler)
	http.HandleFunc("/get", server.GetRoomsRequestHandler)
	http.HandleFunc("/register", server.RegisterUserRequestHandler)
	http.HandleFunc("/login", server.LoginRequestHandler)

	log.Println("Starting Server on :8000")
	http.ListenAndServe(":8000", nil)
}

//Connect to Postgres DB
func connect(connURL string)(*sql.DB, error){
	db, err := sql.Open("postgres", connURL)

	if err != nil {
		return nil, err
	}

	db.Ping()
	fmt.Println(db.Ping())

	return db, nil
}

//Load dot env and return value for key
func loadDotEnv(key string) string {
	err := godotenv.Load(".env")

	if err != nil {
		fmt.Println(err)
	}
	return os.Getenv(key)
}