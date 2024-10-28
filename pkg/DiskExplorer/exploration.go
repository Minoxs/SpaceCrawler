package DiskExplorer

import (
	"log"
	"os"
	"path/filepath"
)

// Map returns the disk info for a particular directory
// This will list out all the contents of the given directory and no more
// Call DiskInfo.Deepen to map out the lower layers of the tree
func Map(path string) (directory DiskInfo) {
	var abs, _ = filepath.Abs(path)
	var info, _ = os.Stat(abs)

	directory = DiskInfo{
		Path:     abs,
		Name:     filepath.Base(abs),
		IsDir:    true,
		Children: []DiskInfo{},
		Mode:     info.Mode(),

		isExplored: true,
		size:       0,
	}

	directory.explore()
	return
}

// explore will iterate over the directory and calculate its size
func (d *DiskInfo) explore() {
	d.isExplored = true

	var files, err = os.ReadDir(d.Path)
	if err != nil {
		log.Println(err)
		return
	}

	for _, file := range files {
		var info, _ = file.Info()

		var child = DiskInfo{
			Path:     filepath.Join(d.Path, file.Name()),
			Name:     file.Name(),
			IsDir:    file.IsDir(),
			Children: []DiskInfo{},
			Mode:     info.Mode(),

			isExplored: !file.IsDir(),
			size:       uint64(info.Size()),
		}

		d.addChild(child)
	}
}

// addChild adds appends a new child to the end of the tree
func (d *DiskInfo) addChild(child DiskInfo) {
	if d.Children == nil {
		d.Children = []DiskInfo{child}
	} else {
		d.Children = append(d.Children, child)
	}
}

// Expand will explore a given node and populate its children
// Returns whether it changed
func (d *DiskInfo) Expand() bool {
	if len(d.Children) > 0 || d.IsExplored() {
		return false
	}
	d.explore()
	return true
}

// Deepen will deepen the exploration by 1 layer
func (d *DiskInfo) Deepen() {
	if !d.Expand() {
		for i := 0; i < len(d.Children); i++ {
			d.Children[i].Deepen()
		}
	}
}
