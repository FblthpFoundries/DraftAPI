package common
import (
	"fmt"
	c "github.com/ostafen/clover"
	"sync"
	"net/http"
	"io"
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

func (base *DataBase) NewPlayer() string{
	doc := c.NewDocument()

	base.Lock()
	defer base.Unlock()

	id, _ := base.db.InsertOne("players",doc)

	return id
}


func (base *DataBase) CreateRoom(r *Room) string{
	doc := c.NewDocumentOf(*r)

	base.Lock()
	defer base.Unlock()

	id, _ := base.db.InsertOne("rooms", doc)
	
	return id
}

func (base *DataBase) PlayerExists(pId string) bool{
	base.Lock()
	defer base.Unlock()

	exists, _ := base.db.Query("players").FindById(pId)

	return exists != nil
}

func (base *DataBase) RoomExists(rId string) bool{
	base.Lock()
	defer base.Unlock()

	exists, _ := base.db.Query("rooms").FindById(rId)

	return exists != nil
}


func (base *DataBase) JoinRoom(pId string, rId string){
	room := &struct{
		Players []string
	}{}
	
	base.Lock()
	defer base.Unlock()

	res, _ := base.db.Query("rooms").FindById(rId)

	res.Unmarshal(room)	

	newPlayers := make([]string, 8)
	insertIdx := 0

	//copy over existing players
	for idx, p := range room.Players{
		if p == ""{
			break
		}

		newPlayers[idx] = p
		insertIdx += 1
	}

	//add new player
	newPlayers[insertIdx] = pId

	//update document
	updates := make(map[string]interface{})
	updates["Players"] = newPlayers

	base.db.Query("rooms").UpdateById(rId, updates)

}

func (base *DataBase) GetSet(set string) []*c.Document {
	base.Lock()
	res, _ := base.db.HasCollection(set)
	base.Unlock()

	//Already retrieved set info so return document ref
	if res{
		cards, _ := base.db.Query(set).FindAll()
		return cards
	}

	url := "https://api.scryfall.com/sets/" + set
	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := http.DefaultClient.Do(req)

	body, _ := io.ReadAll(resp.Body)
	fmt.Println(string(body[:]))

	return make([]*c.Document, 8)
}
