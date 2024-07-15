package quote

import (
	"encoding/json"
	"math/rand"
	"os"

	"main/internal/logger"
)

type quotes struct {
	Text string `json:"text"`
}

var parsedQuotes = []quotes{}

func init() {
	if err := json.Unmarshal(quotesJSON, &parsedQuotes); err != nil {
		logger.Stderr.Error("failed to unmarshal quotes", logger.ErrAttr(err))

		os.Exit(1)
	}
}

func GetRandom() string {
	return parsedQuotes[rand.Intn(len(parsedQuotes))].Text
}
