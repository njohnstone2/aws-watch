package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestHandler(t *testing.T) {
// 	t.Run("Parse EC2 Message", func(t *testing.T) {
// 		event := events.CloudwatchLogsEvent{
// 			AWSLogs: events.CloudwatchLogsRawData{
// 				Data: "H4sIADzdXWQAA11TXYvbMBB8768Ieo5TyZY/38wlPUy5ctRpC22OoLM3RuCvSvKlach/70Z2fKFGYKMZ7cys1mfSgNaigu2pB5KQdbpN90+bPE8fN2RJumMLCreZ63E/CKOYMhe36656VN3QI1KYfSNaLNBAa/YI6BHPjQLR/Hd0/1B3Q7lVQtb7QTsgtHEY8vXwqgsleyO79pOsDShNkl/EQNNvlawq9PBiq27eUOWKnYkssbgb+pxxyiPuBa4X8IgF1A98Pwo58wM38LwwZIx6sRfygDLmUx77YchR00hMbkSDIVgQeZzFbozHguWtI1j+vCNwVfyOhtDajiQ7wlY02pHljgwaVFYiKs0JEeQa7KHlZOnTN0QtrVeyLWQv6qy0WJqtU5beHmZXaplCjQr4TsRRJ1I0SXLfvuQq+RH+oOkaxiNF0Q2tmUrfc28wZvkMp5t2nt1pj/Jzli+iGd3PApfllH8rJ8ilrudQ36HhltLEjRLq/bQFLC3vBlVMNQp3JRrxt2sxyaromnfWrPOAE2Igh2JQ2EI7UKPro/4K1a3f85xYTFuJ7DktS4XRphtxV96Kz0HSClXGvD/yRdbiOLWitrCC3wPe+rNQaOI6ZuPFVVft2dd17Czbbq9hHs377izEUEqzsIks+a0vpi7jlyMYw8XY2ES02nethk1tf5NJdTIznRJlGUd4YU5YhKXDuaBO5MalIw5wYIEXMTgcrNBegRnsqBg1wM3mVEVXzvsUiNeihEOFHi7k8nL58A855N577QMAAA==",
// 			},
// 		}

// 		err := handler(event)
// 		if err != nil {
// 			t.Fatal("Error failed to trigger with an invalid request")
// 		}
// 	})
// }

func TestSlackPost_NoAuth(t *testing.T) {
	t.Run("Slack post invalid auth", func(t *testing.T) {
		token := "INVALID"
		channel := "INVALID"
		message := buildMessage(cloudtrailEvent)

		err := slackPost(token, channel, message)
		assert.Error(t, err)
		assert.ErrorContains(t, err, "invalid_auth")
	})
}
