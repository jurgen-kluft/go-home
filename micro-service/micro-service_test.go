package microservice

import (
	"testing"
)

func TestMatchTopic(t *testing.T) {
	// Matches
	if matchTopic("state/light/automation/", "state/light/automation/") == false {
		t.Fail()
	}
	if matchTopic("state/light/*/", "state/light/automation/") == false {
		t.Fail()
	}

	// Force failures
	if matchTopic("*/switch/*/", "state/light/automation/") == true {
		t.Fail()
	}
	if matchTopic("*/*/switch/", "state/light/automation/") == true {
		t.Fail()
	}
	if matchTopic("state/light/automation", "state/light/automation/") == true {
		t.Fail()
	}
}
