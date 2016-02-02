package dmgutil

import (
	"os"
	"testing"
)

func TestCopyFile(t *testing.T) {
	var badSourceFile = "./TestFiles/FGFailedSourceFileLocation.fg"
	var badDestinationFile = "./TestFiles/NonexistantDirectory/FGFailedDestFileLocation.fg"
	var goodSourceFile = "./TestFiles/DiskImages/dmgutil test_0.0.1.dmg"
	var goodDestinationFile = "./TestFiles/DiskImages/dmgutil test_0.0.1_copy.dmg"
	var symlinkSourceFile = "./TestFiles/DiskImages/TestSubdirectory/TextFiles"
	var symlinkDestinationFile = "./TestFiles/DiskImages/TestSubdirectory/TextFilesCopy.dmg"

	var tests = []struct {
		source       string
		destination  string
		shouldAssert bool
	}{
		{"", goodDestinationFile, false},
		{goodSourceFile, "", false},
		{badSourceFile, goodDestinationFile, false},
		{goodSourceFile, badDestinationFile, false},
		{goodSourceFile, goodDestinationFile, true},
		{symlinkSourceFile, symlinkDestinationFile, true},
	}

	for index, test := range tests {
		err := CopyFile(test.source, test.destination)
		if err != nil && test.shouldAssert == true {
			t.Errorf("Test %d failed but should have passed. "+err.Error(), index)
		} else if err == nil && test.shouldAssert == false {
			t.Errorf("Test %d passed but should have failed.", index)
		} else if err == nil && test.shouldAssert == true {
			if _, err := os.Stat(test.source); err != nil {
				t.Error("CopyFile did not throw an error, but the test file was not created.")
			}

			if err := os.Remove(test.destination); err != nil {
				t.Error("Unable to clean up after running CopyFile tests. " + err.Error())
			}
		}
	}
}

func TestCopyDestinationFile(t *testing.T) {
	var badSourceDirectory = "./TestFiles/BadSourceFolder"
	var goodSourceDirectory = "./TestFiles/DiskImages"
	var goodDestinationDirectory = "./TestFiles/DiskImagesCopy"

	var tests = []struct {
		source       string
		destination  string
		shouldAssert bool
	}{
		{"", goodDestinationDirectory, false},
		{goodSourceDirectory, "", false},
		{badSourceDirectory, goodDestinationDirectory, false},
		{goodSourceDirectory, goodDestinationDirectory, true},
	}

	for index, test := range tests {
		err := CopyDirectory(test.source, test.destination, true)
		if err != nil && test.shouldAssert == true {
			t.Errorf("Test %d failed but should have passed.", index)
		} else if err == nil && test.shouldAssert == false {
			t.Errorf("Test %d passed but should have failed.", index)
		}
	}

	if _, err := os.Stat(goodDestinationDirectory); err != nil {
		t.Error("CopyDirectory tests passed, but the test directory was not created.")
	}

	if err := os.RemoveAll(goodDestinationDirectory); err != nil {
		t.Error("Unable to clean up after running CopyDirectory tests. " + err.Error())
	}
}

func TestMount(t *testing.T) {
	var tests = []struct {
		path         string
		shouldAssert bool
	}{
		{"./TestFiles/TextFiles/dmgutil test_0.0.1.txt", false},
		{"./TestFiles/DiskImages/dmgutil test_0.0.1.dmg", true},
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
	volumePath, err := Mount("./TestFiles/DiskImages/dmgutil test_0.0.1.dmg")
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
	if err := ExtractDMG("./TestFiles/DiskImages/dmgutil test_0.0.1.dmg", destinationDirectory); err != nil {
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
