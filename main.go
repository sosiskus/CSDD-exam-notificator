package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func send(text string, bot string, chat_id string) {

	request_url := "https://api.telegram.org/" + bot + "/sendMessage"

	client := &http.Client{}

	values := map[string]string{"text": text, "chat_id": chat_id}
	json_paramaters, _ := json.Marshal(values)

	req, _ := http.NewRequest("POST", request_url, bytes.NewBuffer(json_paramaters))
	req.Header.Set("Content-Type", "application/json")

	res, err := client.Do(req)

	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(res.Status)
		defer res.Body.Close()
	}
}

func scrape() string {
	client := &http.Client{}
	var data = strings.NewReader(`datums=-1&did=3&datums_txt=`)
	req, err := http.NewRequest("POST", "https://e.csdd.lv/examp/", data)
	if err != nil {
		log.Fatal(err)
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,lv;q=0.8,ru;q=0.7")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "_ga=GA1.1.432155588.1702485824; _ga_KSGMLEJL82=GS1.1.1702494283.3.0.1702494283.0.0.0; PHPSESSID=lheraup71fdttk94nptii9aeu3; eSign=3a2f1edbfd3fc2513016674cb77c89a4; _hjSessionUser_3007240=eyJpZCI6IjNjODE2YTMzLTU5ZDktNWU5YS1iY2QyLTQ5ODU0YzE5OTYxMyIsImNyZWF0ZWQiOjE3MDI3NjM5OTEzMjYsImV4aXN0aW5nIjp0cnVlfQ==; _hjDonePolls=852170; _hjMinimizedPolls=852170; userSawThatSiteUsesCookies=1; _hjIncludedInSessionSample_3007240=0; _hjSession_3007240=eyJpZCI6ImYwNWZlN2U4LTY4OWUtNDgyOS05ZWQwLTg5NmYwN2YwZTY4YSIsImMiOjE3MDI4MjQ3MDAwMDMsInMiOjAsInIiOjAsInNiIjowfQ==; _hjAbsoluteSessionInProgress=0; SimpleSAML=f4dd95a33a50dfbd0c76f0efa5b59a27; SERVERID=s4; SimpleSAMLAuthToken=_fdcd18f9f7c7b287206e721f10e490f883e5b6f6a8; _ga_Q09H2GL8G8=GS1.1.1702824681.5.1.1702825534.0.0.0")
	req.Header.Set("Origin", "https://e.csdd.lv")
	req.Header.Set("Referer", "https://e.csdd.lv/examp/")
	req.Header.Set("Sec-Fetch-Dest", "document")
	req.Header.Set("Sec-Fetch-Mode", "navigate")
	req.Header.Set("Sec-Fetch-Site", "same-origin")
	req.Header.Set("Sec-Fetch-User", "?1")
	req.Header.Set("Upgrade-Insecure-Requests", "1")
	req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36 Edg/120.0.0.0")
	req.Header.Set("sec-ch-ua", `"Not_A Brand";v="8", "Chromium";v="120", "Microsoft Edge";v="120"`)
	req.Header.Set("sec-ch-ua-mobile", "?1")
	req.Header.Set("sec-ch-ua-platform", `"Android"`)

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	return string(bodyText)
}

func main() {

	fmt.Printf("CSDD parse data app. v0.3\n")

	for {

		fmt.Println(time.Now())
		plainHtml := scrape()
		// plainHtml, _ := ioutil.ReadFile("niggger.html")

		var re = regexp.MustCompile(`(?mU)<option\s*value="[0-9]+"\s*>(.+)</option>`)

		str := string(plainHtml)
		res := re.FindAllStringSubmatch(str, -1)

		for i := range res {
			str := res[i][1]
			last_chs := strings.TrimSpace(str[len(str)-2:])
			date := strings.TrimSpace(str[:5])

			dateDay, err := strconv.Atoi(date[:2])
			dateMonth, err := strconv.Atoi(date[3:5])

			if err != nil {
				continue
			}

			// 05.09 less than
			if dateMonth == 12 && dateDay <= 31 && last_chs != "0" {
				fmt.Printf("found\n")
				send(str, "bot6621115120:AAGbttFrh_3ZjTl7YkUFpHaFZVj0kZnIXMc", "773077583")
				break
			}
			if dateMonth == 1 && dateDay < 8 && last_chs != "0" {
				fmt.Printf("found\n")
				send(str, "bot6621115120:AAGbttFrh_3ZjTl7YkUFpHaFZVj0kZnIXMc", "773077583")
				break
			}

			fmt.Printf("%s [%s,%s]\n", str, []byte(date), last_chs)
		}

		if len(res) <= 0 {
			fmt.Printf("session die\n")
			send(str, "bot6621115120:AAGbttFrh_3ZjTl7YkUFpHaFZVj0kZnIXMc", "773077583")
			os.Exit(1)
		}

		time.Sleep(5 * time.Minute)
	}
}
