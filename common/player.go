package common

import (
	"fmt"
	"net/http"
	"github.com/gorilla/websocket"
	"sync"
	"encoding/json"
)

var upgrader = websocket.Upgrader{}

type Player struct{
	Id string 
	Sock websocket.Conn
}

type PlayerServer struct{
	DBRef *DataBase
	sync.Mutex	
}

func NewPlayerServer(d *DataBase) *PlayerServer{
	ps := PlayerServer{
		DBRef: d,
	}
	
	return &ps
}


func (ps *PlayerServer) NewPlayer(w http.ResponseWriter, r * http.Request){
	pId := ps.DBRef.NewPlayer()

	fmt.Println(pId)
	js, err := json.Marshal(map[string]string{"playerId":pId})

	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
