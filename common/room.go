package common

import(
	"fmt"
	"net/http"
	"github.com/google/uuid"
	"time"
	"sync"
	"strconv"
	"encoding/json"
)

//Handles finding rooms from db
type RoomServer struct{
	DBRef *DataBase
	sync.Mutex
}

func NewRoomServer(d *DataBase) *RoomServer{
	rs := RoomServer{DBRef: d}

	return &rs
}

//Representation of room
type Room struct{
	Id uuid.UUID
	Capacity int
	Expire time.Time
	Players []uuid.UUID
}

func (rs *RoomServer) NewRoom(w http.ResponseWriter, r * http.Request) {
	fmt.Println("Huzzah!")
	room := Room{ 
		Id: uuid.New(),
		Capacity: 8,
		Expire: time.Now().Add(time.Hour * 24),
		Players: make([]uuid.UUID, 8),
	}

	w.WriteHeader(http.StatusOK)
	fmt.Println(room.Id)
}

func (room *Room) AddPlayer(w http.ResponseWriter, r *http.Request){
	id, err := strconv.Atoi(r.PathValue("playerId"))

	fmt.Print(id)

	if err != nil{
		http.Error(w, "No Player Id", http.StatusBadRequest)
		return	
	}
	
	js, err := json.Marshal(room)
	if err != nil{
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("content-type", "application/json")
	w.Write(js)

}

func (rs *RoomServer) AddPlayer(w http.ResponseWriter, r *http.Request){
	
}
