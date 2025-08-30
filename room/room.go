package room

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
	Cache  map[uuid.UUID]*Room
	sync.Mutex
}

func NewRoomServer() *RoomServer{
	rs := RoomServer{Cache: make(map[uuid.UUID]*Room)}

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

	rs.Cache[room.Id] = &room 	

	w.WriteHeader(http.StatusOK)
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
