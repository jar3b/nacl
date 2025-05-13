package nacl

import (
	"fmt"
	"github.com/jar3b/grawt"
	"github.com/nats-io/nats.go"
	"log"
	"sync"
	"time"
)

var (
	NatsClient *nats.Conn
	natsLock   = sync.Mutex{}
)

func SetupNatsWithCreds(
	natsURL string,
	credsFile string,
	appName string,
	closeHandler *grawt.CloseHandler,
	errorHandler nats.ErrHandler,
) error {
	natsLock.Lock()
	defer natsLock.Unlock()
	var err error

	// connect
	options := []nats.Option{
		nats.ClosedHandler(func(conn *nats.Conn) {
			if closeHandler != nil {
				closeHandler.Halt(nil)
			}
		}),
		nats.MaxReconnects(5),
		nats.ReconnectWait(time.Second * 2),
		nats.Name(appName),
		nats.ConnectHandler(func(conn *nats.Conn) {
			log.Printf("NATS Connected to %v\n", conn)
		}),
		nats.DisconnectErrHandler(func(conn *nats.Conn, err error) {
			log.Printf("NATS Disconnected: conn=%v, err=%v\n", conn, err)
		}),
	}

	if errorHandler != nil {
		options = append(options, nats.ErrorHandler(errorHandler))
	}

	if credsFile != "" {
		options = append(options, nats.UserCredentials(credsFile))
	}

	NatsClient, err = nats.Connect(
		fmt.Sprintf(natsURL),
		options...,
	)
	if err != nil {
		return fmt.Errorf("cannot connect to NATS: %v", err)
	}

	return nil
}

func SetupNatsV2(
	natsURL string,
	credsFile string,
	appName string,
	closeHandler *grawt.CloseHandler,
	natsOptions ...nats.Option,
) error {
	natsLock.Lock()
	defer natsLock.Unlock()
	var err error

	// connect
	options := []nats.Option{
		nats.ClosedHandler(func(conn *nats.Conn) {
			if closeHandler != nil {
				closeHandler.Halt(nil)
			}
		}),
		nats.MaxReconnects(5),
		nats.ReconnectWait(time.Second * 2),
		nats.Name(appName),
	}
	options = append(options, natsOptions...)

	if credsFile != "" {
		options = append(options, nats.UserCredentials(credsFile))
	}

	NatsClient, err = nats.Connect(
		fmt.Sprintf(natsURL),
		options...,
	)
	if err != nil {
		return fmt.Errorf("cannot connect to NATS: %v", err)
	}

	return nil
}

func FinalizeNats(subscriptions *[]*nats.Subscription) error {
	natsLock.Lock()
	defer natsLock.Unlock()
	if NatsClient == nil {
		return fmt.Errorf("stan client is not initialized")
	}

	if subscriptions != nil {
		for _, subscription := range *subscriptions {
			_ = subscription.Unsubscribe()
		}
	}

	NatsClient.Close()

	return nil
}
