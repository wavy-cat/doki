package cmd

import "github.com/wavy-cat/doki/internal/app"

func Execute() error {
	return NewRootCommand(app.DefaultRunner()).Execute()
}
