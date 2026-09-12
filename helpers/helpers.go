package helpers

import (
	"math/rand/v2"
	"strings"
)

var letters = "abcdifghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"

const idLen int = 20

func GenerateID() string {
	var id string
	lettersArr := strings.Split(letters, "")
	for i := 1; i <= idLen; i++ {
		var rand = rand.IntN(len(lettersArr) - 1)
		id += lettersArr[rand]
	}

	return id
}
