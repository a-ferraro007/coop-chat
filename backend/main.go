package main

import (
	backendServer "backend/server"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	//	"github.com/go-oauth2/oauth2/errors"
	//	"github.com/go-oauth2/oauth2/v4/manage"
	//	"github.com/go-oauth2/oauth2/v4/models"
	//	"github.com/go-oauth2/oauth2/v4/server"
	//	"github.com/go-oauth2/oauth2/v4/store"
)



func main() {
	dbUrl := loadDotEnv("dbURL")
	db, err := connect(dbUrl)

	fmt.Println(db)
	fmt.Println(err)

	backendServer.AllRooms.Init()
	backendServer.DB = db //I need a reference to the DB connection in the signal go module

	//Authorized API Routes
	http.HandleFunc("/create", backendServer.FirebaseAuthRoute(http.HandlerFunc(backendServer.CreateRoomRequestHandler)))
	http.HandleFunc("/join", backendServer.FirebaseAuthRoute(http.HandlerFunc(backendServer.JoinRoomRequestHandler)))
	http.HandleFunc("/get", backendServer.FirebaseAuthRoute(http.HandlerFunc(backendServer.GetRoomsRequestHandler)))

	//Unauthorized API Routes
	http.HandleFunc("/create_user", backendServer.FirebaseAuthRoute(http.HandlerFunc(backendServer.RegisterUserRequestHandler)))
	//http.HandleFunc("/login", backendServer.LoginRequestHandler)

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