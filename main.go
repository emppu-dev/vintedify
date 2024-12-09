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

type Items struct {
	Items                []Item                    `json:"items"`
	DominantBrand        interface{}               `json:"dominant_brand"`
	SearchTrackingParams ItemsSearchTrackingParams `json:"search_tracking_params"`
	Pagination           Pagination                `json:"pagination"`
	Code                 int64                     `json:"code"`
}

type Item struct {
	ID                   int64                    `json:"id"`
	Title                string                   `json:"title"`
	Price                Price                    `json:"price"`
	IsVisible            bool                     `json:"is_visible"`
	Discount             interface{}              `json:"discount"`
	BrandTitle           string                   `json:"brand_title"`
	Path                 string                   `json:"path"`
	User                 User                     `json:"user"`
	Conversion           *Conversion              `json:"conversion"`
	URL                  string                   `json:"url"`
	Promoted             bool                     `json:"promoted"`
	Photo                ItemPhoto                `json:"photo"`
	FavouriteCount       int64                    `json:"favourite_count"`
	IsFavourite          bool                     `json:"is_favourite"`
	Badge                interface{}              `json:"badge"`
	ServiceFee           Price                    `json:"service_fee"`
	TotalItemPrice       Price                    `json:"total_item_price"`
	ViewCount            int64                    `json:"view_count"`
	SizeTitle            string                   `json:"size_title"`
	ContentSource        ContentSource            `json:"content_source"`
	Status               Status                   `json:"status"`
	IconBadges           []interface{}            `json:"icon_badges"`
	SearchTrackingParams ItemSearchTrackingParams `json:"search_tracking_params"`
}

type Conversion struct {
	SellerPrice    string         `json:"seller_price"`
	SellerCurrency SellerCurrency `json:"seller_currency"`
	BuyerCurrency  CurrencyCode   `json:"buyer_currency"`
	FxRoundedRate  string         `json:"fx_rounded_rate"`
	FxBaseAmount   string         `json:"fx_base_amount"`
	FxMarkupRate   string         `json:"fx_markup_rate"`
}

type ItemPhoto struct {
	ID                  int64          `json:"id"`
	ImageNo             int64          `json:"image_no"`
	Width               int64          `json:"width"`
	Height              int64          `json:"height"`
	DominantColor       string         `json:"dominant_color"`
	DominantColorOpaque string         `json:"dominant_color_opaque"`
	URL                 string         `json:"url"`
	IsMain              bool           `json:"is_main"`
	Thumbnails          []Thumbnail    `json:"thumbnails"`
	HighResolution      HighResolution `json:"high_resolution"`
	IsSuspicious        bool           `json:"is_suspicious"`
	FullSizeURL         string         `json:"full_size_url"`
	IsHidden            bool           `json:"is_hidden"`
	Extra               Extra          `json:"extra"`
}

type Extra struct {
}

type HighResolution struct {
	ID          string `json:"id"`
	Timestamp   int64  `json:"timestamp"`
	Orientation *int64 `json:"orientation"`
}

type Thumbnail struct {
	Type         Type   `json:"type"`
	URL          string `json:"url"`
	Width        int64  `json:"width"`
	Height       int64  `json:"height"`
	OriginalSize *bool  `json:"original_size"`
}

type Price struct {
	Amount       string       `json:"amount"`
	CurrencyCode CurrencyCode `json:"currency_code"`
}

type ItemSearchTrackingParams struct {
	Score          int64       `json:"score"`
	MatchedQueries interface{} `json:"matched_queries"`
}

type User struct {
	ID         int64      `json:"id"`
	Login      string     `json:"login"`
	ProfileURL string     `json:"profile_url"`
	Photo      *UserPhoto `json:"photo"`
	Business   bool       `json:"business"`
}

