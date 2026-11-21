package mail

import "github.com/antonybholmes/go-sys/jobs"

type (
	Payload struct {
		DataType string `json:"dataType"`
		Data     string `json:"data"`
	}

	MailItem struct {
		Payload   *Payload      `json:"payload,omitempty"`
		Job       *jobs.BaseJob `json:"job,omitempty"`
		Name      string        `json:"name"`
		To        string        `json:"to"`
		From      string        `json:"from"`
		EmailType string        `json:"type"`
		TTL       string        `json:"ttl,omitempty"`
		LinkUrl   string        `json:"linkUrl,omitempty"`
		Mode      string        `json:"mode,omitempty"`
	}

	// MailQueue defines the interface for sending mail items to a queue
	MailQueue interface {
		SendMail(mail *MailItem) error
	}
)
