package internal

import (
	"fmt"
	"os"
	"os/signal"
	"suplex/internal/config"
	"suplex/internal/log"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

type Suplex struct {
	// Session
	Session *discordgo.Session

	//
	Handler *CommandHandler

	// Database
	//Database *database.DB

	// Config
	Config *config.Configuration

	// Logger
	Logger *log.Logger

	// Runtime Signals
	exitSignal    chan os.Signal
	restartSignal chan os.Signal
}

func NewSuplex(cfg *config.Configuration, logger *log.Logger) (*Suplex, error) {
	var suplex = &Suplex{
		Session: nil,
		Handler: &CommandHandler{cmdMap: map[string]*Command{}},
		//Database: db,
		Config:        cfg,
		Logger:        logger,
		exitSignal:    make(chan os.Signal),
		restartSignal: make(chan os.Signal),
	}

	session, err := discordgo.New("Bot " + suplex.Config.Token)
	if err != nil {
		return nil, err
	}

	suplex.Session = session

	return suplex, err
}

func (bot *Suplex) Start() {
	signal.Notify(
		bot.exitSignal,
		syscall.SIGINT,
		syscall.SIGTERM,
	)
	signal.Notify(
		bot.restartSignal,
		syscall.SIGUSR1,
		syscall.SIGUSR2,
	)

	bot.run()
}

func (bot *Suplex) run() {
	defer bot.handlePanic()

	if err := bot.Session.Open(); err != nil {
		panic(err)
	}

	for {
		select {
		case <-bot.exitSignal:
			bot.Stop()
			return

		case <-bot.restartSignal:
			return
		}

	}
}

func (bot *Suplex) Restart() {}

func (bot *Suplex) Stop() {
	bot.Logger.Info("(*Suplex).Stop()", "stopping Bot...")

	bot.Session.Close()
	//db
	bot.Logger.Close()

	os.Exit(0)
}

func (bot *Suplex) handlePanic() {
	if r := recover(); r != nil {
		bot.Logger.Error("(*Suplex).handlePanic()", "Panic-Reason", fmt.Sprintf("%#v", r))
		bot.Logger.Error("(*Suplex).handlePanic()", "trying to restart...")
		//? restart
	}
}

func (bot *Suplex) ReregisterAllCommands() {
	//unregister all
	cmds, err := bot.Session.ApplicationCommands(
		bot.Session.State.User.ID,
		"",
	)
	if err != nil {
		//?
	}

	if len(cmds) > 0 {
		for _, command := range cmds {
			bot.Session.ApplicationCommandDelete(
				bot.Session.State.User.ID,
				"",
				command.ID,
			)
		}
	}

	//
	for _, cmd := range bot.Handler.cmdMap {
		if _, err := bot.Session.ApplicationCommandCreate(
			bot.Session.State.User.ID,
			"",
			cmd.ApplicationCommand,
		); err != nil {
			bot.Logger.Error("(*Suplex).ReregisterAllCommands()", "skipping because of an error...", cmd.Name, err.Error())
			continue
		}
	}

}
