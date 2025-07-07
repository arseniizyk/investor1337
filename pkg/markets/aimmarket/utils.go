package aimmarket

import "fmt"

func (am *aimmarket) preparePayload(name string) map[string]any {
	payload := map[string]any{
		"operationName": "ApiBotsInventoryCountAndMinPrice",
		"query":         string(am.query),
		"variables": map[string]any{
			"currency": "USD",
			"where": map[string]any{
				"marketHashName": map[string]string{
					"_text": fmt.Sprintf("\"%s\"", name),
				},
			},
		},
	}

	return payload
}
