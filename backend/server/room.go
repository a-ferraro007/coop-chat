package server

import (
	"fmt"
	"sync"

	"github.com/gorilla/websocket"
)

//Hash map of connection pools
type RoomMap struct {
	Mutex sync.RWMutex
	Map map[string][]*Pool //Key type String, ValueType Pool Pointer
}


//Pointer Reciever function modifies RoomMap Struct
//Initializes CoopMap
func (r *RoomMap) Init(){
	fmt.Println(r)
	r.Map = make(map[string][]*Pool)
}

//func (c *CoopMap) getCoopById(id string) []Participant {
//	c.Mutex.RLock()
//	defer c.Mutex.RUnlock()

//	return c.Map[id]
//}

//Create Room generate id and push onto hash map
func (r *RoomMap) createRoom(uuid string) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	//rand.Seed(time.Now().UnixNano())
	//var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXUZ1234567890")

	//b := make([]rune, 8)

	//for i := range b {
	//	b[i] = letters[rand.Intn(len(letters))]
	//}

	//id := string(b)
	pool := NewPool()
	go pool.Start()
	r.Map[uuid] = append(r.Map[uuid], pool)
}

//insert into Room and start reading messages
func (r *RoomMap) insertIntoRoom(id string, w *websocket.Conn) {
	r.Mutex.Lock()
	defer r.Mutex.Unlock()

	fmt.Println("~~~~~~~ INSERT INTO POOL ~~~~~~~~")
	if r.Map[id] != nil {
		pool := r.Map[id][0] //Get the connection pool for the roomID
		participant := &Participant{"", w, pool} //New participant for this room
		fmt.Println("~~~~~~~ PARTICIPANT INSERT INTO POOL ~~~~~~~~")
		fmt.Println(participant)

		 fmt.Println("~~~~~~~~~~~~~~~")
		 fmt.Println("~~~~~~~ POOL INSERT INTO POOL ~~~~~~~~")
		fmt.Println(pool)
fmt.Println("~~~~~~~~~~~~~~~")
		pool.Register <- participant //Add Participant to the connection Pool
		go participant.Read(pool)
	}
}

//Delete Room by ID
func (r *RoomMap) deleteRoom(id string){
	r.Mutex.RLock()
	defer r.Mutex.Unlock()

	delete(r.Map, id)
}
