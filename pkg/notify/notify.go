package notify

import (
	"fmt"
	"github.com/rs/zerolog/log"
	"github.com/wadeling/linear-webhook/api/linear"
	"github.com/wadeling/linear-webhook/pkg/webhook"
)

var creators = make(map[Type]Creator)

// Config contain info that need to create notifier
type Config struct {
	Type      Type          // specify which notify that would be created
	Url       string        // notify server address
	LinConfig linear.Config // linear config of api
}

// Creator create a notifier
type Creator func(config Config) (Notifier, error)

// Notifier define behavior when recv webhook
type Notifier interface {
	Deliver(payload webhook.PayLoad) error
}

// Register record notifier
func Register(name Type, creator Creator) error {
	if creator == nil {
		return fmt.Errorf("could not register nil creator")
	}
	if _, dup := creators[name]; dup {
		return fmt.Errorf("could not register duplicate creator: %v", name)
	}
	creators[name] = creator
	return nil
}

// Open create specify notifier
func Open(cfg Config) (Notifier, error) {
	creator, ok := creators[cfg.Type]
	if !ok {
		return nil, fmt.Errorf("unknown creator %q (forgotten configuration or import?)", cfg.Type)
	}
	return creator(cfg)
}

func DeliverAll(config Config, payload webhook.PayLoad) error {
	for k := range creators {
		// create notifier
		config.Type = k
		n, err := Open(config)
		if err != nil {
			log.Err(err).Str("type", string(k)).Msg("failed to create notifier")
			continue
		}

		// post msg
		err = n.Deliver(payload)
		if err != nil {
			log.Err(err).Str("type", string(k)).Msg("failed to deliver msg")
		} else {
			log.Info().Str("type", string(k)).Msg("deliver msg ok")
		}
	}
	return nil
}
