package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"slices"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/valyala/fasthttp"
)

type Product struct {
	Title string `json:"title"`
	Price struct {
		Amount       string `json:"amount"`
		CurrencyCode string `json:"currency_code"`
	} `json:"price"`
	URL            string    `json:"url"`
	ImageURL       ItemPhoto `json:"photo"`
	BrandTitle     string    `json:"brand_title"`
	SizeTitle      string    `json:"size_title"`
	Status         string    `json:"status"`
	TotalItemPrice struct {
		Amount       string `json:"amount"`
		CurrencyCode string `json:"currency_code"`
	} `json:"total_item_price"`
}

type Response struct {
	Items []Product `json:"items"`
}

type ItemPhoto struct {
	DominantColor string `json:"dominant_color"`
	URL           string `json:"url"`
}

var foundTotal int

func getAccessToken(country string) string {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Error creating cookie jar: %v", err)
	}

	client := &http.Client{
		Jar: jar,
	}
	urlStr := "https://www.vinted." + country
	resp, err := client.Get(urlStr)
	if err != nil {
		log.Fatalf("Error getting URL: %v", err)
	}
	defer resp.Body.Close()

	u, _ := url.Parse(urlStr)
	cookies := jar.Cookies(u)

	for _, cookie := range cookies {
		if cookie.Name == "access_token_web" {
			fmt.Print("Access token generated")
			return cookie.Value
		}
	}

	log.Fatalf("access_token_web cookie not found")
	return ""
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	hakusana := os.Getenv("SEARCH")
	country := os.Getenv("COUNTRY")
	discordWebhook := os.Getenv("DISCORD_WEBHOOK")
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	telegramChatID := os.Getenv("TELEGRAM_CHAT_ID")

	if country == "" {
		log.Fatalf("COUNTRY not defined in .env file")
	}

	if discordWebhook == "" && (telegramBotToken == "" || telegramChatID == "") {
		fmt.Println("Notification method not defined. Please set either Discord Webhook or Telegram Bot Token and Chat ID, or both.")
		return
	}

	if hakusana != "" {
		fmt.Println("Search started with term `" + hakusana + "`...")
	} else {
		fmt.Println("Search started...")
	}

	seen := []string{}
	firstRun := true

	client := &fasthttp.Client{}

	accessToken := getAccessToken(country)

	for {
		foundTotal = 0
		req := fasthttp.AcquireRequest()
		req.SetRequestURI("https://www.vinted." + country + "/api/v2/catalog/items?page=1&per_page=96&search_text=" + hakusana + "&order=newest_first")
		req.Header.SetMethod("GET")
		req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3")
		req.Header.Set("Accept", "application/json")
		req.Header.Set("Cookie", "access_token_web="+accessToken)

		resp := fasthttp.AcquireResponse()
		err := client.Do(req, resp)
		if err != nil {
			log.Fatalf("Error: %s\n", err)
			continue
		}

		var response Response
		body := resp.Body()
		if err := json.Unmarshal(body, &response); err != nil {
			log.Fatalf("Error unmarshalling JSON: %v", err)
			continue
		}

		for _, product := range response.Items {
			if !slices.Contains(seen, product.URL) {
				if !firstRun {
					foundTotal++
					fmt.Println(product.Title + "\n" + product.Price.Amount + " " + product.Price.CurrencyCode + "\n" + product.URL + "\n" + "-----")

					if discordWebhook != "" {
						sendDiscordNotification(discordWebhook, product)
					}

					if telegramBotToken != "" && telegramChatID != "" {
						sendTelegramNotification(telegramBotToken, telegramChatID, product)
					}
				}
				seen = append(seen, product.URL)
			}
			if foundTotal >= 25 {
				break
			}
		}

		fasthttp.ReleaseRequest(req)
		fasthttp.ReleaseResponse(resp)

		if foundTotal > 0 {
			fmt.Println("New notifications found: ", strconv.Itoa(foundTotal)+"\n-----")
		}

		firstRun = false
		time.Sleep(10 * time.Second)
	}
}

func sendDiscordNotification(webhook string, product Product) {
	color := 2895667
	if product.ImageURL.DominantColor != "" {
		dominantColor, err := strconv.ParseInt(product.ImageURL.DominantColor[1:], 16, 64)
		if err == nil {
			color = int(dominantColor)
		}
	}
	ticks := "```"
	payload := fmt.Sprintf(`{"content": null,"embeds": [{"title": "%s","url": "%s","color": %d,"fields": [{"name": "Price","value": "%s%s %s%s","inline": true},{"name": "Brand","value": "%s%s%s","inline": true},{"name": "Size","value": "%s%s%s","inline": true},{"name": "Condition","value": "%s%s%s","inline": true}],"image": {"url": "%s"}}]}`,
		product.Title, product.URL, color, ticks, product.TotalItemPrice.Amount, product.TotalItemPrice.CurrencyCode, ticks, ticks, product.BrandTitle, ticks, ticks, product.SizeTitle, ticks, ticks, product.Status, ticks, product.ImageURL.URL)
	sendWebhook(webhook, payload)
}

func sendTelegramNotification(botToken, chatID string, product Product) {
	message := fmt.Sprintf("%s\nPrice: %s %s\nBrand: %s\nSize: %s\nCondition: %s\n%s\n[Photo](%s)",
		product.Title,
		product.TotalItemPrice.Amount,
		product.TotalItemPrice.CurrencyCode,
		product.BrandTitle,
		product.SizeTitle,
		product.Status,
		product.URL,
		product.ImageURL.URL)

	url := fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", botToken)
	payload := fmt.Sprintf(`{"chat_id": "%s", "text": "%s", "parse_mode": "Markdown"}`, chatID, message)

	sendWebhook(url, payload)
}

func sendWebhook(url, payload string) {
	req := fasthttp.AcquireRequest()
	req.SetRequestURI(url)
	req.Header.SetMethod("POST")
	req.Header.Set("Content-Type", "application/json")
	req.SetBody([]byte(payload))

	resp := fasthttp.AcquireResponse()
	client := &fasthttp.Client{}
	err := client.Do(req, resp)
	if err != nil {
		log.Printf("Error: %s\n", err)
	}

	fasthttp.ReleaseRequest(req)
	fasthttp.ReleaseResponse(resp)
}
