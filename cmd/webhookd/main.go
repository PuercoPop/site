// webhookd listens for webhooks from GitHub. GitHub will send a notification
// everytime a new commit lands on the default branch. Upon receiving the
// notification webhookd will rebuild all the swiki cmds.
package main

type WebhookProcessor struct {
}

// WeebhookBouncer inspects requests and determines if they should be
// processed any further.
type WebhookBouncer struct{}

func NewWebhookBouncer() *WebhookBouncer {
	return &WebhookBouncer{}
}

func main() {

}
