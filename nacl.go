package nacl

import (
	"fmt"
	"github.com/jar3b/grawt"
	"github.com/nats-io/nats.go"
	"sync"
	"time"
)

var (
	NatsClient *nats.Conn
	natsLock   = sync.Mutex{}
)

func SetupNatsWithCreds(host string, port int, credsFile string, closeHandler *grawt.CloseHandler) error {
	natsLock.Lock()
	defer natsLock.Unlock()
	var err error

	// connect
	NatsClient, err = nats.Connect(
		fmt.Sprintf("nats://%s:%d", host, port),
		nats.UserCredentials(credsFile),
		nats.ClosedHandler(func(conn *nats.Conn) {
			if closeHandler != nil {
				closeHandler.Halt(nil)
			}
		}),
		nats.MaxReconnects(5),
		nats.ReconnectWait(time.Second*2),
	)
	if err != nil {
		return fmt.Errorf("cannot connect to NATS: %v", err)
	}

	return nil
}

func SetupNats(host string, port int, user string, pass string, closeHandler *grawt.CloseHandler) error {
	natsLock.Lock()
	defer natsLock.Unlock()
	var err error

	// init connection string
	connectionString := fmt.Sprintf("%s:%d", host, port)
	if user != "" && pass != "" {
		connectionString = fmt.Sprintf("%s:%s@%s", user, pass, connectionString)
	}
	connectionString = fmt.Sprintf("nats://%s", connectionString)

	// connect
	NatsClient, err = nats.Connect(
		connectionString,
		nats.ClosedHandler(func(conn *nats.Conn) {
			if closeHandler != nil {
				closeHandler.Halt(nil)
			}
		}),
		nats.MaxReconnects(5),
		nats.ReconnectWait(time.Second*2),
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
