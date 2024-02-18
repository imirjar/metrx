package agent

import "reflect"

func (a *AgentService) ReadMemValue(ms string) (float64, bool) {
	value := reflect.ValueOf(a.MemStats).FieldByName(ms)
	if value.CanFloat() {
		return value.Float(), true
	} else if value.CanUint() {
		return float64(value.Uint()), true
	} else {
		return 0, false
	}
}
