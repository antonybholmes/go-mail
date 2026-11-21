package mail

import (
	"context"
	"encoding/json"

	"github.com/antonybholmes/go-sys"

	"github.com/antonybholmes/go-sys/log"
	"github.com/panjf2000/ants"
	"github.com/redis/go-redis/v9"
)

type (
	ValkeyHandler struct {
		Addr         string
		User         string
		Password     string
		EmailHandler MailQueue
		pool         *ants.Pool
		backoff      *sys.Backoff
	}
)

func NewValkeyHandler(addr string,
	user string,
	password string,
	emailHandler MailQueue,
	pool *ants.Pool) *ValkeyHandler {
	return &ValkeyHandler{
		Addr:         addr,
		User:         user,
		Password:     password,
		EmailHandler: emailHandler,
		pool:         pool,
		backoff:      sys.NewDefaultBackoff(),
	}
}

func (h *ValkeyHandler) Consume(ctx context.Context) {
	//env.Reload()
	//env.Load("consts.env")
	//env.Load("version.env")

	//var ctx = context.Background()

	var rdb = redis.NewClient(&redis.Options{
		Addr:     h.Addr,
		Username: h.User,
		Password: h.Password, // no password set
		DB:       0,          // use default DB
	})

	subscriber := rdb.Subscribe(ctx, EmailQueueChannel)

	for {
		msg, err := subscriber.ReceiveMessage(ctx)

		if err != nil {
			panic(err)
		}

		err = h.pool.Submit(func() {
			err := h.processMessages(msg)

			if err != nil {
				log.Debug().Msgf("error sending email: %v", err)
			}
		})

		if err != nil {
			log.Debug().Msgf("failed to submit task: %v", err)
		}

	}
}

func (h *ValkeyHandler) processMessages(msg *redis.Message) error {
	var m MailItem

	err := json.Unmarshal([]byte(msg.Payload), &m)

	if err != nil {
		log.Debug().Msgf("email error")
	}

	//log.Debug().Msgf("email %s %v", msg.Payload, qe.EmailType)

	return h.EmailHandler.SendMail(&m)
}
