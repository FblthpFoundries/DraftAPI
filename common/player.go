package common

import (
	"fmt"
	"net/http"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"sync"
	"encoding/json"
)

var upgrader = websocket.Upgrader{}

type Player struct{
	Id uuid.UUID
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
	p := Player{
		Id: uuid.New(),
	}
	
	fmt.Println(p.Id)
	js, err := json.Marshal(map[string]uuid.UUID{"Id":p.Id})

	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
