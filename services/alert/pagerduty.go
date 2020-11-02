package alert

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/rudderlabs/rudder-server/utils/logger"
)

var pagerDutyEndPoint = "https://events.pagerduty.com/v2/enqueue"
var pkgLogger logger.LoggerI

func init() {
	pkgLogger = logger.NewLogger().Child("services").Child("alert")
}

func (ops *PagerDuty) Alert(message string) {

	payload := map[string]interface{}{
		"summary":  message,
		"severity": "critical",
		"source":   ops.instanceName,
	}

	event := map[string]interface{}{
		"payload":      payload,
		"event_action": "trigger",
		"routing_key":  ops.routingKey,
	}

	eventJSON, _ := json.Marshal(event)
	client := &http.Client{}
	resp, err := client.Post(pagerDutyEndPoint, "application/json", bytes.NewBuffer(eventJSON))
	// Not handling errors when sending alert to victorops
	if err != nil {
		pkgLogger.Errorf("Alert: Failed to alert service: %s", err.Error())
		return
	}

	if resp.StatusCode != 200 && resp.StatusCode != 202 {
		pkgLogger.Errorf("Alert: Got error response %d", resp.StatusCode)
	}

	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	pkgLogger.Infof("Alert: Successful %s", string(body))
}

type PagerDuty struct {
	instanceName string
	routingKey   string
}
