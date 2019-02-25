package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	p "github.com/DazWilkin/oauth2/particle"

	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/particle"

	"golang.org/x/oauth2"
)

func randState() string {
	b := make([]byte, 32)
	rand.Read(b)
	return base64.StdEncoding.EncodeToString(b)
}
func main() {

	ctx := context.Background()
	conf := &oauth2.Config{
		ClientID:     os.Getenv("ClientID"),
		ClientSecret: os.Getenv("ClientSecret"),
		Scopes:       []string{""},
		Endpoint:     p.Endpoint,
	}

	url := conf.AuthCodeURL(randState(), oauth2.AccessTypeOffline)
	fmt.Printf("Visit the URL for the auth dialog, %v\n", url)

	var code string
	if _, err := fmt.Scan(&code); err != nil {
		log.Fatal(err)
	}
	tok, err := conf.Exchange(ctx, code)
	if err != nil {
		log.Fatal(err)
	}

	core := particle.NewAdaptor(os.Getenv("DeviceID"), tok.AccessToken)
	led := gpio.NewLedDriver(core, "D7")
	work := func() {
		gobot.Every(1*time.Second, func() {
			led.Toggle()
		})
	}
	robot := gobot.NewRobot("spark",
		[]gobot.Connection{core},
		[]gobot.Device{led},
		work,
	)
	robot.Start()
}
