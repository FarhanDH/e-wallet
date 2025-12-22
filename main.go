package main

import (
	"ewallet/config"
	"ewallet/internal/repository"
	"fmt"
)

func main() {
	db, err := config.ConnectDB()
	if err != nil {
		panic(err)
	}
	defer db.Close()
	repo := repository.NewUserRepository(db)
	user, err := repo.GetUser(1)
	fmt.Println(user)
	user, err = repo.GetUser(2)
	fmt.Println(user)
	err = repo.Transfer(1, 2, 5000)
	if err != nil {
		fmt.Println("Transfer Failed", err)
		panic(err)
	} else {
		fmt.Println("Transfer Success")
	}
	user, err = repo.GetUser(1)
	fmt.Println(user)
	user, err = repo.GetUser(2)
	fmt.Println(user)

}
