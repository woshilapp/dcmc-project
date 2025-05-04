package terminal

import (
	"github.com/desertbit/grumble"
)

//func Println() == app.Println()

func NewTerminal(name string) *grumble.App {
	return grumble.New(&grumble.Config{
		Name:        name,
		Description: "dcmc-project terminal",

		Prompt: ">",
	})
}

func AddCommand(app *grumble.App, name string, help string, args []string, exec func(*grumble.Context) error) {
	app.AddCommand(&grumble.Command{
		Name: name,
		Help: help,

		Args: func(a *grumble.Args) {
			for _, arg := range args {
				a.String(arg, "arg")
			}
		},

		Run: func(c *grumble.Context) error {
			err := exec(c)

			return err
		},
	})
}

func AddMultiArgCommand(app *grumble.App, name string, help string, arg string, exec func(*grumble.Context) error) {
	app.AddCommand(&grumble.Command{
		Name: name,
		Help: help,

		Args: func(a *grumble.Args) {
			a.StringList(arg, "arg")
		},

		Run: func(c *grumble.Context) error {
			err := exec(c)

			return err
		},
	})
}
