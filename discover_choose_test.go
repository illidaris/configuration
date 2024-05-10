package configuration

import (
	"math/rand"
	"testing"
	"time"

	"github.com/nacos-group/nacos-sdk-go/v2/model"
)

// TestDefaultChoose tests the defaultChoose function.
func TestDefaultChoose(t *testing.T) {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator

	// Create a list of instances for testing
	instances := []model.Instance{
		{Weight: 1, InstanceId: "1"},
		{Weight: 2, InstanceId: "2"},
		{Weight: 3, InstanceId: "3"},
	}

	// Call the defaultChoose function
	selectedInstance := defaultChoose(instances...)

	// Check if a valid instance is returned
	if selectedInstance.InstanceId == "" {
		t.Errorf("defaultChoose returned nil, expected a valid model.Instance")
	}
}

// TestChooser tests the Chooser struct and its methods.
func TestChooser(t *testing.T) {
	rand.Seed(time.Now().UnixNano()) // Seed the random number generator

	// Create a list of instances for testing
	instances := []model.Instance{
		{Weight: 1, InstanceId: "1"},
		{Weight: 2, InstanceId: "2"},
		{Weight: 3, InstanceId: "3"},
	}

	// Initialize a new chooser with the instances
	chooser := newChooser(instances)

	// Call the Pick method
	selectedInstance := chooser.Pick()

	// Check if a valid instance is returned
	if selectedInstance.InstanceId == "" {
		t.Errorf("defaultChoose returned nil, expected a valid model.Instance")
	}
}
