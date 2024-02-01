package commands

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

type TranslateRequest struct {
	Q      []string `json:"q"`
	Target string   `json:"target"`
}

type TranslateResponse struct {
	Data struct {
		Translations []struct {
			TranslatedText string `json:"translatedText"`
		} `json:"translations"`
	} `json:"data"`
}

// Prepare map for language name convert to needed code
var LanguageISOCodeMap = map[string]string{
	"english":    "en",
	"russian":    "ru",
	"spanish":    "es",
	"french":     "fr",
	"german":     "de",
	"chinese":    "zh",
	"japanese":   "ja",
	"korean":     "ko",
}

func TranslateText(message string, apiKey string) *discordgo.MessageSend {
	// Regex to match the command format: !translate language "text"
	re := regexp.MustCompile(`!translate\s+(\w+)\s+"([^"]+)"`)
	matches := re.FindStringSubmatch(message)

	if len(matches) < 3 {
		return &discordgo.MessageSend{Content: "Command format error. Use: !translate language \"text\"."}
	}

	languageName := strings.ToLower(matches[1])
	textToTranslate := matches[2]

	// Look up the code for the specified language
	targetLanguage, ok := LanguageISOCodeMap[languageName]
	if !ok {
		return &discordgo.MessageSend{Content: fmt.Sprintf("Unsupported language: %s.", languageName)}
	}

	// Prepare the request body
	requestBody := TranslateRequest{
		Q:      []string{textToTranslate},
		Target: targetLanguage,
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return &discordgo.MessageSend{Content: "Failed to prepare the translation request."}
	}

	// Construct the request URL with the API key
	url := fmt.Sprintf("https://translation.googleapis.com/language/translate/v2?key=%s", apiKey)

	// Make the POST request
	resp, err := http.Post(url, "application/json", bytes.NewBuffer(requestBodyBytes))
	if err != nil {
		return &discordgo.MessageSend{Content: "Failed to connect to the translation service."}
	}
	defer resp.Body.Close()

	// Read and parse the response
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return &discordgo.MessageSend{Content: "Failed to read the translation response."}
	}

	var translateResponse TranslateResponse
	err = json.Unmarshal(responseBody, &translateResponse)
	if err != nil {
		return &discordgo.MessageSend{Content: "Failed to parse the translation response."}
	}

	// First translation is what we want
	if len(translateResponse.Data.Translations) > 0 {
		translatedText := translateResponse.Data.Translations[0].TranslatedText
		return &discordgo.MessageSend{Content: translatedText}
	}

	return &discordgo.MessageSend{Content: "Translation was unsuccessful."}
}
