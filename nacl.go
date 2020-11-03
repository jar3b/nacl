package nacl

import (
	"fmt"
	"github.com/jar3b/grawt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/stan.go"
	"time"
)

var (
	NatsClient *nats.Conn
	StanClient stan.Conn
)

type (
	Msg              = stan.Msg
	NatsMsg          = nats.Msg
	Subscription     = stan.Subscription
	NatsSubscription = nats.Subscription
)

func SetupNats(host string, port int, user string, pass string, closeHandler *grawt.CloseHandler) error {
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

func SetupStan(clusterName string, clientId string, host string, port int, user string, pass string, closeHandler *grawt.CloseHandler) (err error) {
	if err = SetupNats(host, port, user, pass, closeHandler); err != nil {
		return err
	}

	// stan connection
	StanClient, err = stan.Connect(
		clusterName,
		clientId,
		stan.NatsConn(NatsClient),
	)
	if err != nil {
		return fmt.Errorf("cannot connect to STAN: %v", err)
	}

	return nil
}

func FinalizeStan(subscriptions *[]Subscription) error {
	if StanClient == nil {
		return fmt.Errorf("stan client is not initialized")
	}

	if subscriptions != nil {
		for _, subscription := range *subscriptions {
			_ = subscription.Unsubscribe()
		}
	}

	if err := StanClient.Close(); err != nil {
		return fmt.Errorf("cannot disconnect from STAN")
	}

	return nil
}

func FinalizeNats(subscriptions *[]*NatsSubscription) error {
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
