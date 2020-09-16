package main

import (
	"fmt"
	"github.com/go-impatient/gaia/cmd/gaia/create"
	"github.com/go-impatient/gaia/cmd/gaia/server"
	"github.com/urfave/cli/v2"
	"os"
	"time"
)

const Version = "v1.0.0"

var usageStr = `
  ________       .__        
 /  _____/_____  |__|____   
/   \  ___\__  \ |  \__  \  
\    \_\  \/ __ \|  |/ __ \_
 \______  (____  /__(____  /
        \/     \/        \/
`

func run() {
	app := cli.NewApp()
	app.Name = "gaia"
	app.Version = Version
	app.Compiled = time.Now()
	app.Authors = []*cli.Author{
		&cli.Author{
			Name:  "moocss",
			Email: "moocss@gmail.com",
		},
	}
	app.Copyright = "(c) 2020 moocss"
	app.Usage = "一个轻量级的应用服务"
	app.UsageText = usageStr
	app.UseShortOptionHandling = true
	app.EnableBashCompletion = true
	app.Commands = cli.Commands{
		server.Cmd,
		create.Cmd,
	}
	app.Before = func(c *cli.Context) error {
		fmt.Fprintf(c.App.Writer, "brace for impact\n")
		return nil
	}
	app.After = func(c *cli.Context) error {
		fmt.Fprintf(c.App.Writer, "did we lose anyone?\n")
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func main() {
	run()
}
