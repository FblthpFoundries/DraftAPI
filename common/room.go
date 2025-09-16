package common

import(
	"fmt"
	"net/http"
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
	Capacity int
	Expire time.Time
	Players []string
}

func (rs *RoomServer) NewRoom(w http.ResponseWriter, r * http.Request) {
	fmt.Println("Huzzah!")
	room := Room{ 
		Capacity: 8,
		Expire: time.Now().Add(time.Hour * 24),
		Players: make([]string, 8),
	}

	rId := rs.DBRef.CreateRoom(&room)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	js, _ := json.Marshal(map[string]string {"roomId": rId})
	w.Write(js)
	fmt.Println(rId)
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
