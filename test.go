//we finally find it, let's fking go!

package main

import (
	"time"

	"github.com/desertbit/grumble"
)

func main1() {
	app := grumble.New(&grumble.Config{
		Name:        "app",
		Description: "short app description",

		Prompt: ">",

		Flags: func(f *grumble.Flags) {
			f.String("d", "directory", "DEFAULT", "set an alternative directory path")
			f.Bool("v", "verbose", false, "enable verbose mode")
		},
	})

	app.AddCommand(&grumble.Command{
		Name:    "daemon",
		Help:    "run the daemon",
		Aliases: []string{"run"},

		Flags: func(f *grumble.Flags) {
			f.Duration("t", "timeout", 2*time.Second, "timeout duration")
		},

		Args: func(a *grumble.Args) {
			a.String("service", "which service to start", grumble.Default("server"))
		},

		Run: func(c *grumble.Context) error {
			// Parent Flags.
			c.App.Println("directory:", c.Flags.String("directory"))
			c.App.Println("verbose:", c.Flags.Bool("verbose"))
			// Flags.
			c.App.Println("timeout:", c.Flags.Duration("timeout"))
			// Args.
			c.App.Println("service:", c.Args.String("service"))
			return nil
		},
	})

	app.AddCommand(&grumble.Command{
		Name:    "startLog",
		Help:    "start a goroutine",
		Aliases: []string{"sl"},

		Args: func(a *grumble.Args) {
			a.StringList("output", "output string", grumble.Default("log"))
		},

		Run: func(c *grumble.Context) error {
			go func() {
				for {
					time.Sleep(1 * time.Second)
					c.App.Println("[SYSLOG]", c.Args.StringList("output"))
				}
			}()

			return nil
		},
	})

	go func() {
		for {
			time.Sleep(1 * time.Second)
			app.Println("background")
		}
	}()

	err := app.Run()

	if err != nil {
		println(err)
	}
}
