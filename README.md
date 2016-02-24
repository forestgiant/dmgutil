# dmgutil
Package providing methods for mounting, unmounting, and extracting the contents of an OSX disk image (.dmg) using the hdiutil system command.

## Install
`go get -u github.com/forestgiant/dmgutil`

## Methods

#### Mount
Uses the `os/exec` package to issue an `hdiutil attach -nobrowse` command for the given OSX disk image (.dmg).

#### Unmount
Uses the `os/exec` package to issue an `hdiutil unmount` command for the given volume path.

#### ExtractDMG
Mounts the disk image at the source path, copies the contents of the resulting volume to the destination path, then unmounts the volume.
