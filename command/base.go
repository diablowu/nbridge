/////////////////////////////////////////////////////
// 子命令的基础结构体
/////////////////////////////////////////////////////
package command

import (
	"github.com/alecthomas/kingpin"
	"log"
)

const (
	VERSION   = "0.0.1"
	PROG_NAME = "notify-bridge"
	INTRO     = "notify-bridge"
)

type CommandManager struct {
	allowed map[string](Command)
}

func (m *CommandManager) regist(s func() (Command)) {
	if m.allowed == nil {
		log.Fatal("CommandManager not initialized")
	}
	cmd := s()
	m.allowed[cmd.CmdName()] = cmd
}

func (m *CommandManager) get(cmd string) (Command) {
	return m.allowed[cmd]
}

func (m *CommandManager) configFlags(app *kingpin.Application) {
	for k, _ := range m.allowed {
		m.allowed[k].configFlags(app)
	}
}

var manager = new(CommandManager)

func init() {
	manager.allowed = make(map[string](Command), 5)

	manager.regist(NewReceiverCommand)
	manager.regist(NewServerCommand)

}

// 运行时参数和相关上下文对象
type Context struct {
	cmd string
}

//cli实例
type App struct {
	ctx *Context
	cli *kingpin.Application
	cmd *Command
}

//基于命令行参数进行初始化
func (app *App) Init(args []string) {
	app.cli = kingpin.New(PROG_NAME, INTRO)

	app.cli.Version(VERSION)

	//配置支持的cmd参数
	manager.configFlags(app.cli)

	app.ctx.cmd = kingpin.MustParse(app.cli.Parse(args[1:]))

}

//执行
func (app *App) Run() {
	manager.get(app.ctx.cmd).Exec(app.ctx)
}

func NewApp() *App {
	app := new(App)
	app.ctx = new(Context)
	return app
}

//子命令的接口
type Command interface {
	CmdName() string

	Exec(c *Context)

	//配置子命令的flags
	configFlags(cmd *kingpin.Application)
}
