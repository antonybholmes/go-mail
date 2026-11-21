package mail

import "github.com/antonybholmes/go-sys/jobs"

type (
	Payload struct {
		DataType string `json:"dataType"`
		Data     string `json:"data"`
	}

	MailItem struct {
		Payload *Payload      `json:"payload,omitempty"`
		Job     *jobs.BaseJob `json:"job,omitempty"`
		// Name of the user, if applicable.
		UserID string `json:"userId,omitempty"`
		//Name      string    `json:"name,omitempty"`
		To        string    `json:"to"`
		From      string    `json:"from"`
		EmailType EmailType `json:"type"`
		TTL       string    `json:"ttl,omitempty"`
		LinkUrl   string    `json:"linkUrl,omitempty"`
		Mode      string    `json:"mode,omitempty"`
	}

	// MailQueue defines the interface for sending mail items to a queue
	MailQueue interface {
		SendMail(mail *MailItem) error
	}
)

func NewJobSubmittedMailItem(job *jobs.BaseJob, userID, to, from string) *MailItem {
	return &MailItem{
		UserID:    userID,
		Job:       job,
		To:        to,
		From:      from,
		EmailType: EmailTypeJobSubmitted,
	}
}
