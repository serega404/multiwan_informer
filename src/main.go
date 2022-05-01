package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"
)

func messFormat(iface Interface, successful bool) string {
	if successful {
		return fmt.Sprintf("✅ Inteface *%s* now is *up*!", iface.DisplayName)
	} else {
		return fmt.Sprintf("❌ Inteface *%s* now is *down!*", iface.DisplayName)
	}
}

func main() {

	config := loadConfig()

	validateConfig(config)

	pingable := []PingData{}

	data, dataOk := loadData()

	if !dataOk {
		sendMessage("📁 *Data file is not readable*", config.TelegramConf)
	}

	for iface_indx, iface := range config.Interfaces {
		pingResult := ping("-w " + fmt.Sprint(config.WaitTimeSec) + " -I " + iface.IpOrInterfaceName + " " + config.PingAddr)

		pingable = append(pingable, PingData{iface.IpOrInterfaceName, pingResult, time.Now()})

		if dataOk {
			getIndex := func() int {
				for i, element := range data {
					if element.Name == iface.IpOrInterfaceName {
						return (i)
					}
				}
				return (-1)
			}

			idx := getIndex()

			if idx == -1 { // send message if interface not exist in data
				sendMessage(messFormat(iface, pingResult), config.TelegramConf)
			} else if data[idx].Successful != pingResult { // send messages if the result is different from what is in the data
				if pingResult {
					deltaTime := pingable[iface_indx].LastEdit.Sub(data[idx].LastEdit)
					sendMessage(
						messFormat(iface, pingResult)+
							fmt.Sprintf("\n⏳ Down time: *%d* days *%d* hours *%d* min *%d* sec", int(deltaTime.Hours())/24, int(deltaTime.Hours())%24, int(deltaTime.Minutes())%60, int(deltaTime.Seconds())%60%60),
						config.TelegramConf)
				} else {
					sendMessage(
						messFormat(iface, pingResult),
						config.TelegramConf)
				}
			}
		} else { // send message if data not loaded
			sendMessage(messFormat(iface, pingResult), config.TelegramConf)
		}

	}

	filse, _ := json.Marshal(pingable)

	_ = ioutil.WriteFile("data.json", filse, 0644)

}

func ping(params string) bool {
	Command := fmt.Sprintf("ping -c 1 %s  > /dev/null && echo true || echo false", params)
	output, err := exec.Command("/bin/sh", "-c", Command).Output()
	if err != nil {
		fmt.Println(err)
		return false
	}
	real_ip := strings.TrimSpace(string(output))
	if real_ip == "false" {
		return false
	} else {
		return true
	}
}

func sendMessage(mess string, tg_conf Telegram) bool {
	resp, err := http.Get("https://api.telegram.org/bot" + tg_conf.BotToken + "/sendMessage?parse_mode=Markdown&chat_id=" + tg_conf.ChatID + "&disable_notification=" + tg_conf.SendSilent + "&text=" + url.QueryEscape(mess))
	if err != nil {
		log.Fatalln(err)
		return false
	}

	if resp.StatusCode == 200 {
		return true
	}

	return false
}

func loadData() ([]PingData, bool) {
	file, _ := os.Open("data.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	data := []PingData{}
	err := decoder.Decode(&data)
	if err != nil {
		log.Println("Data file error:", err)
		return []PingData{}, false
	}

	return data, true
}

func loadConfig() Configuration {
	file, _ := os.Open("conf.json")
	defer file.Close()
	decoder := json.NewDecoder(file)
	configuration := Configuration{}
	err := decoder.Decode(&configuration)
	if err != nil {
		log.Fatalln("error:", err)
	}

	return configuration
}

func validateConfig(conf Configuration) {
	if !(len(conf.Interfaces) > 0) {
		log.Fatalln("You haven't added interfaces")
	}

	if conf.PingAddr == "" || !(conf.WaitTimeSec > 0) {
		log.Fatalln("Bad config")
	}

	// Telegram
	resp, err := http.Get("https://api.telegram.org/bot" + conf.TelegramConf.BotToken + "/getMe")
	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	var result map[string]interface{}
	json.Unmarshal([]byte(body), &result)

	if resp.StatusCode != 200 || fmt.Sprintf("%v", result["ok"]) != "true" {
		log.Fatalln("Failed to connect to bot")
	}
}

// config structs
type Configuration struct {
	Interfaces  []Interface
	WaitTimeSec int
	PingAddr    string

	TelegramConf Telegram
}

type Interface struct {
	DisplayName       string
	IpOrInterfaceName string
}

type Telegram struct {
	BotToken   string
	ChatID     string
	SendSilent string
}

// data struct

type PingData struct {
	Name       string
	Successful bool
	LastEdit   time.Time
}
