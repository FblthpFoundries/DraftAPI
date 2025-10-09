package common
import (
	"encoding/json"
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

func (base *DataBase) NewPlayer(p *Player) string{
	doc := c.NewDocumentOf(p)

	base.Lock()
	defer base.Unlock()

	id, _ := base.db.InsertOne("players",doc)

	return id
}


func (base *DataBase) CreateRoom(r *Room) string{
	doc := c.NewDocumentOf(r)

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
	room := &Room{}
	
	base.Lock()
	defer base.Unlock()

	res, _ := base.db.Query("rooms").FindById(rId)

	res.Unmarshal(room)	
	fmt.Println(room)

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

func (base *DataBase) GetPlayers(rId string) []string{
	players := &Room{}

	base.Lock()
	defer base.Unlock()

	res, _ := base.db.Query("rooms").FindById(rId)

	res.Unmarshal(players)
	fmt.Println(players)

	return players.Players
}

type cardData struct {
	scryfall_id	string		`json:"id"`
	name		string		`json:"name"`	
	rarity		string		`json:"rarity"`
	images		cardImages	`json:"image_uris"`
}

type cardImages struct{
	small	string	`json:"small"`
	normal	string	`json:"normal"`
	large	string	`json:"large"`
	png		string	`json:"png"`
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
	
	//Get scryfall page with all cards in set
	url := "https://api.scryfall.com/sets/" + set
	req, _ := http.NewRequest("GET", url, nil)
	resp, _ := http.DefaultClient.Do(req)

	body, _ := io.ReadAll(resp.Body)
		
	bodyJson := map[string]string{}
	json.Unmarshal(body, &bodyJson)
	resp.Body.Close()

	//request cards in set
	page := map[string]any{}
	var cards []*cardData
	req, _ = http.NewRequest("GET", bodyJson["search_uri"], nil)
	resp, _ = http.DefaultClient.Do(req)
	body, _ = io.ReadAll(resp.Body)

	json.Unmarshal(body, &page)
	cards = unmarshalCards(page["data"].([]any))
	resp.Body.Close()

	//depaginate subsequent pages
	for page["has_more"].(bool){
		url = page["next_page"].(string)

		req, _ = http.NewRequest("GET", url, nil)
		resp, _ = http.DefaultClient.Do(req)
		body, _ = io.ReadAll(resp.Body)
		page = map[string]any{}
		json.Unmarshal(body, &page)
		resp.Body.Close()
		
		nextCards := unmarshalCards(page["data"].([]any))

		cards = append(cards, nextCards...)
	}

	fmt.Println(cards[len(cards)-1])

	//Put all card "objects" into document array
	cardDocs := make([]*c.Document, len(cards))
	for i , card := range cards{
		cardDocs[i] = c.NewDocumentOf(card)
	}

	base.Lock()
	defer base.Unlock()

	base.db.CreateCollection(set)
	base.db.Insert(set, cardDocs...)

	
		
	//fmt.Println(page["data"].([]any)[0].(map[string]any)["name"])

	return cardDocs
}

func unmarshalCards(data []any) []*cardData{
	var cards []*cardData
	var card map[string]any
	var imageData map[string]any
	for _, c:= range data{
		card = c.(map[string]any)
		imageData = card["image_uris"].(map[string]any)
		cardRef := &cardData{
					scryfall_id : card["id"].(string),
					name		: card["name"].(string),
					rarity		: card["rarity"].(string),
					images		: cardImages{
									small: 	imageData["small"].(string),
									normal: imageData["normal"].(string),
									large:	imageData["large"].(string),
									png:	imageData["png"].(string),
								},
				}
		cards = append(cards, cardRef)
	}

	return cards
}
