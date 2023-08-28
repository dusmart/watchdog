package main

import (
	"encoding/json"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const defaultInterval = time.Minute
var config Config
var alarms map[string]*time.Timer

func init() {
	configFile := WithoutError(os.Open("config.json"))
	defer configFile.Close()
	json.NewDecoder(configFile).Decode(&config)

	bot := WithoutError(tgbotapi.NewBotAPI(config.Key))
	alarms = make(map[string]*time.Timer)
	for _, p_ := range config.Projects {
		p := p_
		println("project: " + p.Id)
		p.Interval = defaultInterval
		alarms[p.Id] = time.NewTimer(p.Interval)
		Execute(func() {
			for {
				<-alarms[p.Id].C
				println(p.Id + " stop response")
				msg := tgbotapi.NewMessage(p.RoomId, p.Id + " stop response")
				bot.Send(msg)
				p.Interval *= 2
				alarms[p.Id].Reset(p.Interval)
			}
		})
	}
}

func main() {
	r := gin.Default()
	r.GET("/ping/:id", func(c *gin.Context) {
		id := c.Params.ByName("id")
		if p, ok := config.getProject(id); ok {
			p.Interval = defaultInterval
		}
		if alarm, ok := alarms[id]; ok {
			alarm.Reset(defaultInterval)
		}
		c.JSON(http.StatusOK, gin.H{
			"message": "pong" + id,
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

type Project struct {
	Id       string        `json:"id"`
	RoomId   int64         `json:"roomid"`
	Interval time.Duration `json:"interval"`
}
type Config struct {
	Projects []*Project `json:"projects"`
	Key      string    `json:"key"`
}

func (c *Config) getProject(id string) (*Project, bool) {
	for _, p := range c.Projects {
		if p.Id == id {
			return p, true
		}
	}
	return &Project{}, false
}


// bot.Debug = true
// log.Printf("Authorized on account %s", bot.Self.UserName)
// u := tgbotapi.NewUpdate(0)
// u.Timeout = 60
// updates := bot.GetUpdatesChan(u)
// for update := range updates {
// 	if update.Message != nil { // If we got a message
// 		log.Printf("[%v] %s", update.Message.Chat.ID, update.Message.Text)

// 		msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
// 		msg.ReplyToMessageID = update.Message.MessageID

// 		bot.Send(msg)
// 	}
// }
