package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/CodeArena-Org/codearena-ws/routers"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("unable to load env: %v", err)
	}
	PORT := os.Getenv("WS_PORT")
	// db, err := sql.Open("postgres", os.Getenv("POSTGRESS_URL"))
	// if err != nil {
	// 	log.Printf("unable to connect to db: %v", err)
	// }
	// err = db.Ping()
	// if err != nil {
	// 	log.Printf("unable to ping db: %v", err)
	// }

	if PORT == "" {
		PORT = ":5050"
	} else {
		PORT = ":" + PORT
	}
	r := routers.Router()
	log.Println("Server started on PORT ", PORT)
	log.Println(http.ListenAndServe(PORT, r))
}
