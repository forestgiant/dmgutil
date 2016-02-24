package dmgutil

import (
	"os"
	"testing"
)

func TestMount(t *testing.T) {
	var tests = []struct {
		path         string
		shouldAssert bool
	}{
		{"./TestFiles/test_0.0.1.txt", false},
		{"./TestFiles/dmgutil test_0.0.1.txt", false},
		{"./TestFiles/dmgutil test_0.0.1.dmg", true},
	}

	for index, test := range tests {
		volumePath, err := Mount(test.path)
		if err != nil && test.shouldAssert == true {
			t.Errorf("Test %d failed but should have passed.", index)
		} else if err == nil && test.shouldAssert == false {
			t.Errorf("Test %d passed but should have failed.", index)
		} else if err == nil && test.shouldAssert == true {
			if _, err := os.Stat(volumePath); err != nil {
				t.Errorf("Test %d reported a successful mount, but the resulting volume path does not exist.", index)
			}

			if err := Unmount(volumePath); err != nil {
				t.Errorf("Test %d successfully mounted the volume, but was unable to unmount the volume.", index)
			}
		}
	}
}

func TestUnmount(t *testing.T) {
	volumePath, err := Mount("./TestFiles/dmgutil test_0.0.1.dmg")
	if err != nil {
		t.Error("Unable to run TestUnmount.  Mount of test image failed.")
		return
	}

	if _, err := os.Stat(volumePath); err != nil {
		t.Error("Unable to run TestUnmount.  Mount of test image reported no errors, however the volume is not present.")
		return
	}

	if err := Unmount(volumePath); err != nil {
		t.Error(err.Error())
		return
	}

	if _, err := os.Stat(volumePath); err == nil {
		t.Error("TestUnmount failed.  Unmount reported no errors, however the volume is still listed.")
		return
	}
}

func TestExtractDMG(t *testing.T) {
	var destinationDirectory = "./TestFiles/Applications"
	if err := ExtractDMG("./TestFiles/dmgutil test_0.0.1.dmg", destinationDirectory); err != nil {
		t.Error("TestExtractDMG failed. " + err.Error())
		return
	}

	if _, err := os.Stat(destinationDirectory); err != nil {
		t.Error("ExtractDMG indicated no errors, however the expected destination directory was not created.")
	}

	if err := os.RemoveAll(destinationDirectory); err != nil {
		t.Error("Unable to clean up after running TestExtractDMG.")
	}
}
