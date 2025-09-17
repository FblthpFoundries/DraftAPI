package main

import(
	"fmt"
	"net/http"
	"github.com/FblthpFoundries/DraftAPI/common"
)

func main(){
	fmt.Println("Here?")
	dataBase := common.OpenDB()
	defer dataBase.CloseDB()
	
	mux := http.NewServeMux()

	rs := common.NewRoomServer(dataBase)
	ps := common.NewPlayerServer(dataBase)

	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprint(w, "World")
	})


	mux.HandleFunc("POST /newRoom", rs.NewRoom)
	mux.HandleFunc("POST /newPlayer", ps.NewPlayer)
	mux.HandleFunc("POST /register/{roomId}/{playerId}", rs.AddPlayer)

	fmt.Println("This Bitch is Serving")



	http.ListenAndServe("localhost:8081", mux)
}
