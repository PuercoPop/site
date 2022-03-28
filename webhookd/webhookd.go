// webhookd contains the code to process webhooks from GitHub.
package webhookd

import (
	"errors"
	"encoding/json"
	"net/http"
)

type WebhookProcessor struct {
}

// Parse reads an http.Request and extracts the relevant information
// into a WebhookEvent.
func Parse(req *http.Request) (*WebhookEvent, error) {
	ev := &WebhookEvent{}
	ev.evtype = req.Header.Get("X-GitHub-Event")
	ev.sig = req.Header.Get("X-Hub-Signature-256")
	body := struct {
		ref    string `json:"ref"`
		commit string `json:"after"`
	}{}
	err := json.NewDecoder(req.Body).Decode(&body)
	defer req.Body.Close()
	if err != nil {
		return nil, err
	}
	ev.branch = body.ref
	ev.commit = body.commit
	return ev, nil
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

var (
	UnknownEventType  error = errors.New("not a push event")
	SignatureMismatch error = errors.New("signature mismatch")
	WrongBranch       error = errors.New("wrong branch")
)

func (b *WebhookBouncer) Do() error {
	b.checksig("iou")
	b.checkevent()
	b.checkbranch()
	return b.ret
}

func (b *WebhookBouncer) checksig(secret string) {}
func (b *WebhookBouncer) checkevent() {
	if b.ev.evtype != "push" {
		b.ret = UnknownEventType
	}
}
func (b *WebhookBouncer) checkbranch() {
	if b.ev.branch != "default" {
		b.ret = WrongBranch
	}
}
