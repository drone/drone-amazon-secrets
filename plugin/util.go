// Copyright 2018 Drone.IO Inc
// Use of this software is governed by the Business Source License
// that can be found in the LICENSE file.

package plugin

import (
	"fmt"
	"strings"
)

// helper function extracts the branch filters from the
// secret payload in key value format.
func extractBranches(params map[string]interface{}) []string {
	for key, value := range params {
		if strings.EqualFold(key, "X-Drone-Branches") {
			return parseCommaSeparated(value)
		}
	}
	return nil
}

// helper function extracts the repository filters from the
// secret payload in key value format.
func extractRepos(params map[string]interface{}) []string {
	for key, value := range params {
		if strings.EqualFold(key, "X-Drone-Repos") {
			return parseCommaSeparated(value)
		}
	}
	return nil
}

// helper function extracts the event filters from the
// secret payload in key value format.
func extractEvents(params map[string]interface{}) []string {
	for key, value := range params {
		if strings.EqualFold(key, "X-Drone-Events") {
			return parseCommaSeparated(value)
		}
	}
	return nil
}

func parseCommaSeparated(s interface{}) []string {
	str := fmt.Sprintf("%v", s)
	parts := strings.Split(str, ",")
	if len(parts) == 1 && parts[0] == "" {
		return nil
	}
	return parts
}
