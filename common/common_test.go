package common

import(
	"io"
	"testing"
	"net/http"
	"net/http/httptest"
	"github.com/google/uuid"
	"encoding/json"
)

type servers struct{
	packServer *PackServer
	roomServer *RoomServer
	playerServer *PlayerServer
	db			*DataBase
}

func setup() *servers{
	dataBase := OpenDB()
	return &servers{
		packServer: NewPackServer(dataBase),
		roomServer: NewRoomServer(dataBase),
		playerServer: NewPlayerServer(dataBase),
		db: dataBase,
	}
}

func testHttp(method string, path string) (*httptest.ResponseRecorder, *http.Request){
	return httptest.NewRecorder(), httptest.NewRequest(method, path, nil) 
}

func testHttpDefault() (*httptest.ResponseRecorder, *http.Request){
	return httptest.NewRecorder(), httptest.NewRequest("GET", "/", nil) 
}

func (s *servers) cleanup(){
	s.db.CloseDB()
}

func TestSetup(t *testing.T){
	s := setup()
	defer s.cleanup()
}


func TestNewRoom(t *testing.T){
	s := setup()
	defer s.cleanup()
	
	type jsonData struct{
		RoomId string	`json:"roomId"` 
	}
	var data jsonData
	w, req := testHttpDefault()
	s.roomServer.NewRoom(w, req)
	r := w.Result()
	body, _ := io.ReadAll(r.Body)
	json.Unmarshal(body, &data)


	t.Run("RoomExists", func(t *testing.T){
		exists := s.db.RoomExists(data.RoomId)
		if !exists{
			t.Error("Room not created in database")
		}
	})

	t.Run("RoomDoesNotExist", func(t *testing.T){
		exists := s.db.RoomExists(uuid.New().String())
		if exists{
			t.Error("Database has garbage room id")
		}
	})
}

func TestNewPlayer(t *testing.T){
	s := setup()
	defer s.cleanup()

	var data struct{
		PlayerId string `json:"playerId"`
	}
	w, req := testHttpDefault()
	s.playerServer.NewPlayer(w, req)

	r := w.Result()
	json.NewDecoder(r.Body).Decode(&data)

	t.Run("PlayerExists", func(t *testing.T){
		exists := s.db.PlayerExists(data.PlayerId)
		if !exists{
			t.Error("Player not created in database")
		}
	})
	t.Run("PlayerDoesNotExists", func(t *testing.T){
		exists := s.db.PlayerExists(uuid.New().String())
		if exists{
			t.Error("Database has garbage player id")
		}
	})
}

func TestJoinRoom(t *testing.T){
	s := setup()
	defer s.cleanup()
	
	var pData struct{
		PlayerId string `json:"playerId"`
	}
	var rData struct{
		RoomId string `json:"roomId"`
	}
	w, req := testHttpDefault()
	s.playerServer.NewPlayer(w, req)
	r := w.Result()
	json.NewDecoder(r.Body).Decode(&pData)

	w, req = testHttpDefault()
	s.roomServer.NewRoom(w, req)
	r = w.Result()
	json.NewDecoder(r.Body).Decode(&rData)

	s.roomServer.AddPlayer(testHttp("PUT", "/register?rId=" + rData.RoomId + "&pId=" + pData.PlayerId))

	players := s.db.GetPlayers(rData.RoomId)

	if players[0] != pData.PlayerId{
		t.Error("Player not added to database")
	}
}

func TestGetSet(t *testing.T){
	s := setup()
	defer s.cleanup()
	
	t.Run("GetsNewSet", func(t *testing.T){
		cards := s.db.GetSet("eoe")
		if cards == nil{
			t.Error("failed to get new set")
		}
	})

	t.Run("RetrieveSet", func(t *testing.T){
		cards := s.db.GetSet("eoe")
		if cards == nil{
			t.Error("Failed to retrieve set")
		}
	})
}
