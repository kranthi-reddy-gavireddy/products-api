package cmd

import "github.com/kranthi-reddy-gavireddy/message-broker/internal/app"

func main() {
	app := app.New()
	app.Start()
}
