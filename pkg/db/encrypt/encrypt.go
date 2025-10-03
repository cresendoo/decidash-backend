package encrypt

import (
	"fmt"
	"os"
	"strings"

	"github.com/cresendoo/decidash-backend/pkg/utils"
)

var secret string

type Cipher interface {
	Encrypt(int64, []byte) (string, error)
	Decrypt(int64, string) ([]byte, error)
}

func init() {
	s := strings.TrimSpace(os.Getenv("DATABASE_SECRET"))
	if s != "" {
		secret = string(s)
	} else {
		if utils.IsProductionPhase() {
			panic("DATABASE_SECRET environment variable must be set!!!")
		}
		fmt.Println("!!!! DATABASE_SECRET not set. Use default key !!!!")
		secret = "CK4EaJE7FuY6LQpcMn53Rw=="
	}
}
