package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Incorrect Command. Usage: add|del|list [flags]")
		os.Exit(1)
	}

	localAppData := filepath.Join(os.Getenv("LOCALAPPDATA"), "NLockBox")

	_, err := os.Stat(localAppData)

	if err != nil {
		os.Mkdir(localAppData, 0700)
		os.Create(filepath.Join(localAppData, "secrets.txt"))
	}

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	delCmd := flag.NewFlagSet("del", flag.ExitOnError)
	_ = flag.NewFlagSet("list", flag.ExitOnError)

	addUsername := addCmd.String("username", "", "username of the secret")
	addPassword := addCmd.String("password", "", "password of the secret")

	delUsername := delCmd.String("username", "", "username of the secret")

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		fmt.Println(*addUsername)
		fmt.Println(*addPassword)
	case "del":
		delCmd.Parse(os.Args[2:])
		fmt.Println(*delUsername)
	case "list":
		fmt.Println("list")
	default:
		os.Exit(1)

	}

}
