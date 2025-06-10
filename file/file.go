package file

import (
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/navidabedi92/NLockBox.git/encryption"
)

type Secret struct {
	Key   string
	Value string
}

var secretFilePath string
var localAppData string

func Init() {

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

func Write(secrets []Secret, path string) {
	var filePath string
	if path == "" {
		filePath = secretFilePath
	} else {
		filePath = path
	}
	file, err := os.Create(filePath)
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
		sb.WriteString(secret.Key + "\t" + secret.Value + "\n")
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
			secrets = append(secrets, Secret{Key: up[0], Value: up[1]})
		}
	}

	return secrets

}

func RenewFolders() {
	os.RemoveAll(localAppData)
	CreateFolders()
}

func Copy(path string) {
	secrets := ReadFile()
	path = path + "NLockBox-" + time.Now().Format("2006-01-02") + ".txt"
	var decrypteSecrects []Secret
	for _, secret := range secrets {
		decodedBytes, _ := base64.StdEncoding.DecodeString(secret.Value)
		decryptedArray, _ := encryption.Decrypt(decodedBytes)
		decrypted := string(decryptedArray)
		newSecret := Secret{Key: secret.Key, Value: decrypted}
		decrypteSecrects = append(decrypteSecrects, newSecret)
	}
	Write(decrypteSecrects, path)

}
