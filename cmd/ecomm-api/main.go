package main

import (
	"log"

	"github.com/teasec4/ecomm-go-backend/db"
)

func main(){
	db, err := db.NewDatabase()
	if err != nil{
		log.Fatalf("error opening db: %v", err)
	}
	defer db.Close()
	log.Println("successfully connected to db")
}