package main

import (
	"log"

	"github.com/teasec4/ecomm-go-backend/db"
	"github.com/teasec4/ecomm-go-backend/ecomm-api/handler"
	"github.com/teasec4/ecomm-go-backend/ecomm-api/server"
	"github.com/teasec4/ecomm-go-backend/ecomm-api/storer"
	"github.com/ianschenck/envflag"

)

const minSecrtKeySize = 32

func main(){
	var secretKey = envflag.String("SECRET_KEY", "01234567890123456789012345678914567891", "secret key for JWT")
	if len(*secretKey) < minSecrtKeySize{
		log.Fatalf("secret key must be at least %d characters long", minSecrtKeySize)
	}
	db, err := db.NewDatabase()
	if err != nil{
		log.Fatalf("error opening db: %v", err)
	}
	defer db.Close()
	log.Println("successfully connected to db")
	
	st := storer.NewMySQLStorer(db.GetDB())
	serv := server.NewServer(st)
	ndl := handler.NewHandler(serv, *secretKey)
	handler.RegisterRoutes(ndl)
	handler.Start(":8080")
}