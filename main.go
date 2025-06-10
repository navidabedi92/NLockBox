package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/joho/godotenv"
	"github.com/navidabedi92/NLockBox.git/encryption"
	"github.com/navidabedi92/NLockBox.git/file"
	"github.com/samber/lo"
)

func main() {

	if len(os.Args) < 2 {
		log.Fatal("Incorrect Command. Usage: add|del|list [flags]")
		os.Exit(1)
	}

	path, err := os.Executable()
	if err != nil {
		panic(err)
	}
	dir := filepath.Dir(path)
	envPath := filepath.Join(dir, ".env")
	godotenv.Load(envPath) // 👈 load .env file

	file.Init()

	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	editCmd := flag.NewFlagSet("edit", flag.ExitOnError)
	delCmd := flag.NewFlagSet("del", flag.ExitOnError)
	findCmd := flag.NewFlagSet("find", flag.ExitOnError)
	backupCmd := flag.NewFlagSet("backup", flag.ExitOnError)
	_ = flag.NewFlagSet("renew", flag.ExitOnError)
	_ = flag.NewFlagSet("list", flag.ExitOnError)

	addkey := addCmd.String("key", "", "key of the secret")
	addvalue := addCmd.String("value", "", "value of the secret")

	editkey := editCmd.String("key", "", "key of the secret")
	editvalue := editCmd.String("value", "", "value of the secret")

	findkey := findCmd.String("key", "", "key of the secret")

	delkey := delCmd.String("key", "", "key of the secret")

	backupPath := backupCmd.String("path", "", "key of the secret")

	secrets := file.ReadFile()

	switch os.Args[1] {
	case "add":
		addCmd.Parse(os.Args[2:])

		_, exist := lo.Find(secrets, func(secret file.Secret) bool {
			return secret.Key == *addkey
		})
		if !exist {
			encrypted, _ := encryption.Encrypt([]byte(*addvalue))
			encodedvalue := base64.StdEncoding.EncodeToString(encrypted)
			newSecret := file.Secret{Key: *addkey, Value: encodedvalue}

			secrets = append(secrets, newSecret)
			file.Write(secrets, "")
		} else {
			log.Fatal("key Is Not Unique")
		}
	case "edit":
		editCmd.Parse(os.Args[2:])
		newkey := *editkey
		secret, exist := lo.Find(secrets, func(secret file.Secret) bool {
			return secret.Key == newkey
		})
		if exist {
			secrets = lo.Reject(secrets, func(x file.Secret, index int) bool {
				return x == secret
			})
			encrypted, _ := encryption.Encrypt([]byte(*editvalue))
			encodedvalue := base64.StdEncoding.EncodeToString(encrypted)
			newSecret := file.Secret{Key: newkey, Value: encodedvalue}
			secrets = append(secrets, newSecret)
			file.Write(secrets, "")
		} else {
			log.Fatal("key Is Not In The List")
		}
	case "del":
		delCmd.Parse(os.Args[2:])
		secrets = lo.Reject(secrets, func(x file.Secret, index int) bool {
			return x.Key == *delkey
		})
		file.Write(secrets, "")
	case "find":
		findCmd.Parse(os.Args[2:])
		secret, exist := lo.Find(secrets, func(secret file.Secret) bool {
			return secret.Key == *findkey
		})

		if exist {
			decodedBytes, _ := base64.StdEncoding.DecodeString(secret.Value)
			decryptedArray, _ := encryption.Decrypt(decodedBytes)
			decrypted := string(decryptedArray)
			fmt.Print(decrypted)
		}
	case "list":
		secrets := file.ReadFile()
		for index, secret := range secrets {
			decodedBytes, _ := base64.StdEncoding.DecodeString(secret.Value)
			decryptedArray, _ := encryption.Decrypt(decodedBytes)
			decrypted := string(decryptedArray)
			fmt.Printf("%d) key: %s	value: %s\n", index+1, secret.Key, decrypted)
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
