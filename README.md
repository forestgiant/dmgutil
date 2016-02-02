# dmgutil
Package providing methods for mounting, unmounting, and extracting the contents of an OSX disk image (.dmg) using the hdiutil system command.

## Install
`go get -u github.com/forestgiant/dmgutil`

## Methods

#### CopyFile(source, destination string) error
Copies the file at the source path to the provided destination.

#### CopyDirectory(source string, destination string, recursive bool) error
Copies the directory at the source path to the provided destination, with the option of recursively copying subdirectories.

#### Mount(sourcePath string) (volumePath string, err error)
Uses the `os/exec` package to issue an `hdiutil attach -nobrowse` command for the given OSX disk image (.dmg).

#### Unmount(volumePath string) error
Uses the `os/exec` package to issue an `hdiutil unmount` command for the given volume path.

#### ExtractDMG(sourcePath, destinationPath string) error
Mounts the disk image at the source path, copies the contents of the resulting volume to the destination path, then unmounts the volume.
