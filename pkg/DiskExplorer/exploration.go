package DiskExplorer

import (
	"os"
	"path/filepath"
)

// Map returns the disk info for a particular directory
// This will list out all the contents of the given directory and no more
// Call DiskInfo.Deepen to map out the lower layers of the tree
// Use DiskInfo.Expand to expand specific nodes
func Map(path string) (d DiskInfo) {
	// Set Info name
	var abs, _ = filepath.Abs(path)
	d = DiskInfo{
		Path: abs,
		Name: filepath.Base(abs),
	}

	// Grab metadata from disk
	var info, err = os.Stat(abs)
	if err != nil {
		d.denied = true
		return
	}

	const InvalidFile = os.ModeSymlink | os.ModeNamedPipe | os.ModeSocket | os.ModeDevice | os.ModeCharDevice | os.ModeIrregular
	d.IsDir = info.IsDir()
	d.Mode = info.Mode()
	d.denied = d.Mode.Type()&InvalidFile != 0
	if d.denied {
		return
	}

	d.explore()
	return
}

// Expand will explore a given node and populate its children
// Returns whether it changed
func (d *DiskInfo) Expand() bool {
	if d.Expanded() {
		return false
	}

	d.explore()
	return true
}

// explore will iterate over the directory and list out subfolders and files
func (d *DiskInfo) explore() {
	var files, err = os.ReadDir(d.Path)
	if err != nil {
		d.denied = true
		return
	}

	var tmp = make([]DiskInfo, len(files))
	for i, file := range files {
		var info, _ = file.Info()

		var child = DiskInfo{
			Path:     filepath.Join(d.Path, file.Name()),
			Name:     file.Name(),
			IsDir:    file.IsDir(),
			Children: nil,
			Mode:     info.Mode(),

			size: uint64(info.Size()),
		}

		tmp[i] = child
	}

	d.Children = tmp
}
