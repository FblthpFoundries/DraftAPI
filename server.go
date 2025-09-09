package main

import(
	"fmt"
	"net/http"
	"github.com/FblthpFoundries/DraftAPI/room"
	"github.com/FblthpFoundries/DraftAPI/player"
)

func main(){
	mux := http.NewServeMux()
	rs := room.NewRoomServer()
	ps := player.NewPlayerServer()
	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprint(w, "World")
	})


	mux.HandleFunc("POST /newRoom", rs.NewRoom)
	mux.HandleFunc("POST /register", ps.NewPlayer)

	fmt.Println("This Bitch is Serving")



	http.ListenAndServe("localhost:8081", mux)
}
