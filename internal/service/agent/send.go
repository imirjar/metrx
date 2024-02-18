package agent

func (a *AgentService) SendMetrix(url string) {
	var counter = 0
	gauges, err := a.Storage.ReadAll("gauge")
	if err != nil {
		panic(err)
	}
	counters, err := a.Storage.ReadAll("counter")
	if err != nil {
		panic(err)
	}

	for _, g := range gauges {
		a.MetricsClient.SendJSON(&g)
		counter++
	}

	for _, c := range counters {
		a.MetricsClient.SendJSON(&c)
		counter++
	}
}
