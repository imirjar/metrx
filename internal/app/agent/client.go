package agent

import (
	"fmt"
	"net/http"
	"time"
)

func sendMetric(metricType string, metric string, appURL string, value any) {
	client := &http.Client{
		Timeout: time.Second * 1,
	}

	path := fmt.Sprintf(appURL+"/update/%s/%s/%v", metricType, metric, value)
	_, err := client.Post(path, "text/plain", nil)

	if err != nil {
		fmt.Println(err)
	}
	// _, err = io.ReadAll(resp.Body)

	// resp.Body.Close()
	// if err != nil {
	// 	fmt.Println(err)
	// 	return
	// }
}
