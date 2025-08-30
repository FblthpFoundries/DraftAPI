package main

import(
	"fmt"
	"net/http"
	"github.com/FblthpFoundries/DraftAPI/room"
)

func main(){
	mux := http.NewServeMux()
	rs := room.NewRoomServer()
	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprint(w, "World")
	})


	mux.HandleFunc("POST /newRoom", rs.NewRoom)

	fmt.Println("This Bitch is Serving")



	http.ListenAndServe("localhost:8081", mux)
}
