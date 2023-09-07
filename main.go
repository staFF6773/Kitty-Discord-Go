package main

import (
	"suplex/internal"
	"suplex/internal/command"
	"suplex/internal/config"
	"suplex/internal/event"
	"suplex/internal/log"

	"github.com/bwmarrin/discordgo"
)

const (
	configFilePath = "./config.json"
)

var (
	cfg    *config.Configuration
	logger *log.Logger
	//db
)

func init() {
	var err error

	// Config
	cfg, err = config.ParseConfigFromJSONFile(configFilePath)
	if err != nil {
		panic(err)
	}

	// Logger
	logger, err = log.NewLogger(
		cfg.Logging.Method,
		cfg.Logging.File,
		cfg.Logging.Level,
	)
	if err != nil {
		panic(err)
	}

	//?db

}

func main() {

	suplex, err := internal.NewSuplex(cfg, logger)
	if err != nil {
		panic(err)
	}

	//?
	suplex.Session.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)

	/*
	 * Register Events
	 */
	suplex.Session.AddHandler(event.NewReadyHandler(suplex).Exec)
	suplex.Session.AddHandler(event.NewMessageAddHandler(suplex).Exec)
	suplex.Session.AddHandler(event.NewInteractionAdd(suplex).Exec)
	suplex.Session.AddHandler(event.NewDefaultHandler(suplex).Exec)

	/*
	 * Register Commands
	 */

	suplex.Handler.RegisterCommand(command.NewTestCommand(suplex))
	suplex.Handler.RegisterCommand(command.NewNekoCommand(suplex))

	// Start Bot
	suplex.Start()

	// Stop Bot
	suplex.Stop()
}

/*
func main() {
    // Register the slash command with Discord
    command := &discordgo.ApplicationCommand{
        Name:        "ping",
        Description: "Responds with Pong!",
    }
    _, err = discord.ApplicationCommandCreate(discord.State.User.ID, "", command)
    if err != nil {
        log.Fatal("Error creating slash command: ", err)
    }

}

func handleSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
    if i.ApplicationCommandData().Name == "ping" {
        s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{
                Content: "Pong!",
            },
        })
    }
}
*/
