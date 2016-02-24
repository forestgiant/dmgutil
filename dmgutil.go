//Package dmgutil is a simple package that provides methods for mounting,
//unmounting, and extracting the contents of an OSX disk image (.dmg) using the
//hdiutil system command.
package dmgutil

import (
	"errors"
	"os/exec"
	"strings"

	"github.com/forestgiant/fsutil"
)

//Mount uses the `os/exec` package to issue an `hdiutil attach -nobrowse` command for the given OSX disk image (.dmg).
func Mount(sourcePath string) (volumePath string, err error) {
	var command = exec.Command("hdiutil", "attach", "-nobrowse", sourcePath)
	outputBytes, err := command.Output()
	if err != nil {
		return "", err
	}

	volumes := strings.Split(string(outputBytes), "\t")
	for index, volumeName := range volumes {
		volumes[index] = strings.TrimSpace(volumeName)
	}

	volumePath = volumes[len(volumes)-1]
	return volumePath, nil
}

//Unmount uses the `os/exec` package to issue an `hdiutil unmount` command for the given volume path.
func Unmount(volumePath string) error {
	var command = exec.Command("hdiutil", "unmount", volumePath)
	if err := command.Start(); err != nil {
		return err
	}

	if err := command.Wait(); err != nil {
		return err
	}

	return nil
}

//ExtractDMG mounts the disk image at the source path, copies the contents of the resulting volume to the destination path, then unmounts the volume.
func ExtractDMG(sourcePath, destinationPath string) error {
	volumePath, err := Mount(sourcePath)
	if err != nil {
		return errors.New("Failed to mount the volume. " + err.Error())
	}

	if err := fsutil.CopyDirectory(volumePath, destinationPath, true); err != nil {
		return errors.New("Failed to copy contents of the volume. " + err.Error())
	}

	if err := Unmount(volumePath); err != nil {
		return errors.New("Failed to unmount the volume. " + err.Error())
	}

	return nil
}
