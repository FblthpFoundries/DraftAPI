package common

import(
	"fmt"
	"net/http"
	"time"
	"sync"
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

	w.Header().Set("content-type", "application/json")
	js, _ := json.Marshal(map[string]string {"roomId": rId})
	w.Write(js)
	fmt.Println(rId)
}


func (rs *RoomServer) AddPlayer(w http.ResponseWriter, r *http.Request){
	
	rid := r.PathValue("roomId")

	fmt.Print(rid)

	if rid == ""{
		http.Error(w, "No room Id", http.StatusBadRequest)
		return	
	}

	exists := rs.DBRef.RoomExists(rid)
	
	if ! exists{
		http.Error(w, "Room does not exist", http.StatusBadRequest)
		return
	}
	
	pid := r.PathValue("playerId")

	fmt.Print(pid)

	if pid == ""{
		http.Error(w, "No Player Id", http.StatusBadRequest)
		return	
	}

	exists = rs.DBRef.PlayerExists(pid)
	
	if ! exists{
		http.Error(w, "Player does not exist", http.StatusBadRequest)
		return
	}

	rs.DBRef.JoinRoom(pid, rid)
	
	w.WriteHeader(http.StatusOK)
}
