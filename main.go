package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/navidabedi92/NLockBox.git/encryption"
	"github.com/navidabedi92/NLockBox.git/file"
)

func main() {
	godotenv.Load() // 👈 load .env file

	var secretFilePath string
	if len(os.Args) < 2 {
		log.Fatal("Incorrect Command. Usage: add|del|list [flags]")
		os.Exit(1)
	}

	localAppData := filepath.Join(os.Getenv("LOCALAPPDATA"), "NLockBox")
	secretFilePath = filepath.Join(localAppData, "secrets.txt")
	_, err := os.Stat(localAppData)
	if err != nil {
		os.Mkdir(localAppData, 0700)
		os.Create(secretFilePath)
	}

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	editCmd := flag.NewFlagSet("edit", flag.ExitOnError)
	delCmd := flag.NewFlagSet("del", flag.ExitOnError)
	_ = flag.NewFlagSet("list", flag.ExitOnError)

	addUsername := addCmd.String("username", "", "username of the secret")
	addPassword := addCmd.String("password", "", "password of the secret")

	// editUsername := editCmd.String("username", "", "username of the secret")
	// editPassword := editCmd.String("password", "", "password of the secret")

	delUsername := delCmd.String("username", "", "username of the secret")

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		// secrets := file.ReadFile(secretFilePath)

		encrypted, _ := encryption.Encrypt([]byte(*addPassword))

		file.Write(secretFilePath, *addUsername+"		"+string(encrypted))
	case "edit":
		editCmd.Parse(os.Args[2:])
		encrypted, _ := encryption.Encrypt([]byte(*addPassword))
		file.Write(secretFilePath, *addUsername+"		"+string(encrypted))
	case "del":
		delCmd.Parse(os.Args[2:])
		fmt.Println(*delUsername)
	case "list":
		secrets := file.ReadFile(secretFilePath)
		fmt.Print(secrets)
	default:
		os.Exit(1)

	}

}
