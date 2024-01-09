package main

import (
	"bytes"
	"encoding/json"
	"fmt"
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

var curl string = `curl 'https://e.csdd.lv/examp/' \
-H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7' \
-H 'Accept-Language: en-US,en;q=0.9,lv;q=0.8,ru;q=0.7' \
-H 'Cache-Control: max-age=0' \
-H 'Connection: keep-alive' \
-H 'Content-Type: application/x-www-form-urlencoded' \
-H 'Cookie: _ga=GA1.1.432155588.1702485824; _ga_KSGMLEJL82=GS1.1.1702494283.3.0.1702494283.0.0.0; eSign=3a2f1edbfd3fc2513016674cb77c89a4; _hjSessionUser_3007240=eyJpZCI6IjNjODE2YTMzLTU5ZDktNWU5YS1iY2QyLTQ5ODU0YzE5OTYxMyIsImNyZWF0ZWQiOjE3MDI3NjM5OTEzMjYsImV4aXN0aW5nIjp0cnVlfQ==; _hjDonePolls=852170; _hjMinimizedPolls=852170; userSawThatSiteUsesCookies=1; PHPSESSID=2adi4h37ocbe4i9sm4a397tq25; _hjIncludedInSessionSample_3007240=0; _hjSession_3007240=eyJpZCI6ImNmNWZiMDAzLWUxOWItNDVkYS04YjcwLTI2ZDgxNTAwNTc2NSIsImMiOjE3MDM3NjI0NzEzMTMsInMiOjAsInIiOjAsInNiIjowfQ==; _hjAbsoluteSessionInProgress=0; SimpleSAML=3de3b99052a47d556c8ebef1e7511c9b; SERVERID=s6; SimpleSAMLAuthToken=_3815688ec46904625c9c496b9564b28719b7971c8e; _ga_Q09H2GL8G8=GS1.1.1703762470.32.1.1703763266.0.0.0' \
-H 'Origin: https://e.csdd.lv' \
-H 'Referer: https://e.csdd.lv/examp/' \
-H 'Sec-Fetch-Dest: document' \
-H 'Sec-Fetch-Mode: navigate' \
-H 'Sec-Fetch-Site: same-origin' \
-H 'Sec-Fetch-User: ?1' \
-H 'Upgrade-Insecure-Requests: 1' \
-H 'User-Agent: Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36 Edg/120.0.0.0' \
-H 'sec-ch-ua: "Not_A Brand";v="8", "Chromium";v="120", "Microsoft Edge";v="120"' \
-H 'sec-ch-ua-mobile: ?1' \
-H 'sec-ch-ua-platform: "Android"' \
--data-raw 'datums=-1&did=3&datums_txt=' \
--compressed`

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
	// client := &http.Client{}
	// var data = strings.NewReader(`datums=-1&did=3&datums_txt=`)
	// req, err := http.NewRequest("POST", "https://e.csdd.lv/examp/", data)
	// if err != nil {
	// 	log.Println(err)
	// 	return ""
	// }

	// req.Header.Set("Accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7")
	// req.Header.Set("Accept-Language", "en-US,en;q=0.9,lv;q=0.8,ru;q=0.7")
	// req.Header.Set("Cache-Control", "max-age=0")
	// req.Header.Set("Connection", "keep-alive")
	// req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	// req.Header.Set("Cookie", "_ga=GA1.1.432155588.1702485824; _ga_KSGMLEJL82=GS1.1.1702494283.3.0.1702494283.0.0.0; eSign=3a2f1edbfd3fc2513016674cb77c89a4; _hjSessionUser_3007240=eyJpZCI6IjNjODE2YTMzLTU5ZDktNWU5YS1iY2QyLTQ5ODU0YzE5OTYxMyIsImNyZWF0ZWQiOjE3MDI3NjM5OTEzMjYsImV4aXN0aW5nIjp0cnVlfQ==; _hjDonePolls=852170; _hjMinimizedPolls=852170; userSawThatSiteUsesCookies=1; PHPSESSID=2adi4h37ocbe4i9sm4a397tq25; _hjIncludedInSessionSample_3007240=0; _hjSession_3007240=eyJpZCI6ImNmNWZiMDAzLWUxOWItNDVkYS04YjcwLTI2ZDgxNTAwNTc2NSIsImMiOjE3MDM3NjI0NzEzMTMsInMiOjAsInIiOjAsInNiIjowfQ==; _hjAbsoluteSessionInProgress=0; SimpleSAML=3de3b99052a47d556c8ebef1e7511c9b; SERVERID=s6; SimpleSAMLAuthToken=_3815688ec46904625c9c496b9564b28719b7971c8e; _ga_Q09H2GL8G8=GS1.1.1703762470.32.1.1703763266.0.0.0")
	// req.Header.Set("Origin", "https://e.csdd.lv")
	// req.Header.Set("Referer", "https://e.csdd.lv/examp/")
	// req.Header.Set("Sec-Fetch-Dest", "document")
	// req.Header.Set("Sec-Fetch-Mode", "navigate")
	// req.Header.Set("Sec-Fetch-Site", "same-origin")
	// req.Header.Set("Sec-Fetch-User", "?1")
	// req.Header.Set("Upgrade-Insecure-Requests", "1")
	// req.Header.Set("User-Agent", "Mozilla/5.0 (Linux; Android 6.0; Nexus 5 Build/MRA58N) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/120.0.0.0 Mobile Safari/537.36 Edg/120.0.0.0")
	// req.Header.Set("sec-ch-ua", `"Not_A Brand";v="8", "Chromium";v="120", "Microsoft Edge";v="120"`)
	// req.Header.Set("sec-ch-ua-mobile", "?1")
	// req.Header.Set("sec-ch-ua-platform", `"Android"`)

	// resp, err := client.Do(req)
	// if err != nil {
	// 	log.Println(err)
	// 	return ""
	// }
	// defer resp.Body.Close()
	// bodyText, err := io.ReadAll(resp.Body)
	// if err != nil {
	// 	log.Println(err)
	// 	return ""
	// }
	// return string(bodyText)

	// delete all new lines in curl variable
	curl = strings.ReplaceAll(curl, "\n", "")
	// delete all ^ in curl variable
	curl = strings.ReplaceAll(curl, "^", "")
	curl = strings.ReplaceAll(curl, "\\", "")
	curl = strings.ReplaceAll(curl, "'", "\"")
	curl = strings.ReplaceAll(curl, "  ", " ")

	// delete curl
	curl = strings.ReplaceAll(curl, "curl ", "")

	fmt.Println(curl)

	cmd := exec.Command("curl", "-O", curl)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("error occured")
		fmt.Println(string(out))
		log.Fatal(err)
	}

	return string(out)
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
			curl = update.Message.Text
			msg.Text = "Curl updated"

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

	fmt.Printf("CSDD parse data app. v1.1\n")

	out := scrape()
	fmt.Println("captured data")
	// save to file
	f, err := os.Open("niggger.html")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	f.WriteString(out)

	os.Exit(0)

	// Parse configs
	f, err = os.Open("config/config.yml")
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
