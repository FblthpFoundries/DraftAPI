package common

import(
	"testing"
	"net/http"
	"bytes"
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

func dummyHttp() (http.ResponseWriter, *http.Request){
	req, _ := http.NewRequest("GET","", bytes.NewReader(make([]byte, 8))) 
	return dummyWriter{} , req 
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

	rId := s.roomServer.NewRoom(dummyHttp())

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
