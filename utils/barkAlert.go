package utils

import (
	"bytes"
	"fmt"
	"github.com/golang/glog"
	"github.com/skeyic/ark-robot/config"
	"net/http"
	"strings"
)

var (
	neuronServerURL = config.Config.NeuronServer.URL + "/users/" + config.Config.NeuronServer.User + "/send"
)

func SendAlert(title, content string) error {
	var (
		barkURL = "https://api.day.app/kMHL4X8KSWDWzhZyZY3hgk/%s/%s"
	)

	fmt.Println(SendRequest(http.MethodPost, fmt.Sprintf(barkURL, title, content), nil))

	return nil
}

// http://www.xiaxuanli.com:7474/users/2db982e4-9492-4202-a4c9-e615e01883f9/send -H "accept: application/json" -H "Content-Type: application/json" -d "{ \"content\": \"futu rate speaker\", \"title\": \"test\"}"
func SendAlertV2(title, content string) error {
	glog.V(4).Infof("TRY SENDING ALERT, title: %s, content: %s", title, content)
	body := bytes.NewBufferString(fmt.Sprintf("{\"content\": \"%s\", \"title\": \"%s\"}", strings.ReplaceAll(content, "\n", "    "), title))

	rCode, rBody, rError := SendRequest(http.MethodPost, neuronServerURL, body)
	if rError != nil {
		glog.Errorf("failed to send alert, rCode: %d, rBody: %v, rError: %v", rCode, rBody, rError)
		return rError
	}
	glog.V(4).Infof("SEND ALERT SUCCESSFULLY, title: %s, content: %s", title, content)

	return nil
}
