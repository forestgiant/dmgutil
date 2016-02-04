//Package dmgutil is a simple package that provides methods for mounting,
//unmounting, and extracting the contents of an OSX disk image (.dmg) using the
//hdiutil system command.
package dmgutil

import (
	"errors"
	"io"
	"os"
	"os/exec"
	"strings"
)

//CopyFile copies the file at the source path to the provided destination.
func CopyFile(source, destination string) error {
	//Validate the source and destination paths
	if len(source) == 0 {
		return errors.New("You must provide a source file path.")
	}

	if len(destination) == 0 {
		return errors.New("You must provide a destination file path.")
	}

	//Verify the source path refers to a regular file
	sourceFileInfo, err := os.Lstat(source)
	if err != nil {
		return err
	}

	//Handle regular files differently than symbolic links and other non-regular files.
	if sourceFileInfo.Mode().IsRegular() {
		//open the source file
		sourceFile, err := os.Open(source)
		if err != nil {
			return err
		}
		defer sourceFile.Close()

		//create the destinatin file
		destinationFile, err := os.Create(destination)
		if err != nil {
			return err
		}
		defer destinationFile.Close()

		//copy the source file contents to the destination file
		if _, err = io.Copy(destinationFile, sourceFile); err != nil {
			return err
		}

		//replicate the source file mode for the destination file
		if err := os.Chmod(destination, sourceFileInfo.Mode()); err != nil {
			return err
		}
	} else if sourceFileInfo.Mode()&os.ModeSymlink != 0 {
		linkDestinaton, err := os.Readlink(source)
		if err != nil {
			return errors.New("Unable to read symlink. " + err.Error())
		}

		if err := os.Symlink(linkDestinaton, destination); err != nil {
			return errors.New("Unable to replicate symlink. " + err.Error())
		}
	} else {
		return errors.New("Unable to use io.Copy on file with mode " + string(sourceFileInfo.Mode()))
	}

	return nil
}

//CopyDirectory copies the directory at the source path to the provided destination, with the option of recursively copying subdirectories.
func CopyDirectory(source string, destination string, recursive bool) error {
	if len(source) == 0 || len(destination) == 0 {
		return errors.New("File paths must not be empty.")
	}

	//get properties of the source directory
	sourceInfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	//create the destination directory
	err = os.MkdirAll(destination, sourceInfo.Mode())
	if err != nil {
		return err
	}

	sourceDirectory, err := os.Open(source)
	if err != nil {
		return err
	}
	defer sourceDirectory.Close()

	objects, err := sourceDirectory.Readdir(-1)
	if err != nil {
		return err
	}

	for _, object := range objects {
		if object.Name() == ".Trashes" {
			continue
		}

		sourceObjectName := source + "/" + object.Name()
		destObjectName := destination + "/" + object.Name()

		if object.IsDir() {
			//create sub-directories
			err = CopyDirectory(sourceObjectName, destObjectName, true)
			if err != nil {
				return err
			}
		} else {
			err = CopyFile(sourceObjectName, destObjectName)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

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

	if err := CopyDirectory(volumePath, destinationPath, true); err != nil {
		return errors.New("Failed to copy contents of the volume. " + err.Error())
	}

	if err := Unmount(volumePath); err != nil {
		return errors.New("Failed to unmount the volume. " + err.Error())
	}

	return nil
}
