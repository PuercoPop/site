// webhookd contains the code to process webhooks from GitHub.
package webhookd

import "http"

type WebhookProcessor struct {
}

// WeebhookBouncer inspects requests and determines if they should be
// processed any further.
type WebhookBouncer struct{}

func NewWebhookBouncer() *WebhookBouncer {
	return &WebhookBouncer{}
}

func (b *WebhookBouncer)validateSignature(req *http.HttpRequest){

}
