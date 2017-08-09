package command

import (
	"github.com/alecthomas/kingpin"
	"net"
	"net/url"
	"log"
	"github.com/gorilla/websocket"
	"strconv"
	"net/http"
	"strings"
)

//部署客户端工具命令
type ReceiverCommand struct {
	Name string
	Help string

	Port   *int
	IP     *net.IP
	WSPath *string
	Token  *string
}

const BEAR_CHAT_BOOT_URI = "https://hook.bearychat.com/=bwBlp/incoming/"

func (c *ReceiverCommand) Exec(ctx *Context) {

	u := url.URL{Scheme: "ws", Host: c.IP.String() +
		":" +
		strconv.Itoa(*c.Port), Path: *c.WSPath}
	log.Printf("connecting to %s", u.String())

	conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		log.Fatal("dial:", err)
	}
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			return
		}
		log.Printf("recv: %s", message)

		sr := strings.NewReader(string(message))
		http.Post(BEAR_CHAT_BOOT_URI + *c.Token, "application/json", sr)

		//bs, _ := ioutil.ReadAll(resp.Body)
		//
		//defer resp.Body.Close()
		//log.Print(string(bs))
	}

}

func (c *ReceiverCommand) configFlags(cmd *kingpin.Application) {
	sub := cmd.Command(c.Name, c.Help)
	c.Port = sub.Flag("port", "port").Required().Int()
	c.IP = sub.Flag("ip", "ip").Required().IP()
	c.WSPath = sub.Flag("ws-path", "http-ctx-path").
		Default("/_notify_bridge/ws").String()
	c.Token = sub.Flag("token", "robot token").Required().String()
}

func (c *ReceiverCommand) CmdName() string {
	return c.Name
}

func NewReceiverCommand() Command {
	c := new(ReceiverCommand)
	c.Name = "receiver"
	c.Help = "公网出口的receiver"
	return c
}
