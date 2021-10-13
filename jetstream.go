package nacl

import (
	"fmt"
	"github.com/jar3b/grawt"
	"github.com/nats-io/nats.go"
	"log"
)

var (
	JsClient nats.JetStreamContext
)

func SetupJetStream(host string, port int, credsFile string, closeHandler *grawt.CloseHandler, opts ...nats.JSOpt) error {
	var err error
	if err := SetupNatsWithCreds(host, port, credsFile, closeHandler); err != nil {
		return err
	}

	JsClient, err = NatsClient.JetStream(opts...)
	if err != nil {
		return fmt.Errorf("cannot init JS: %v", err)
	}

	return nil
}

func FinalizeJetStream(subscriptions *[]*nats.Subscription) error {
	return FinalizeNats(subscriptions)
}

func AddOrUpdateStream(cfg *nats.StreamConfig, opts ...nats.JSOpt) (*nats.StreamInfo, error) {
	_, err := JsClient.StreamInfo(cfg.Name)
	if err != nil && err.Error() != "nats: stream not found" {
		return nil, err
	}

	if err == nil {
		sc, err := JsClient.UpdateStream(cfg, opts...)
		if err != nil {
			return nil, err
		}
		return sc, nil
	} else {
		sc, err := JsClient.AddStream(cfg, opts...)
		if err != nil {
			return nil, err
		}
		return sc, nil
	}
}

func AddOrUpdateConsumer(stream string, consumer string, cfg *nats.ConsumerConfig, allowDelete bool, opts ...nats.JSOpt) (*nats.ConsumerInfo, error) {
	ci0, err := JsClient.ConsumerInfo(stream, consumer, opts...)
	if err != nil && err.Error() != "nats: consumer not found" {
		return nil, err
	}

	if err == nil && allowDelete {
		if err := JsClient.DeleteConsumer(stream, consumer, opts...); err != nil {
			return nil, err
		}
	}

	ci, err := JsClient.AddConsumer(stream, cfg, opts...)
	if err != nil {
		if !allowDelete {
			log.Printf("Consumer '%s' was exists", consumer)
			return ci0, nil
		} else {
			return nil, err
		}
	}
	return ci, nil
}
