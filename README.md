# discordbotrus
[![GitHub Build](https://github.com/outdead/discordbotrus/workflows/build/badge.svg)](https://github.com/goutdead/discordbotrus/actions)
[![Go Coverage](https://github.com/outdead/discordbotrus/wiki/coverage.svg)](https://raw.githack.com/wiki/outdead/discordbotrus/coverage.html)
[![Go Report Card](https://goreportcard.com/badge/github.com/outdead/discordbotrus)](https://goreportcard.com/report/github.com/outdead/discordbotrus)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/outdead/discordbotrus)

[Logrus](https://github.com/sirupsen/logrus) hook for [Discord](https://discordapp.com/) using [Discord application](https://discordapp.com/developers/applications/).

## Install

```text
go get github.com/outdead/discordbotrus
```

See [Changelog](CHANGELOG.md) for release details.

## Requirements

Go 1.23 or higher

## Usage

```go
package main

import (
	"log"
	"os"

	"github.com/outdead/discordbotrus"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := discordbotrus.NewDefaultConfig(os.Getenv("LDH_TOKEN"), os.Getenv("LDH_CHANNEL"))

	hook, err := discordbotrus.New(cfg)
	if err != nil {
		log.Fatalf("expected nil got error: %s", err)
	}

	defer hook.Close()

	logger := &logrus.Logger{
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}
	logger.AddHook(hook)

	logger.Info("My spoon is too big")
}
```

### With an existing discordgo session

If you wish to initialize a Discord Hook with an already initialized discordgo session, you can use the SetSession option:

```go
package main

import (
	"log"
	"os"

	"github.com/bwmarrin/discordgo"
	"github.com/outdead/discordbotrus"
	"github.com/sirupsen/logrus"
)

func main() {
	token := os.Getenv("LDH_TOKEN")
	channelID := os.Getenv("LDH_CHANNEL")

	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("expected nil got error: %s", err)
	}

	// In this case, you can use the session without opening a web socket.
	// But to establish a stable connection, it is better to do this.
	if err := session.Open(); err != nil {
		log.Fatalf("open discord session error: %s", err)
	}

	defer session.Close()

	cfg := &discordbotrus.Config{ChannelID: channelID}

	hook, err := discordbotrus.New(cfg, discordbotrus.WithSession(session))
	if err != nil {
		log.Fatalf("expected nil got error: %s", err)
	}

	logger := &logrus.Logger{
		Formatter: new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     logrus.InfoLevel,
	}
	logger.AddHook(hook)

	logger.Info("My spoon is too big")
}
```

## Config

In most cases, the package will not exist on its own, it will be part of an application. For ease of configuration, 
yaml and json tags are affixed to Config. You can use configuration files similar to this:

```yaml
disabled: false # it is possible to disable the hook from the configuration file.
token: "" # required to create own connection to discord bot.
channel_id: "" # required.
format: "json" # supported formats: text, json, embed
min_level: "info"
levels:
  - "error"
  - "warning"
  - "info"
  - "trace"
```

If only min_level is specified, then the hook will fire for all levels above the specified one. If only the levels 
list is specified, then the hook will work only for all listed levels. The parameters min_level and levels are used 
together and the intersection between them is calcalated. If both of them are specified, then the levels below min_level are cut off from all the listed levels

## Tests

To run tests you need to add the environment variables LDH_TOKEN and LDH_CHANNEL.

## License

MIT License, see [LICENSE](LICENSE)
