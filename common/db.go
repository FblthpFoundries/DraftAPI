package common
import (
	c "github.com/ostafen/clover/v2"
	"sync"
)

type DataBase struct{
	db *c.DB 
	sync.Mutex
}

func OpenDB() *DataBase{
	dbRef, _ := c.Open("draft-db")

	return &DataBase{db : dbRef}
}

func (base *DataBase) CloseDB(){
	base.db.Close()
}
