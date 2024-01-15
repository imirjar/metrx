package agent

import (
	"fmt"
	"net/http"
)

func sendMetric(metricType string, metric string, appURL string, value any) {
	path := fmt.Sprintf(appURL+"/update/%s/%s/%v", metricType, metric, value)
	resp, err := http.Post(path, "text/plain", nil)
	defer resp.Body.Close()
	if err != nil {
		fmt.Println(err)
	}

}
