package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

const GeocodingURL string = "http://api.openweathermap.org/geo/1.0/direct?"
const WeatherURL string = "https://api.openweathermap.org/data/2.5/weather?"
var (
	OpenWeatherToken       string
)

type WeatherDetails struct {
    Temp      float64 `json:"temp"`
    FeelsLike float64 `json:"feels_like"`
    TempMin   float64 `json:"temp_min"`
    TempMax   float64 `json:"temp_max"`
    Pressure  int     `json:"pressure"`
    Humidity  int     `json:"humidity"`
}

type Wind struct {
    Speed float64 `json:"speed"`
    Deg   int     `json:"deg"`
}

type Clouds struct {
    All int `json:"all"`
}

type Precipitation struct {
    OneHour   float64 `json:"1h,omitempty"`
    ThreeHour float64 `json:"3h,omitempty"`
}

type Sys struct {
    Sunrise int64 `json:"sunrise"`
    Sunset  int64 `json:"sunset"`
}

type WeatherData struct {
    Weather       []struct {
        Description string `json:"description"`
    } `json:"weather"`
    Main          WeatherDetails `json:"main"`
    Wind          Wind           `json:"wind"`
    Clouds        Clouds         `json:"clouds"`
    Rain          *Precipitation `json:"rain,omitempty"`
    Snow          *Precipitation `json:"snow,omitempty"`
    Name          string         `json:"name"`
    Visibility    int            `json:"visibility"`
    Sys           Sys            `json:"sys"`
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

func GetCurrentWeather(message string) *discordgo.MessageSend {

	r, _ := regexp.Compile(`\S+\s*$`)
	cityName := strings.TrimSpace(r.FindString(message))

	if cityName == "" {
		return &discordgo.MessageSend{
			Content: "Write me a proper city name",
		}
	}

	lat, lon, err := getCoordinates(cityName)
	if err != nil {
		return &discordgo.MessageSend{
			Content: "No such city",
		}
	}

	weatherURL := fmt.Sprintf("%slat=%f&lon=%f&units=Metric&appid=%s", WeatherURL, lat, lon, OpenWeatherToken)

	// Create new HTTP client & set timeout
	client := http.Client{Timeout: 5 * time.Second}

	// Query OpenWeather API
	response, err := client.Get(weatherURL)
	if err != nil {
		return &discordgo.MessageSend{
			Content: "Weather not found",
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
    sunrise := time.Unix(data.Sys.Sunrise, 0).Format("15:04")
    sunset := time.Unix(data.Sys.Sunset, 0).Format("15:04")

    // Build Discord embed response 
    embed := &discordgo.MessageEmbed{
        Title:       "Current Weather in " + city,
        Description: fmt.Sprintf("**%s**\n*%s*", conditions, city),
        Color:       0x00ff00, 
        Fields: []*discordgo.MessageEmbedField{
            {Name: "Temperature", Value: fmt.Sprintf("%0.2f°C (Feels like: %0.2f°C)", data.Main.Temp, data.Main.FeelsLike), Inline: true},
            {Name: "Min/Max Temperature", Value: fmt.Sprintf("Min: %0.2f°C, Max: %0.2f°C", data.Main.TempMin, data.Main.TempMax), Inline: true},
            {Name: "Pressure", Value: fmt.Sprintf("%d hPa", data.Main.Pressure), Inline: true},
            {Name: "Humidity", Value: fmt.Sprintf("%d%%", data.Main.Humidity), Inline: true},
            {Name: "Wind", Value: fmt.Sprintf("%0.2f m/s at %d°", data.Wind.Speed, data.Wind.Deg), Inline: true},
            {Name: "Cloudiness", Value: fmt.Sprintf("%d%%", data.Clouds.All), Inline: true},
            {Name: "Visibility", Value: fmt.Sprintf("%d meters", data.Visibility), Inline: true},
            {Name: "Sunrise", Value: sunrise, Inline: true},
            {Name: "Sunset", Value: sunset, Inline: true},
        },
    }

    if data.Rain != nil {
        embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{Name: "Rain (last hour)", Value: fmt.Sprintf("%0.2f mm", data.Rain.OneHour), Inline: true})
    }
    if data.Snow != nil {
        embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{Name: "Snow (last hour)", Value: fmt.Sprintf("%0.2f mm", data.Snow.OneHour), Inline: true})
    }

    return &discordgo.MessageSend{Embeds: []*discordgo.MessageEmbed{embed}}
}