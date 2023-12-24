package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v2"
)

var bot *tgbotapi.BotAPI
var globalStatus [][]string
var priorityChatID string = ""

var cookie string = ""

type Config struct {
	Telegram struct {
		BotID  string   `yaml:"bot_id"`
		ChatID []string `yaml:"chat_id"`
	} `yaml:"telegram"`
	Scraper struct {
		WaitTimeMin int    `yaml:"wait_time_min"`
		Date        string `yaml:"date"`
	} `yaml:"scraper"`
	Admin struct {
		Password string `yaml:"password"`
	} `yaml:"admin"`
}

func remove(slice []string, s int) []string {
	return append(slice[:s], slice[s+1:]...)
}

func send(text string, bot string, chat_id []string) {

	request_url := "https://api.telegram.org/bot" + bot + "/sendMessage"

	client := &http.Client{}

	for i := range chat_id {

		if priorityChatID != "" && chat_id[i] == priorityChatID {
			values := map[string]string{"text": text, "chat_id": chat_id[i]}
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

			chat_id = remove(chat_id, i)

			break
		} else if priorityChatID == "" {
			values := map[string]string{"text": text, "chat_id": chat_id[i]}
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

	}

	if priorityChatID != "" {
		go sendOther(text, bot, chat_id)
	}

}

func sendOther(text string, bot string, chat_id []string) {
	time.Sleep(3 * time.Minute)

	fmt.Println("OTHERS")

	request_url := "https://api.telegram.org/bot" + bot + "/sendMessage"

	client := &http.Client{}

	for i := range chat_id {

		fmt.Println("sennding message to " + chat_id[i])

		values := map[string]string{"text": text, "chat_id": chat_id[i]}
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
}

func scrape() string {
	client := &http.Client{}
	var data = strings.NewReader(`datums=-1&did=3&datums_txt=`)
	req, err := http.NewRequest("POST", "https://e.csdd.lv/examp/", data)
	if err != nil {
		log.Println(err)
		return ""
	}

	req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	req.Header.Set("Accept-Language", "en-US,en;q=0.9,lv;q=0.8,ru;q=0.7")
	req.Header.Set("Cache-Control", "max-age=0")
	req.Header.Set("Connection", "keep-alive")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	if cookie != "" {
		req.Header.Set("Cookie", cookie)
		cookie = ""
	} else {
		req.Header.Set("Cookie", "_ga=GA1.1.432155588.1702485824; _ga_KSGMLEJL82=GS1.1.1702494283.3.0.1702494283.0.0.0; eSign=3a2f1edbfd3fc2513016674cb77c89a4; _hjSessionUser_3007240=eyJpZCI6IjNjODE2YTMzLTU5ZDktNWU5YS1iY2QyLTQ5ODU0YzE5OTYxMyIsImNyZWF0ZWQiOjE3MDI3NjM5OTEzMjYsImV4aXN0aW5nIjp0cnVlfQ==; _hjDonePolls=852170; _hjMinimizedPolls=852170; userSawThatSiteUsesCookies=1; PHPSESSID=1bi0ngsmd5r21o76185gvstqnk; _hjIncludedInSessionSample_3007240=0; _hjSession_3007240=eyJpZCI6ImUzMjZkZjBjLThhM2EtNDIwOS04NTNlLWI5YjI0NjQzMDc2NiIsImMiOjE3MDMwODUxMDIxMTksInMiOjAsInIiOjAsInNiIjowfQ==; _hjAbsoluteSessionInProgress=0; SimpleSAML=13f3d17b2928f90408fe2c6abb6eb890; SimpleSAMLAuthToken=_e312727e3d4ae7a836d06ed14d000c3f7770574bf7; SERVERID=s8; _ga_Q09H2GL8G8=GS1.1.1703085034.16.1.1703085189.0.0.0")
	}

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
		log.Println(err)
		return ""
	}
	defer resp.Body.Close()
	bodyText, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return ""
	}
	return string(bodyText)
}

func telegramBotUpdater(api string, adminPassword string) {
	bot, err := tgbotapi.NewBotAPI(api)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message updates
			continue
		}

		if !update.Message.IsCommand() { // ignore any non-command Messages
			continue
		}

		// Create a new MessageConfig. We don't have text yet,
		// so we leave it empty.
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "")

		// Extract the command from the Message.
		switch update.Message.Command() {
		case "status":
			var total string
			for i := range globalStatus {
				total += globalStatus[i][1] + "\n"
			}
			if len(total) == 0 {
				msg.Text = "No entries yet"
			} else {
				msg.Text = total
			}

		case "priority":
			msg.Text = "Incorrect password"
			res := strings.Split(update.Message.Text, " ")
			if len(res) > 1 {
				if res[1] == adminPassword {
					priorityChatID = strconv.Itoa(int(update.Message.Chat.ID))
					msg.Text = "priority set to" + priorityChatID
				}
			}

		case "rpriority":
			msg.Text = "Incorrect password"
			res := strings.Split(update.Message.Text, " ")
			if len(res) > 1 {
				if res[1] == adminPassword {
					msg.Text = "priority removed"
				}
			}

		case "curl":
			findIn := update.Message.Text
			n := strings.Index(findIn, "Cookie: ")

			for i := n; ; i++ {
				if string(findIn[i]) == "'" {
					break
				}
				cookie += string(findIn[i])
			}
			cookie = strings.ReplaceAll(cookie, "Cookie: ", "")
			msg.Text = cookie

		default:
			msg.Text = "I don't know that command"
		}

		if _, err := bot.Send(msg); err != nil {
			log.Panic(err)
		}
	}
}

func restart() {
	cmd := exec.Command(os.Args[0])
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err == nil {
		os.Exit(0)
	}
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

	go telegramBotUpdater(cfg.Telegram.BotID, cfg.Admin.Password)

	untilDay, err := strconv.Atoi(cfg.Scraper.Date[:2])
	untilMonth, err1 := strconv.Atoi(cfg.Scraper.Date[3:5])
	untilYear, err2 := strconv.Atoi(cfg.Scraper.Date[6:10])
	if err != nil || err1 != nil || err2 != nil {
		log.Fatal(err)
	}

	defer send("Program die", cfg.Telegram.BotID, cfg.Telegram.ChatID)

	var found bool = false
	for {
		fmt.Println(time.Now())

		plainHtml := scrape()
		// n, _ := ioutil.ReadFile("niggger.html")
		// plainHtml := string(n)

		var re = regexp.MustCompile(`(?mU)<option\s*value="[0-9]+"\s*>(.+)</option>`)
		res := re.FindAllStringSubmatch(plainHtml, -1)

		if len(res) <= 0 {
			fmt.Printf("session die\n")

			fmt.Println(plainHtml)
			if found {
				restart()
			}

			go send("Session die", cfg.Telegram.BotID, cfg.Telegram.ChatID)
			time.Sleep(time.Duration(cfg.Scraper.WaitTimeMin) * time.Minute)
			continue
		}

		res[len(res)-1][1] = time.Now().String()
		globalStatus = res

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
				found = true
				go send(str, cfg.Telegram.BotID, cfg.Telegram.ChatID)
				break
			} else {
				found = false
			}
		}
		time.Sleep(time.Duration(cfg.Scraper.WaitTimeMin) * time.Minute)
	}
}
