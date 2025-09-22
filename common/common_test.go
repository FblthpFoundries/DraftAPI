package common

import(
	"testing"
	"net/http"
	"net/http/httptest"
	"github.com/google/uuid"
)

type servers struct{
	packServer *PackServer
	roomServer *RoomServer
	playerServer *PlayerServer
	db			*DataBase
}

type dummyWriter struct{
}

func (w dummyWriter) Header() http.Header{
	return http.Header{}
}

func (w dummyWriter) Write(data []byte) (int, error){
	return 0, nil

}

func (w dummyWriter) WriteHeader(status int){
	return
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

func testHttp(method string, path string) (http.ResponseWriter, *http.Request){
	return httptest.NewRecorder(), httptest.NewRequest(method, path, nil) 
}

func testHttpDefault() (http.ResponseWriter, *http.Request){
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

	rId := s.roomServer.NewRoom(testHttpDefault())

	t.Run("RoomExists", func(t *testing.T){
		exists := s.db.RoomExists(rId)
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

	pId := s.playerServer.NewPlayer(testHttpDefault())

	t.Run("PlayerExists", func(t *testing.T){
		exists := s.db.PlayerExists(pId)
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

	pId := s.playerServer.NewPlayer(testHttpDefault())
	rId := s.roomServer.NewRoom(testHttpDefault())

	s.roomServer.AddPlayer(testHttp("PUT", "/register?rId=" + rId + "&pId=" + pId))

	players := s.db.GetPlayers(rId)

	if players[0] != pId{
		t.Error("Player not added to database")
	}
}
