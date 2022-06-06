// webhookd listens for webhooks from GitHub. GitHub will send a notification
// everytime a new commit lands on the default branch. Upon receiving the
// notification webhookd will rebuild all the swiki cmds.
package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/PuercoPop/site/webhookd"
)

func handleWebhook(w http.ResponseWriter, r *http.Request) {
	bouncer := webhookd.NewWebhookBouncer()
	err := bouncer.Do()
	if err != nil {
		fmt.Println("go build goes here")
	}
}

func main() {

	http.HandleFunc("/", handleWebhook)
	log.Fatal(http.ListenAndServe(":8080", nil))

}
