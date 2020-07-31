package terminal

import (
	"context"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/stripe/stripe-cli/pkg/config"
	"github.com/stripe/stripe-cli/pkg/proxy"
	"github.com/stripe/stripe-cli/pkg/stripe"
	"github.com/stripe/stripe-cli/pkg/validators"
	"github.com/stripe/stripe-cli/pkg/websocket"
)

func processWebhookEvent(msg websocket.IncomingMessage) {
	fmt.Println(msg.WebhookEvent.WebhookID)
	fmt.Println(msg.WebhookEvent.Type)
	fmt.Println(msg.WebhookEvent.EventPayload)
	// TODO: forward events to reader
}

// Tunnel tunnels PaymentIntent events to a Terminal Reader on the local network
func Tunnel(cfg *config.Config) error {
	key, err := cfg.Profile.GetAPIKey(false)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	err = validators.APIKeyNotRestricted(key)
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	/* TODO: connect to reader before listening for events
	tsCtx := SetTerminalSessionContext(cfg)
	tsCtx, err = p400.RegisterAndActivateReader(tsCtx)

	if err != nil {
		if err.Error() == promptui.ErrInterrupt.Error() {
			os.Exit(1)
		} else {
			return fmt.Errorf(err.Error())
		}
	}

	fmt.Println("Connected to reader!")
	fmt.Println(tsCtx.DeviceInfo)
	*/

	events := []string{"payment_intent.created", "payment_intent.canceled"}
	p := proxy.New(&proxy.Config{
		DeviceName:       cfg.Profile.DeviceName,
		Key:              key,
		APIBaseURL:       stripe.DefaultAPIBaseURL,
		Log:              log.StandardLogger(),
		WebSocketFeature: "webhooks",
		EventHandler:     processWebhookEvent,
	}, events)

	err = p.Run(context.Background())
	if err != nil {
		return fmt.Errorf(err.Error())
	}

	return nil
}
