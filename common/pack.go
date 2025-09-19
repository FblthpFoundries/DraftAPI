package common

import (
	"sync"
	"net/http"
)

type Pack struct{
	Cards []string
	sync.Mutex
}

type PackServer struct{
	dbRef *DataBase
}

func NewPackServer(db *DataBase) *PackServer{
	return &PackServer{dbRef : db}
}

func (ps *PackServer) GeneratePacks(w http.ResponseWriter, r * http.Request) {
	ps.dbRef.GetSet("eoe")

}

