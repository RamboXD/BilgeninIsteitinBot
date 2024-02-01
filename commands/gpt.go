package commands

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	GPTKey       string
)

func HandleGPTCommand(s *discordgo.Session, m *discordgo.MessageCreate) {
    question := strings.TrimPrefix(m.Content, "!gpt ")

    maxTokens := 200 

    // Prepare the JSON payload
    payload := map[string]interface{}{
        "model": "gpt-3.5-turbo",
        "messages": []map[string]interface{}{
            {
                "role": "system",
                "content": "You are a helpful assistant. Please provide an answer in at most 50 words.",
            },
            {
                "role": "user",
                "content": question,
            },
        },
        "temperature": 0.7,
        "max_tokens": maxTokens,
    }
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        log.Printf("Error encoding request data: %v", err)
        s.ChannelMessageSend(m.ChannelID, "Failed to encode request.")
        return
    }
    // Create the request
    req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewReader(payloadBytes))
    if err != nil {
        log.Printf("Error creating request: %v", err)
        s.ChannelMessageSend(m.ChannelID, "Failed to create request.")
        return
    }
    req.Header.Add("Authorization", "Bearer "+GPTKey)
    req.Header.Add("Content-Type", "application/json")

    // Perform the request
    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        log.Printf("Error making request: %v", err)
        s.ChannelMessageSend(m.ChannelID, "Failed to get response.")
        return
    }
    defer resp.Body.Close()

    // Read the response
    var result map[string]interface{}
    body, err := io.ReadAll(resp.Body)
    if err != nil {
        log.Printf("Error reading response body: %v", err)
        s.ChannelMessageSend(m.ChannelID, "Failed to read response.")
        return
    }

    err = json.Unmarshal(body, &result)
	log.Printf("Raw response from OpenAI: %s\n", string(body))

    if err != nil {
        log.Printf("Error decoding response: %v", err)
        s.ChannelMessageSend(m.ChannelID, "Failed to decode response.")
        return
    }

    // Extract the text response
	if choices, found := result["choices"].([]interface{}); found && len(choices) > 0 {
		if firstChoice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := firstChoice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					s.ChannelMessageSend(m.ChannelID, content)
					return
				}
			}
		}
	}
    s.ChannelMessageSend(m.ChannelID, "Failed to extract response text.")
}
