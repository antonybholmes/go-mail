package mailqueue

import (
	"sync"

	"github.com/antonybholmes/go-jobs"
)

var (
	queue jobs.MailQueue
	once  sync.Once
)

func InitMailQueue(q jobs.MailQueue) {
	once.Do(func() {
		queue = q
	})
}

// SendMail adds an email to the mail queue
func Instance() jobs.MailQueue {
	return queue
}
