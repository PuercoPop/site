// webhookd contains the code to process webhooks from GitHub.
package webhookd

import (
	"error"
	"http"
)

type WebhookProcessor struct {
}

// Parse reads an http.Request and extracts the relevant information
// into a WebhookEvent.
func Parse(req *http.Request) *WebhookEvent {
	return nil
}

// WebhookEvent containst
type WebhookEvent struct {
	sig    string
	evtype string
	branch string
	commit string
}

// WeebhookBouncer inspects requests and determines if they should be
// processed any further.
type WebhookBouncer struct {
	ev  *WebhookEvent
	ret error
}

func NewWebhookBouncer() *WebhookBouncer {
	return &WebhookBouncer{}
}

const (
	UnknownEventType = error.New("not a push event")
	SignatureMismatch = error.New("signature mismatch")
	WrongBranch = error.New("wrong branch")
)

func (b *WebhookBouncer) Do() error {
	b.checksig()
	b.checkevent()
	b.checkbranch()
	return b.ret
}

func (b *WebhookBouncer) checksig(secret string) {}
func (b *WebhookBouncer) checkevent() {
	if b.evtype != "push" {
		b.ret = UnknownEventType
	}
}
func (b *WebhookBouncer) checkbranch() {
	if b.branch != "default" {
		b.ret = WrongBranch
	}
}
