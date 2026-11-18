package mailqueue

import (
	"sync"

	email "github.com/antonybholmes/go-mail"
)

var (
	queue email.MailQueue
	once  sync.Once
)

func InitMailQueue(q email.MailQueue) {
	once.Do(func() {
		queue = q
	})
}

// SendMail adds an email to the mail queue
func SendMail(email *email.MailItem) error {
	return queue.SendMail(email)
}
