package main

import (
	"fmt"
	"monkey/repl"
	"os"
	"os/user"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	var name string
	if user.Name == "" {
		name = user.Username
	} else {
		name = user.Name
	}

	fmt.Printf("Hello %s! This is the Monkey programming language\n", name)
	fmt.Println("Feel free to type in commands")
	repl.Start(os.Stdin, os.Stdout)
}
