package common
import (
	"fmt"
	c "github.com/ostafen/clover"
	"sync"
)

type DataBase struct{
	db *c.DB 
	sync.Mutex
}

func OpenDB() *DataBase{
	fmt.Println("Starting database")
	dbRef, _ := c.Open("draft-db")

	dbRef.CreateCollection("rooms")

	dbRef.CreateCollection("players")

	return &DataBase{db : dbRef}
}

func (base *DataBase) CloseDB(){
	base.db.Close()
}



func (base *DataBase) CreateRoom(r *Room) string{
	doc := c.NewDocumentOf(*r)

	base.Lock()
	defer base.Unlock()

	id, _ := base.db.InsertOne("rooms", doc)
	
	return id
}
