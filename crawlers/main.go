package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	"net/url"

	"github.com/atotto/clipboard"
)

type ApiResponseData struct {
	Media      []Media `json:"media"`
	NextCursor string  `json:"next_cursor"`
}

type Media struct {
	Type  string `json:"type"`
	Image Image  `json:"image"`
}

type Image struct {
	ID             string `json:"_id"`
	AdaptiveBase   string `json:"adaptive_base"`
	SiteProfileURL string `json:"site_profile_image_url"`
	ResponsiveURL  string `json:"responsive_url"`
	ShareLink      string `json:"share_link"`
	// add other fields as needed
}

func main() {
	bearerTokenPtr := flag.String("", "", "Bearer token")
	flag.Parse()

	bearerToken := *bearerTokenPtr

	baseUrl := "https://vsco.co/api/3.0/medias/profile"
	paramsMap := map[string]string{
		"site_id": "51936585",
		"limit":   "50",
		"cursor":  "",
	}
	urlValues := urlValuesFromMap(paramsMap)
	apiUrl := baseUrl + "?" + urlValues.Encode()
	fmt.Println("apiUrl", apiUrl)

	req, err := http.NewRequest(http.MethodGet, apiUrl, nil)
	if err != nil {
		fmt.Println("error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/122.0.0.0 Safari/537.36")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", bearerToken))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error making request:", err)
		return
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("API error:", resp.StatusCode)
		return
	}

	// make resp go readable
	var respData ApiResponseData
	err = json.NewDecoder(resp.Body).Decode(&respData)
	if err != nil {
		fmt.Println("Error decoding JSON response:", err)
		return
	}

	var urls []string
	for _, imgData := range respData.Media {
		urls = append(urls, imgData.Image.ResponsiveURL)
	}
	// final json data
	formattedJSON, err := json.MarshalIndent(urls, "", "    ")
	if err != nil {
		fmt.Println("Error formatting JSON:", err)
		return
	}
	fmt.Println("Response body:", string(formattedJSON))

	clipboard.WriteAll(string(formattedJSON))
}

func urlValuesFromMap(data map[string]string) url.Values {
	values := url.Values{}
	for key, value := range data {
		values.Add(key, value)
	}
	return values
}
