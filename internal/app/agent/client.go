package agent

import (
	"fmt"
	"net/http"
)

func sendMetric(metricType string, metric string, appURL string, value any) {
	path := fmt.Sprintf(appURL+"/update/%s/%s/%v", metricType, metric, value)
	_, err := http.Post(path, "text/plain", nil)
	if err != nil {
		fmt.Println(err)
	}

}
