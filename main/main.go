package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	graphqlEndpoint = "https://api.github.com/graphql"
)

type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

type GraphQLResponse struct {
	Data   map[string]interface{} `json:"data"`
	Errors []interface{}          `json:"errors,omitempty"`
}

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("読み込み失敗")
		return
	}
	token := os.Getenv("GITHUB_TOKEN")
	userName := os.Getenv("GITHUB_ACCOUNT")

	// 昨日の日付を取得
	yesterday := time.Now().AddDate(0, 0, -1).Format("2006-01-02")

	query := `
		query ($userName:String!, $startDate:DateTime!, $endDate:DateTime!) {
			user(login: $userName) {
				contributionsCollection(from: $startDate, to: $endDate) {
					contributionCalendar {
						totalContributions
					}
				}
			}
		}
	`

	variables := map[string]interface{}{
		"userName":  userName,
		"startDate": yesterday + "T00:00:00Z", // 開始日時
		"endDate":   yesterday + "T23:59:59Z", // 終了日時
	}

	request := GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	jsonRequest, err := json.Marshal(request)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		return
	}

	req, err := http.NewRequest("POST", graphqlEndpoint, bytes.NewBuffer(jsonRequest))
	if err != nil {
		fmt.Println("Error creating HTTP request:", err)
		return
	}

	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending HTTP request:", err)
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	var graphQLResp GraphQLResponse
	err = json.Unmarshal(body, &graphQLResp)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return
	}

	// totalContributionsを表示
	totalContributions := graphQLResp.Data["user"].(map[string]interface{})["contributionsCollection"].(map[string]interface{})["contributionCalendar"].(map[string]interface{})["totalContributions"]
	fmt.Println("Total Contributions:", totalContributions)
}
