package file

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
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

func Write(secrets []Secret) {
	file, err := os.Create(secretFilePath)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer func() {
		if err := file.Close(); err != nil {
			fmt.Println("Error closing file:", err)
		}
	}()

	// Use strings.Builder for efficient concatenation
	var sb strings.Builder
	for _, secret := range secrets {
		sb.WriteString(secret.Username + "\t" + secret.Password + "\n")
	}

	// Write to file
	if _, err := file.WriteString(sb.String()); err != nil {
		fmt.Println("Error writing to file:", err)
	} else {
		fmt.Println("Success: Secrets saved")
	}
}

func ReadFile() []Secret {

	var secrets []Secret

	data, _ := os.ReadFile(secretFilePath)
	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		if line != "" {
			up := strings.Split(line, "	")
			secrets = append(secrets, Secret{Username: up[0], Password: up[1]})
		}
	}

	return secrets

}

func RenewFolders() {
	os.RemoveAll(localAppData)
	CreateFolders()
}
