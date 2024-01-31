package bot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const GeocodingURL string = "http://api.openweathermap.org/geo/1.0/direct?"
const WeatherURL string = "https://api.openweathermap.org/data/2.5/weather?"

type WeatherData struct {
	Weather []struct {
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Name string `json:"name"`
}

func getCoordinates(cityName string) (float64, float64, error) {
	// fmt.Println(cityName)
	geocodingURL := fmt.Sprintf("%sq=%s&limit=1&appid=%s", GeocodingURL, cityName, OpenWeatherToken)

	client := http.Client{Timeout: 5 * time.Second}
	response, err := client.Get(geocodingURL)
	if err != nil {
		return 0, 0, err
	}

	body, _ := io.ReadAll(response.Body)
	defer response.Body.Close()

	var data []struct {
		Lat float64 `json:"lat"`
		Lon float64 `json:"lon"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil || len(data) == 0 {
		return 0, 0, fmt.Errorf("invalid location")
	}

	return data[0].Lat, data[0].Lon, nil
}

func getCurrentWeather(message string) *discordgo.MessageSend {

	r, _ := regexp.Compile(`\S+\s*$`)
	cityName := strings.TrimSpace(r.FindString(message))

	if cityName == "" {
		return &discordgo.MessageSend{
			Content: "Ондай қала жоқ сияқты",
		}
	}

	lat, lon, err := getCoordinates(cityName)
	if err != nil {
		return &discordgo.MessageSend{
			Content: "Қала табылмады",
		}
	}

	weatherURL := fmt.Sprintf("%slat=%f&lon=%f&units=imperial&appid=%s", WeatherURL, lat, lon, OpenWeatherToken)

	// Create new HTTP client & set timeout
	client := http.Client{Timeout: 5 * time.Second}

	// Query OpenWeather API
	response, err := client.Get(weatherURL)
	if err != nil {
		return &discordgo.MessageSend{
			Content: "Ауа райы таба алмадым",
		}
	}

	// Open HTTP response body
	body, _ := io.ReadAll(response.Body)
	defer response.Body.Close()

	// Convert JSON
	var data WeatherData
	json.Unmarshal([]byte(body), &data)

	// Pull out desired weather info & Convert to string if necessary
	city := data.Name
	conditions := data.Weather[0].Description
	temperature := strconv.FormatFloat(data.Main.Temp, 'f', 2, 64)
	humidity := strconv.Itoa(data.Main.Humidity)
	wind := strconv.FormatFloat(data.Wind.Speed, 'f', 2, 64)

	// Build Discord embed response
	embed := &discordgo.MessageSend{
		Embeds: []*discordgo.MessageEmbed{{
			Type:        discordgo.EmbedTypeRich,
			Title:       "Current Weather",
			Description: "Weather for " + city,
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Conditions",
					Value:  conditions,
					Inline: true,
				},
				{
					Name:   "Temperature",
					Value:  temperature + "°F",
					Inline: true,
				},
				{
					Name:   "Humidity",
					Value:  humidity + "%",
					Inline: true,
				},
				{
					Name:   "Wind",
					Value:  wind + " mph",
					Inline: true,
				},
			},
		},
		},
	}

	return embed
}