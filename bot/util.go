package bot

import (
	"os"
	"regexp"
)

func getenv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

func stripBotName(text string) string {
	re := regexp.MustCompile(`^<@.*> *`)
	return re.ReplaceAllString(text, "")
}
