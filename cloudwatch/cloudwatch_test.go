package cloudwatch

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"testing"
)

func TestUnmarshalASG(t *testing.T) {
	asg, err := ReadFixture("data/cloudwatch_event_ec2_launch.json")
	if err != nil {
		t.Fatalf("Err: %v", err)
	}

	var e Event
	if err := json.Unmarshal(asg, &e); err != nil {
		t.Fatalf("Err: %v", err)
	}
}

func ReadFixture(path string) ([]byte, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	return b, nil
}
