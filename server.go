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

	roomServer := common.NewRoomServer(dataBase)
	playerServer := common.NewPlayerServer(dataBase)
	packServer := common.NewPackServer(dataBase)
	

	mux.HandleFunc("GET /hello", func(w http.ResponseWriter, r *http.Request){
		fmt.Fprint(w, "World")
	})


	mux.HandleFunc("POST /newRoom", roomServer.NewRoom)
	mux.HandleFunc("POST /newPlayer", playerServer.NewPlayer)
	mux.HandleFunc("POST /register", roomServer.AddPlayer)
	mux.HandleFunc("GET /getSet", packServer.GeneratePacks)

	fmt.Println("This Bitch is Serving")



	http.ListenAndServe("localhost:8081", mux)
}
