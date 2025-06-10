package terminal

import (
	"github.com/desertbit/grumble"
)

//func Println() == app.Println()

func NewTerminal(name string) *grumble.App {
	return grumble.New(&grumble.Config{
		Name:        name,
		Description: name,

		Prompt: ">",
	})
}

func AddCommand(app *grumble.App, name string, help string, args []string, listargs string, exec func(*grumble.Context) error) {
	app.AddCommand(&grumble.Command{
		Name: name,
		Help: help,

		Args: func(a *grumble.Args) {
			if listargs != "" {
				a.StringList(listargs, "args")
			}

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

func SetPrompt(app *grumble.App, prompt string) {
	app.SetPrompt(prompt)
}
