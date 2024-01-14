package main

import (
	"context"
	"fmt"
	"github.com/codewithtoucans/goweb/models"
	"os"
)

func main() {
	config := models.DefaultPostgresConfig()
	conn, err := models.Open(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	err = conn.Ping(context.Background())
	if err != nil {
		return
	}

	//_, err = conn.Exec(context.Background(), `
	//		CREATE TABLE users (
	//			id SERIAL primary key,
	//			email TEXT UNIQUE not null,
	//			password_hash TEXT not null
	//		)
	//`)
	//
	//_, err = conn.Exec(context.Background(), `
	//	create table sessions (
	//		id serial primary key,
	//		user_id int unique,
	//		token_hash text unique not null,
	//		foreign key (user_id) references users(id) on delete cascade
	//	)
	//`)

	//us := models.SessionService{DB: conn}
	//user, err := us.Create(123456)
	//if err != nil {
	//	fmt.Println(err)
	//	fmt.Println("create users table is error")
	//	return
	//}
	//fmt.Printf("user is %+v\n", user)
	//rows, err := conn.Query(context.Background(), "drop table sessions")
	//if err != nil {
	//	log.Fatal(err)
	//}
	//defer rows.Close()

	//for rows.Next() {
	//	var id any
	//	var userID any
	//	var tokenHash any
	//	if err := rows.Scan(&id, &userID, &tokenHash); err != nil {
	//		log.Fatal(err)
	//	}
	//	fmt.Println(id, userID, tokenHash)
	//}

	fmt.Println("connect success")
}
