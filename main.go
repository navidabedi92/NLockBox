package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/navidabedi92/NLockBox.git/encryption"
	"github.com/navidabedi92/NLockBox.git/file"
	"github.com/samber/lo"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Incorrect Command. Usage: add|del|list [flags]")
		os.Exit(1)
	}
	file.Init()

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	editCmd := flag.NewFlagSet("edit", flag.ExitOnError)
	delCmd := flag.NewFlagSet("del", flag.ExitOnError)
	_ = flag.NewFlagSet("renew", flag.ExitOnError)
	_ = flag.NewFlagSet("list", flag.ExitOnError)

	addUsername := addCmd.String("username", "", "username of the secret")
	addPassword := addCmd.String("password", "", "password of the secret")

	// editUsername := editCmd.String("username", "", "username of the secret")
	// editPassword := editCmd.String("password", "", "password of the secret")

	delUsername := delCmd.String("username", "", "username of the secret")

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])

		secrets := file.ReadFile()

		_, exist := lo.Find(secrets, func(secret file.Secret) bool {
			return secret.Username == *addUsername
		})
		if !exist {
			encrypted, _ := encryption.Encrypt([]byte(*addPassword))
			file.Write(*addUsername + "		" + string(encrypted))
		}
	case "edit":
		editCmd.Parse(os.Args[2:])
		encrypted, _ := encryption.Encrypt([]byte(*addPassword))
		file.Write(*addUsername + "		" + string(encrypted))
	case "del":
		delCmd.Parse(os.Args[2:])
		fmt.Println(*delUsername)
	case "list":
		secrets := file.ReadFile()
		fmt.Print(secrets)
	case "renew":
		file.RenewFolders()

	default:
		os.Exit(1)

	}

}
