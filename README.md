# discordbotrus
[![Go Report Card](https://goreportcard.com/badge/github.com/outdead/discordbotrus)](https://goreportcard.com/report/github.com/outdead/discordbotrus)
[![Build Status](https://travis-ci.org/outdead/discordbotrus.svg?branch=master)](https://travis-ci.org/outdead/discordbotrus)
[![Coverage](https://gocover.io/_badge/github.com/outdead/discordbotrus?0 "coverage")](https://gocover.io/github.com/outdead/discordbotrus)
[![GoDoc](https://img.shields.io/badge/godoc-reference-blue.svg)](https://godoc.org/github.com/outdead/discordbotrus)

A Discord Bot hook for Logrus.

## Install

```text
go get github.com/outdead/discordbotrus
```

Or use dependency manager such as dep or vgo.

## Usage

```go
package main

import (
	"io/ioutil"
	"os"
	"log"

	"github.com/outdead/discordbotrus"
	"github.com/sirupsen/logrus"
)

func main() {
	cfg := hook.NewDefaultConfig(os.Getenv("LDH_TOKEN"), os.Getenv("LDH_CHANNEL"))
	hooker, err := hook.New(cfg)
	if err != nil {
		log.Fatalf("expected nil got error: %s", err)
	}

	logger := &logrus.Logger{Out: ioutil.Discard, Formatter: new(logrus.JSONFormatter), Hooks: make(logrus.LevelHooks), Level: logrus.InfoLevel}
	logger.AddHook(hooker)

	defer func() {
		if err := hooker.Close(); err != nil {
			log.Fatalf("expected nil got error: %s", err)
		}
	}()

	logger.Info("My spoon is too big")	
}
```

## Config
```yaml
disabled: false
token: "" # required to create own connection to discord bot.
channel_id: "" # required.
format: "json"
min_level: "info"
levels:
- "error"
- "warning"
- "info"
- "trace"
```

## Tests

To run tests you need to add the environment variables LDH_TOKEN and LDH_CHANNEL.
