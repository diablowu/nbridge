package command

import (
	"github.com/alecthomas/kingpin"
	"log"
	"net/http"
	"github.com/gorilla/websocket"
	"strconv"
	"github.com/go-macaron/binding"
	"fmt"
	"gopkg.in/macaron.v1"
	"encoding/json"
)

//agent模式的命令
type ServerCommand struct {
	Name   string
	Help   string
	Port   *int
	WSPath *string
	IMPath *string
}

type BearChatIncoming struct {
	Text         string `json:"text" form:"text"`
	Notification string `json:"notification"`
	MarkDown     bool `json:"markdown"`
	Channel      string `json:"channel"`
	User         string `json:"user"`
}

var wsConn *websocket.Conn

var upgrader = websocket.Upgrader{}

func (c *ServerCommand) Exec(ctx *Context) {

	m := macaron.Classic()

	m.Any(*c.WSPath, wsHandle)

	m.Post(*c.IMPath,
		binding.BindIgnErr(BearChatIncoming{}),
		func(im BearChatIncoming) string {
			ss, _ := json.Marshal(im)
			wsConn.WriteMessage(1, ss)
			return fmt.Sprintf("%s", im)
		})

	log.Print("starting to http server")

	err := http.ListenAndServe("0.0.0.0:"+strconv.Itoa(*c.Port), m)
	if err != nil {
		log.Fatal(err)
	}

}

func wsHandle(w http.ResponseWriter, r *http.Request) {
	var err error
	wsConn, err = upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade: ", err)
		return
	}
}

func (c *ServerCommand) configFlags(cmd *kingpin.Application) {
	sub := cmd.Command(c.Name, c.Help)
	c.WSPath = sub.Flag("ws-path", "websocket path").
		Default("/_notify_bridge/ws").String()
	c.IMPath = sub.Flag("im-path", "http post im path").
		Default("/_notify_bridge/im").String()
	c.Port = sub.Flag("port", "port").Required().Int()
}

func (c *ServerCommand) CmdName() string {
	return c.Name
}

func NewServerCommand() Command {
	c := new(ServerCommand)
	c.Name = "server"
	c.Help = "server bridge"
	return c
}
