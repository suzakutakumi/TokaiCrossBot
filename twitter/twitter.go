package twitter

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type Attachments struct {
	MediaKeys []string `json:"media_keys"`
}

type Tweet struct {
	Id     string       `json:"id"`
	Text   string       `json:"text"`
	Attach *Attachments `json:"attachments"`
}

type MetaData struct {
	Count    int     `json:"result_count"`
	NewestId string  `json:"newest_id"`
	OldestId string  `json:"oldest_id"`
	Next     *string `json:"next_token"`
}

type TimeLine struct {
	Data []Tweet  `json:"data"`
	Meta MetaData `json:"meta"`
}

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Username string `json:"username"`
}

func GetBearerToken() string {
	err := godotenv.Load(".env")
	if err != nil {
		panic(err)
	}
	env := os.Getenv("BearerToken")
	return env
}

func GetMention(token string, userId string, querys map[string]string) TimeLine {
	url := "https://api.twitter.com/2/users/" + userId + "/mentions"

	req, err := http.NewRequest("GET", url, nil)

	if err != nil {
		panic(err)
	}

	req.Header.Set("Authorization", "Bearer "+token)

	params := req.URL.Query()
	for key, val := range querys {
		params.Add(key, val)
	}
	req.URL.RawQuery = params.Encode()

	client := new(http.Client)
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	byteArray, _ := ioutil.ReadAll(resp.Body)
	//fmt.Println(string(byteArray))
	var timeLine TimeLine
	_ = json.Unmarshal(byteArray, &timeLine)
	return timeLine
}
