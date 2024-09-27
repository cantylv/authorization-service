package functions

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"strings"

	mc "github.com/cantylv/authorization-service/internal/utils/myconstants"
	"golang.org/x/crypto/argon2"
)

// HashData хэширует данные используя алгоритм Argon2. Принимает на вход данные для хэширования, соль
// и возвращает хэш в виде строки шестнадцетеричных цифр.
func HashData(payload, salt []byte) string {
	saltCopy := make([]byte, len(salt))
	copy(saltCopy, salt)
	hashedPassword := argon2.IDKey(payload, salt, mc.HashTime, mc.HashMemory, mc.HashThreads, mc.HashKeylen)
	return hex.EncodeToString(hashedPassword)
}

func GetHashedPassword(pwdPass string) (string, error) {
	salt, err := generateNewSalt()
	if err != nil {
		return "", err
	}
	hashedPwd := HashData([]byte(pwdPass), []byte(salt))
	return fmt.Sprintf("%s.%s", hashedPwd, salt), nil
}

func generateNewSalt() (string, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}
	return hex.EncodeToString(salt), nil
}

func IsPasswordsEqual(pwdPass, pwdDB string) bool {
	saltDB := strings.Split(pwdDB, ".")[1]
	hash := HashData([]byte(pwdPass), []byte(saltDB))

	return fmt.Sprintf("%s.%s", hash, saltDB) == pwdDB
}
