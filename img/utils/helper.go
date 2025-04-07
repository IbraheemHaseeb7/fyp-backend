package utils

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/IbraheemHaseeb7/pubsub"
)

var Requests map[string](chan pubsub.PubsubMessage)

func GenerateHash(data map[string]string) string {

	input := ""
	hasher := sha256.New()

	for _, val := range data {
		input += val
	}
	hasher.Write([]byte(input))
	hash := hasher.Sum(nil)

	return hex.EncodeToString(hash)
}
