package player

import(
	"fmt"
	"net/http"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	//"time"
	"sync"
	//"strconv"
	"encoding/json"
)

var upgrader = websocket.Upgrader{}

type Player struct{
	Id uuid.UUID
	Sock websocket.Conn
}

type PlayerServer struct{
	Cache map[uuid.UUID]*Player
	sync.Mutex	
}

func NewPlayerServer() *PlayerServer{
	ps := PlayerServer{
		Cache : make(map[uuid.UUID]*Player),
	}
	
	return &ps
}


func (ps *PlayerServer) NewPlayer(w http.ResponseWriter, r * http.Request){
	p := Player{
		Id: uuid.New(),
	}
	ps.Lock()
	defer ps.Unlock()
	ps.Cache[p.Id] = &p
	
	fmt.Println(p.Id)
	js, err := json.Marshal(map[string]uuid.UUID{"Id":p.Id})

	if err != nil{
		delete(ps.Cache, p.Id)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
