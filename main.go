package main

import (
	"TokaiCrossBot/db"
	"TokaiCrossBot/twitter"
	"fmt"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

type Value struct {
	Id  string `db:"id"`
	Val string `db:"val"`
}

type Member struct {
	Id       string `db:"member"`
	Val      string `db:"content"`
	IsFinish bool   `db:"isFinish"`
}

func main() {

	go updateCross()
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		membersName := []string{
			"てつや",
			"しばゆー",
			"虫眼鏡",
			"りょう",
			"としみつ",
			"ゆめまる",
		}
		html := "<html><body>"
		for i := 0; i < 6; i++ {
			mem := []Member{}
			db.Select(&mem, "select * from cross where member=?", i)
			html += "<h1>" + membersName[i] + "</h1>"
			html += "<ul>"
			for _, cross := range mem {
				html += "<li>"
				html += cross.Val
				if cross.IsFinish {
					html += "<font size=\"8em\" color=\"red\">済</font>"
				}
				html += "</li>"
			}
			html += "</ul>"
		}

		c.Writer.WriteString(html)
	})
	router.Run(":3000")
}

func updateCross() {
	members := [][]string{{"てつや", "てっちゃん"},
		{"しばゆー", "ハイオクマンタン"},
		{"虫眼鏡", "虫さん"},
		{"りょう", "イタリア人"},
		{"としみつ", "鈴木"},
		{"ゆめまる", "ごわす"},
	}

	token := twitter.GetBearerToken()
	userId := "1528038485390524422"
	//userId := "1050383043330854913"
	for {
		time.Sleep(time.Second * 5)
		params := map[string]string{
			"max_results": "100",
		}
		if id, err := getSince(); err == nil {
			params["since_id"] = id
		}
		tweets := twitter.GetMention(token, userId, params)
		if len(tweets.Data) == 0 {
			continue
		}
		setSince(tweets.Data[0].Id)
		for _, data := range tweets.Data {
			for _, content := range strings.Split(data.Text, "\n") {
				text := strings.Split(content, " ")
				fmt.Println(text)
				if len(text) < 3 {
					continue
				}
				if text[0] == "add" {
					for i, v := range members {
						if arrayContains(v, text[1]) {
							db.Push("insert into cross values(?,?,?)", i, text[2], false)
							break
						}
					}
				} else if text[0] == "delete" {
					for _, v := range members {
						if arrayContains(v, text[1]) {
							db.Push("update cross set isFinish=? where content = ?", 1, text[2])
							break
						}
					}
				}
			}
		}
	}
}

func getSince() (string, error) {
	value := []Value{}
	db.Select(&value, "select * from value where id=?", "since_id")
	if len(value) == 0 {
		return "0", fmt.Errorf("Error: %s\n", "since id is not found")
	}
	return value[0].Val, nil
}

func setSince(val string) {
	db.Push("update value set val=? where id=?", val, "since_id")
}
func arrayContains(arr []string, str string) bool {
	for _, v := range arr {
		if v == str {
			return true
		}
	}
	return false
}
