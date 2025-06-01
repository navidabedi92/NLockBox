package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/navidabedi92/NLockBox.git/encryption"
)

type Secret struct {
	Username string
	Password string
}

var secretFilePath string
var localAppData string

func Init() {
	godotenv.Load() // 👈 load .env file

	localAppData = filepath.Join(os.Getenv("LOCALAPPDATA"), "NLockBox")
	secretFilePath = filepath.Join(localAppData, "secrets.txt")

	CreateFolders()
}

func CreateFolders() string {

	_, err := os.Stat(localAppData)
	if err != nil {
		os.Mkdir(localAppData, 0700)
		os.Create(secretFilePath)
	}

	return secretFilePath
}

func Write(text string) {
	file, err := os.OpenFile(secretFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// Write to the file
	if _, err := file.WriteString(text + "\n"); err != nil {
		fmt.Println("Error writing to file:", err)
	} else {
		fmt.Println("Successfully appended to file.")
	}
}

func ReadFile() []Secret {

	var secrets []Secret

	data, _ := os.ReadFile(secretFilePath)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if line != "" {
			up := strings.Split(line, "		")
			decrypted, _ := encryption.Decrypt([]byte(up[1]))
			secrets = append(secrets, Secret{Username: up[0], Password: string(decrypted)})
		}
	}

	return secrets

}

func RenewFolders() {
	os.RemoveAll(localAppData)
	CreateFolders()
}
