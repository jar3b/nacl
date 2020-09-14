# nacl

Client helper for NATS/STAN

## examples

connect to NATS

```go
package main

import (
    "fmt"
    "github.com/jar3b/grawt"
    "github.com/jar3b/nacl"
    "github.com/kelseyhightower/envconfig"
)

func main() {
    // get params
    conf := config.NewConfig()

	// init waiter (github.com/jar3b/grawt)
	var waiter = grawt.NewWaiter()

	// init nats
	var subscriptions []*nacl.NatsSubscription
	if err := nacl.SetupNats(conf.Nats.Host, conf.Nats.Port, conf.Nats.User, conf.Nats.Pass,
		// handler called before app closed
        // we need to terminate sub's properly and sometimes doing another actions (finalizers)
        waiter.AddCloseHandler(func() {
			nacl.FinalizeNats(&subscriptions)
		}, false),
	); err != nil {
		waiter.Halt(fmt.Errorf("cannot connect to nats: %v", err))
	}
	defer nacl.NatsClient.Close()
	
    // here we add some subscriptions and wait (using blocking call)
    // ...
}
```