package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/yonasadiel/charon/auth"
	"github.com/yonasadiel/helios"
)

func main() {
	var name, username, password, userType string
	var user auth.User
	reader := bufio.NewReader(os.Stdin)

	helios.App.Initialize()
	helios.App.Migrate()

	fmt.Printf("Input name: ")
	name, _ = reader.ReadString('\n')
	name = strings.TrimSpace(name)

	fmt.Printf("Input username: ")
	username, _ = reader.ReadString('\n')
	username = strings.TrimSpace(username)

	fmt.Printf("Input password: ")
	password, _ = reader.ReadString('\n')
	password = strings.TrimSpace(password)

	user = auth.NewUser(name, username, password)

	fmt.Printf("Input user type [A]dmin / [O]rganize / [L]ocal / [P]articipant: ")
	userType, _ = reader.ReadString('\n')
	userType = strings.TrimSpace(userType)

	if userType == "A" {
		user.Role = "admin"
	} else if userType == "o" {
		user.Role = "organizer"
	} else if userType == "L" {
		user.Role = "local"
	} else if userType == "P" {
		user.Role = "participant"
	} else {
		fmt.Println("Unknown user type.")
		return
	}
	helios.DB.Create(&user)

	fmt.Println("Success creates user.")
}
