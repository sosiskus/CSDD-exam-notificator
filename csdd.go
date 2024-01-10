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
	"runtime"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gopkg.in/yaml.v2"
)

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

var bot *tgbotapi.BotAPI
var globalStatus [][]string
var priorityChatID string = ""

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

func parseCurl(curll string) []string {
	var res []string

	// delete all new lines in curl variable
	// curll = strings.ReplaceAll(curll, "\n", "")
	// delete all ^ in curl variable
	curll = strings.ReplaceAll(curll, "^", "")
	curll = strings.ReplaceAll(curll, "\\", "")
	// curll = strings.ReplaceAll(curll, "'", "\"")
	// curll = strings.ReplaceAll(curll, "  ", " ")

	// fmt.Println(curll)

	// res = append(res, "curl")

	// extract url
	var re = regexp.MustCompile(`(?mU)curl\s*'(.+)'`)
	for _, match := range re.FindAllStringSubmatch(curll, -1) {
		res = append(res, match[1])
	}

	re = regexp.MustCompile(`(?mU)-H\s*'(.+)'`)

	for _, match := range re.FindAllStringSubmatch(curll, -1) {
		res = append(res, "-H")
		res = append(res, match[1])
	}
	res = append(res, "--data-raw")
	// extract data-raw
	re = regexp.MustCompile(`(?mU)--data-raw\s*'(.+)'`)
	for _, match := range re.FindAllStringSubmatch(curll, -1) {
		res = append(res, match[1])
	}

	// if os == windows then remove --compressed
	if runtime.GOOS == "linux" {
		res = append(res, "--compressed")
	}

	return res
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

	com := parseCurl(curl)
	for i := range com {
		fmt.Println(com[i])
	}

	cmd := exec.Command("curl", com...)
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

func main() {

	fmt.Printf("CSDD parse data app. v1.1\n")

	// Parse configs
	f, err := os.Open("config/config.yml")
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
				go send(str, cfg.Telegram.BotID, cfg.Telegram.ChatID)
				break
			}
		}
		time.Sleep(time.Duration(cfg.Scraper.WaitTimeMin) * time.Minute)
	}
}