type UserPhoto struct {
	ID                  int64               `json:"id"`
	Width               int64               `json:"width"`
	Height              int64               `json:"height"`
	TempUUID            interface{}         `json:"temp_uuid"`
	URL                 string              `json:"url"`
	DominantColor       DominantColor       `json:"dominant_color"`
	DominantColorOpaque DominantColorOpaque `json:"dominant_color_opaque"`
	Thumbnails          []Thumbnail         `json:"thumbnails"`
	IsSuspicious        bool                `json:"is_suspicious"`
	Orientation         interface{}         `json:"orientation"`
	HighResolution      HighResolution      `json:"high_resolution"`
	FullSizeURL         string              `json:"full_size_url"`
	IsHidden            bool                `json:"is_hidden"`
	Extra               Extra               `json:"extra"`
}

type Pagination struct {
	CurrentPage  int64 `json:"current_page"`
	TotalPages   int64 `json:"total_pages"`
	TotalEntries int64 `json:"total_entries"`
	PerPage      int64 `json:"per_page"`
	Time         int64 `json:"time"`
}

type ItemsSearchTrackingParams struct {
	SearchCorrelationID   string `json:"search_correlation_id"`
	SearchSessionID       string `json:"search_session_id"`
	GlobalSearchSessionID string `json:"global_search_session_id"`
}

type ContentSource string

const (
	Search ContentSource = "search"
)

type CurrencyCode string

const (
	Eur CurrencyCode = "EUR"
)

type SellerCurrency string

const (
	Dkk SellerCurrency = "DKK"
	Pln SellerCurrency = "PLN"
	Sek SellerCurrency = "SEK"
)

type Type string

const (
	Thumb100     Type = "thumb100"
	Thumb150     Type = "thumb150"
	Thumb150X210 Type = "thumb150x210"
	Thumb20      Type = "thumb20"
	Thumb310     Type = "thumb310"
	Thumb310X430 Type = "thumb310x430"
	Thumb364X428 Type = "thumb364x428"
	Thumb428X624 Type = "thumb428x624"
	Thumb50      Type = "thumb50"
	Thumb624X428 Type = "thumb624x428"
	Thumb70X100  Type = "thumb70x100"
)

type Status string

const (
	Good           Status = "Good"
	NewWithTags    Status = "New with tags"
	NewWithoutTags Status = "New without tags"
	Satisfactory   Status = "Satisfactory"
	VeryGood       Status = "Very good"
)

type DominantColor string

const (
	The31Abc2 DominantColor = "#31abc2"
)

type DominantColorOpaque string

const (
	C1E6Ed DominantColorOpaque = "#C1E6ED"
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

var foundTotal int

func getAccessToken() string {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalf("Error creating cookie jar: %v", err)
	}

	client := &http.Client{
		Jar: jar,
	}

	urlStr := "https://www.vinted.fi"
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

	hakusana := os.Getenv("HAKUSANA")
	discordWebhook := os.Getenv("DISCORD_WEBHOOK")
	telegramBotToken := os.Getenv("TELEGRAM_BOT_TOKEN")
	telegramChatID := os.Getenv("TELEGRAM_CHAT_ID")

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

	accessToken := getAccessToken()
	//fmt.Println("Access Token:", accessToken)

	for {
		foundTotal = 0
		req := fasthttp.AcquireRequest()
		req.SetRequestURI("https://www.vinted.fi/api/v2/catalog/items?page=1&per_page=96&search_text=" + hakusana + "&order=newest_first")
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
		//fmt.Println("Response: ", string(resp.Body()))

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
	payload := fmt.Sprintf(`{"content": null,"embeds": [{"title": "%s","url": "%s","color": %d,"fields": [{"name": "Hinta","value": "%s%s EUR%s","inline": true},{"name": "Merkki","value": "%s%s%s","inline": true},{"name": "Koko","value": "%s%s%s","inline": true},{"name": "Kunto","value": "%s%s%s","inline": true}],"image": {"url": "%s"}}]}`,
		product.Title, product.URL, color, ticks, product.TotalItemPrice.Amount, ticks, ticks, product.BrandTitle, ticks, ticks, product.SizeTitle, ticks, ticks, product.Status, ticks, product.ImageURL.URL)
	sendWebhook(webhook, payload)
}

func sendTelegramNotification(botToken, chatID string, product Product) {
	message := fmt.Sprintf("%s\nPrice: %s %s\n%s\n[Photo](%s)",
		product.Title,
		product.Price.Amount,
		product.Price.CurrencyCode,
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
