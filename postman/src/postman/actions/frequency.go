package actions

import (
	"postman/processor"
)

type FrequencyMsg struct {
	Domain string `json:"domain"`
	Action string `json:"action"`
	Value  string `json:"value"`
}

func Frequency(args interface{}) {
	f := args.(*FrequencyMsg)
	switch f.Action {
	case "delete":
		processor.DeleteDliverFrequency(f.Domain)
	case "update":
		processor.SaveDeliverFrequency(f.Domain, f.Value)
	}
}
