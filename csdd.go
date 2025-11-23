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

var curl string = `'https://e.csdd.lv/examp/' \
	-H 'Accept: text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.7' \
	-H 'Accept-Language: en-US,en;q=0.9,lv;q=0.8,ru;q=0.7' \
	-H 'Cache-Control: max-age=0' \
	-H 'Connection: keep-alive' \
	-H 'Content-Type: application/x-www-form-urlencoded' \
	-b 'PHPSESSID=cv700lms24m9sfevhsj9q6i9ka; eSign=8027a52c9d1126fe82a75fcbd22ce50c; SERVERID=s6; SimpleSAML=3dbece22c2c81f2f5d90bfd81022df83; SimpleSAMLAuthToken=_d44f5fc5664cfb56c050d27bfcf094b44294294e9c' \
	-H 'Origin: https://e.csdd.lv' \
	-H 'Referer: https://e.csdd.lv/examp/' \
	-H 'Sec-Fetch-Dest: document' \
	-H 'Sec-Fetch-Mode: navigate' \
	-H 'Sec-Fetch-Site: same-origin' \
	-H 'Sec-Fetch-User: ?1' \
	-H 'Upgrade-Insecure-Requests: 1' \
	-H 'User-Agent: Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/138.0.0.0 Safari/537.36 Edg/138.0.0.0' \
	-H 'sec-ch-ua: "Not)A;Brand";v="8", "Chromium";v="138", "Microsoft Edge";v="138"' \
	-H 'sec-ch-ua-mobile: ?0' \
	-H 'sec-ch-ua-platform: "Linux"' \
	--data-raw 'veids=5&did=2&kods=B&veids_txt=B&savs_tl_txt=&capcha=EjXmLaZ'`

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

func parseCurl(command string) []string {
	// 1. Normalize the command string
	// Replace newlines and backslashes used for line continuation
	command = strings.ReplaceAll(command, "\\\n", " ")
	command = strings.ReplaceAll(command, "\n", " ")

	// 2. Split the command into arguments
	var args []string
	// This regex splits by spaces, but keeps quoted sections together.
	// It handles single and double quotes.
	r := regexp.MustCompile(`[^\s"']+|"([^"]*)"|'([^']*)'`)
	matches := r.FindAllString(command, -1)

	for _, match := range matches {
		// 3. Remove the quotes from the matched arguments
		if len(match) > 1 && (match[0] == '"' && match[len(match)-1] == '"' || match[0] == '\'' && match[len(match)-1] == '\'') {
			args = append(args, match[1:len(match)-1])
		} else {
			args = append(args, match)
		}
	}

	// The first element is "curl", which should be omitted for exec.Command arguments
	if len(args) > 0 && args[0] == "curl" {
		return args[1:]
	}

	return args
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

			go sendOther(text, bot, []string{chat_id[i]})
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
		log.Fatal("error occurred: ", err)
		fmt.Println(string(out))
		log.Fatal(err)
	}

	return string(out)
}

func telegramBotUpdater(api string, adminPassword string, cfg Config) {
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
					priorityChatID = ""
					msg.Text = "priority removed"
				}
			}

		case "curl":
			curl = update.Message.Text
			msg.Text = "Curl updated"

		case "test":
			go send("TEST CHAT ID", cfg.Telegram.BotID, cfg.Telegram.ChatID)
			msg.Text = "Test message sent to all chat IDs"

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

	go telegramBotUpdater(cfg.Telegram.BotID, cfg.Admin.Password, cfg)

	untilDay, err := strconv.Atoi(cfg.Scraper.Date[:2])
	untilMonth, err1 := strconv.Atoi(cfg.Scraper.Date[3:5])
	untilYear, err2 := strconv.Atoi(cfg.Scraper.Date[6:10])
	if err != nil || err1 != nil || err2 != nil {
		log.Fatal(err)
	}

	defer send("Program die", cfg.Telegram.BotID, cfg.Telegram.ChatID)

	for {
		fmt.Println("Scraping...")

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
