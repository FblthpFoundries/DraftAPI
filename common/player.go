package common

import (
	"fmt"
	"net/http"
	"sync"
	"encoding/json"
)


type Player struct{
	IsBot bool
	Cards []string
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


func (ps *PlayerServer) NewPlayer(w http.ResponseWriter, r * http.Request) string{
	p := &Player{
		IsBot: false,
		Cards: make([]string, 45),
	}
	pId := ps.DBRef.NewPlayer(p)

	fmt.Println(pId)
	js, err := json.Marshal(map[string]string{"playerId":pId})

	if err != nil{
		http.Error(w, err.Error(), http.StatusBadRequest)
		return "BAD"
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
	return pId
}
