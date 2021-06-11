package main

import (
	"backend/server"
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)



func main() {
	db, err := connect("postgres://krworokvzlpzdf:8dad2e8ca848991a9e1f1a307227f6d0eee2db1b9eabb761d350f3fa8a9d0f84@ec2-54-163-97-228.compute-1.amazonaws.com:5432/d158lk1cpmebhj")

	fmt.Println(db)
	fmt.Println(err)

	server.AllRooms.Init()
	server.DB = db

	http.HandleFunc("/create", server.CreateRoomRequestHandler)
	http.HandleFunc("/join", server.JoinRoomRequestHandler)
	http.HandleFunc("/get", server.GetRoomsRequestHandler)

	log.Println("Starting Server on :8000")
	http.ListenAndServe(":8000", nil)
}

func connect(connURL string)(*sql.DB, error){
	db, err := sql.Open("postgres", connURL)

	//fmt.Println(db)

	if err != nil {
		return nil, err
	}

	db.Ping()

	fmt.Println(db.Ping())

	return db, nil
}