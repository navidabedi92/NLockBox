package main

import (
	"encoding/base64"
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
	backupCmd := flag.NewFlagSet("backup", flag.ExitOnError)
	_ = flag.NewFlagSet("renew", flag.ExitOnError)
	_ = flag.NewFlagSet("list", flag.ExitOnError)

	addUsername := addCmd.String("username", "", "username of the secret")
	addPassword := addCmd.String("password", "", "password of the secret")

	editUsername := editCmd.String("username", "", "username of the secret")
	editPassword := editCmd.String("password", "", "password of the secret")

	delUsername := delCmd.String("username", "", "username of the secret")

	backupPath := backupCmd.String("path", "", "username of the secret")

	secrets := file.ReadFile()
	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])
		_, exist := lo.Find(secrets, func(secret file.Secret) bool {
			return secret.Username == *addUsername
		})
		if !exist {
			encrypted, _ := encryption.Encrypt([]byte(*addPassword))
			encodedPassword := base64.StdEncoding.EncodeToString(encrypted)
			newSecret := file.Secret{Username: *addUsername, Password: encodedPassword}
			secrets = append(secrets, newSecret)
			file.Write(secrets, "")
		} else {
			log.Fatal("Username Is Not Unique")
		}
	case "edit":
		editCmd.Parse(os.Args[2:])
		newUsername := *editUsername
		secret, exist := lo.Find(secrets, func(secret file.Secret) bool {
			return secret.Username == newUsername
		})
		if exist {
			secrets = lo.Reject(secrets, func(x file.Secret, index int) bool {
				return x == secret
			})
			encrypted, _ := encryption.Encrypt([]byte(*editPassword))
			encodedPassword := base64.StdEncoding.EncodeToString(encrypted)
			newSecret := file.Secret{Username: newUsername, Password: encodedPassword}
			secrets = append(secrets, newSecret)
			file.Write(secrets, "")
		} else {
			log.Fatal("Username Is Not In The List")
		}
	case "del":
		delCmd.Parse(os.Args[2:])
		secrets = lo.Reject(secrets, func(x file.Secret, index int) bool {
			return x.Username == *delUsername
		})
		file.Write(secrets, "")
	case "list":
		secrets := file.ReadFile()
		for index, secret := range secrets {
			decodedBytes, _ := base64.StdEncoding.DecodeString(secret.Password)
			decryptedArray, _ := encryption.Decrypt(decodedBytes)
			decrypted := string(decryptedArray)
			fmt.Printf("%d) Username: %s	Password: %s\n", index+1, secret.Username, decrypted)
		}
	case "renew":
		file.RenewFolders()
	case "backup":
		backupCmd.Parse(os.Args[2:])
		file.Copy(*backupPath)

	default:
		os.Exit(1)

	}

}
