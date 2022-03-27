package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestWebhookBouncer(t *testing.T) {
	body, err := os.Open("./webhookbody.sample")
	if err != nil {
		t.Fatalf("Could not open webookbody.sample. %s", err)
	}
	req := httptest.NewRequest(http.MethodPost, "swiki-webhook", body)
	req.Header.Add("content-type", "application/json")
	req.Header.Add("User-Agent", "GitHub-Hookshot/f05835d")
	req.Header.Add("X-GitHub-Delivery", "1eb307f2-ad7f-11ec-98bb-f307568b0602")
	req.Header.Add("X-GitHub-Event", "push")
	req.Header.Add("X-GitHub-Hook-ID", "350114828")
	req.Header.Add("X-GitHub-Hook-Installation-Target-ID", "469844163")
	req.Header.Add("X-GitHub-Hook-Installation-Target-Type", "repository")
	req.Header.Add("X-Hub-Signature", "sha1=ef4ebbf956a7624035958e0e933a9102a84d8bba")
	req.Header.Add("X-Hub-Signature-256", "sha256=17fb2c8e454cf328f1e73f46f113d10a3306960b03751c0c2f3d18911d31dac9")

	bouncer := NewWebhookBouncer()
	err = bouncer.Do(req)
	if err != nil {

	}
	t.Run("it detects events of interest", func(t *testing.T) {})
	t.Run("fails if the signature doesnt match", func(t *testing.T) {})
	t.Run("fails if the event is not push", func(t *testing.T) {})
	t.Run("fails if the push is not to the default branch", func(t *testing.T) {})
}
