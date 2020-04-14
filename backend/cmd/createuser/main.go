package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/yonasadiel/charon/backend/auth"
	"github.com/yonasadiel/helios"
)

func main() {
	var name, username, password, userRoleString string
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

	fmt.Printf("Input user type (admin / organizer / local / participant): ")
	userRoleString, _ = reader.ReadString('\n')
	userRoleString = strings.ToLower(strings.TrimSpace(userRoleString))

	var userRole uint
	if userRoleString == "admin" {
		userRole = auth.UserRoleAdmin
	} else if userRoleString == "organizer" {
		userRole = auth.UserRoleOrganizer
	} else if userRoleString == "local" {
		userRole = auth.UserRoleLocal
	} else if userRoleString == "participant" {
		userRole = auth.UserRoleParticipant
	} else {
		fmt.Println("Unknown user type.")
		return
	}

	auth.UserFactorySaved(auth.User{Name: name, Username: username, Password: password, Role: userRole})

	fmt.Println("Success creates user.")
}
