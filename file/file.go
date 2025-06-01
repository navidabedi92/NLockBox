package file

import (
	"fmt"
	"os"
	"strings"

	"github.com/navidabedi92/NLockBox.git/encryption"
)

type Secret struct {
	Username string
	Password string
}

func Write(path string, text string) {
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
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

func ReadFile(path string) []Secret {

	var secrets []Secret

	data, _ := os.ReadFile(path)
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
