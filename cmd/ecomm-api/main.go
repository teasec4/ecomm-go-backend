package main

import (
	"log"

	"github.com/teasec4/ecomm-go-backend/db"
	"github.com/teasec4/ecomm-go-backend/ecomm-api/handler"
	"github.com/teasec4/ecomm-go-backend/ecomm-api/server"
	"github.com/teasec4/ecomm-go-backend/ecomm-api/storer"
)

func main(){
	db, err := db.NewDatabase()
	if err != nil{
		log.Fatalf("error opening db: %v", err)
	}
	defer db.Close()
	log.Println("successfully connected to db")
	
	st := storer.NewMySQLStorer(db.GetDB())
	serv := server.NewServer(st)
	ndl := handler.NewHandler(serv)
	handler.RegisterRoutes(ndl)
	handler.Start(":8080")
}