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

	"gopkg.in/yaml.v2"
)

type Config struct {
	Telegram struct {
		BotID  string `yaml:"bot_id"`
		ChatID string `yaml:"chat_id"`
	} `yaml:"telegram"`
	Scraper struct {
		WaitTimeMin int    `yaml:"wait_time_min"`
		Date        string `yaml:"date"`
	} `yaml:"scraper"`
}

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
	req.Header.Set("Accept-Language", "en-US,en;q=0.9")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Cookie", "PHPSESSID=gsl7rrh7suuqt35e3c3ne408ip; eSign=7cf6bd2046f334665e572e9504bfe40c; SERVERID=s6; _hjFirstSeen=1; _hjIncludedInSessionSample_3007240=0; _hjSession_3007240=eyJpZCI6IjUzYTAyZGMwLTk3MDAtNGJjZS05YmU2LTUyZjUwYTcwZjVkZSIsImMiOjE3MDI4NDU5MTU0MzgsInMiOjAsInIiOjAsInNiIjowfQ==; _hjAbsoluteSessionInProgress=0; SimpleSAML=837545cf14ad6ad527dee4b0b26fab1c; SimpleSAMLAuthToken=_06d3a64f432545c8e5517dc278a9f2ce73eadfb396; _ga=GA1.1.2074492875.1702845949; _hjSessionUser_3007240=eyJpZCI6IjU0ODhmZTJlLTliMDYtNWMxMC1iNzY1LWUzODZmZWE1ZTNhZCIsImNyZWF0ZWQiOjE3MDI4NDU5MTU0MzYsImV4aXN0aW5nIjp0cnVlfQ==; _ga_Q09H2GL8G8=GS1.1.1702845949.1.1.1702845992.0.0.0; _hjMinimizedPolls=852170")
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

	// Parse configs
	f, err := os.Open("config.yml")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		log.Fatal(err)
	}

	untilDay, err := strconv.Atoi(cfg.Scraper.Date[:2])
	untilMonth, err1 := strconv.Atoi(cfg.Scraper.Date[3:5])
	untilYear, err2 := strconv.Atoi(cfg.Scraper.Date[6:10])
	if err != nil || err1 != nil || err2 != nil {
		log.Fatal(err)
	}

	for {
		fmt.Println(time.Now())

		plainHtml := scrape()
		// n, _ := ioutil.ReadFile("niggger.html")
		// plainHtml := string(n)

		var re = regexp.MustCompile(`(?mU)<option\s*value="[0-9]+"\s*>(.+)</option>`)
		res := re.FindAllStringSubmatch(plainHtml, -1)

		if len(res) <= 0 {
			fmt.Printf("session die\n")
			send("Session die", cfg.Telegram.BotID, cfg.Telegram.ChatID)
			time.Sleep(time.Duration(cfg.Scraper.WaitTimeMin) * time.Minute)
			continue
		}

		for i := range res {
			str := res[i][1]
			last_chs := strings.TrimSpace(str[len(str)-2:])
			date := strings.TrimSpace(str[:10])

			dateDay, err1 := strconv.Atoi(date[:2])
			dateMonth, err := strconv.Atoi(date[3:5])
			dateYear, err2 := strconv.Atoi(date[6:10])

			if err != nil || err1 != nil || err2 != nil {
				continue
			}

			fmt.Printf("%s [%s,%s]\n", str, []byte(date), last_chs)

			end := time.Date(untilYear, time.Month(untilMonth), untilDay, 0, 0, 0, 0, time.UTC)

			dateToCheck := time.Date(dateYear, time.Month(dateMonth), dateDay, 0, 0, 0, 0, time.UTC)

			if dateToCheck.Before(end) && last_chs != "0" {
				fmt.Printf("found\n")
				send(str, cfg.Telegram.BotID, cfg.Telegram.ChatID)
				break
			}
		}
		time.Sleep(time.Duration(cfg.Scraper.WaitTimeMin) * time.Minute)
	}
}
