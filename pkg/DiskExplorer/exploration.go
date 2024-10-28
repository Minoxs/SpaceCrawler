package DiskExplorer

import (
	"log"
	"os"
	"path/filepath"
	"sync/atomic"
)

// Map returns the disk info for a particular directory
// This will list out all the contents of the given directory and no more
// Call DiskInfo.Deepen to map out the lower layers of the tree
func Map(path string) (directory DiskInfo) {
	var abs, _ = filepath.Abs(path)
	var info, _ = os.Stat(abs)

	directory = DiskInfo{
		Path:       abs,
		Name:       filepath.Base(abs),
		IsDir:      true,
		IsExplored: true,
		Size:       0,
		Children:   []DiskInfo{},
		Mode:       info.Mode(),
	}

	directory.explore()
	return
}

// explore will iterate over the directory and calculate its size
func (d *DiskInfo) explore() {
	d.IsExplored = true

	var files, err = os.ReadDir(d.Path)
	if err != nil {
		log.Println(err)
		return
	}

	for _, file := range files {
		var info, _ = file.Info()

		var child = DiskInfo{
			Path:       filepath.Join(d.Path, file.Name()),
			Name:       file.Name(),
			IsDir:      file.IsDir(),
			IsExplored: !file.IsDir(),
			Size:       uint64(info.Size()),
			Children:   []DiskInfo{},
			Mode:       info.Mode(),
		}

		d.addChild(child)
		d.Size += child.Size
		d.IsExplored = d.IsExplored && child.IsExplored
	}
}

// addChild adds appends a new child to the end of the tree
func (d *DiskInfo) addChild(directory DiskInfo) {
	if d.Children == nil {
		d.Children = []DiskInfo{directory}
	} else {
		d.Children = append(d.Children, directory)
	}
}

// addSize will atomically add the size to the size of the node
func (d *DiskInfo) addSize(size uint64) {
	atomic.AddUint64(&d.Size, size)
}

// GetSize will atomically get the size of the node
func (d *DiskInfo) GetSize() uint64 {
	return atomic.LoadUint64(&d.Size)
}

// Expand will explore a given node and populate its children
// Returns whether it changed
func (d *DiskInfo) Expand() bool {
	if d.IsExplored || len(d.Children) > 0 {
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
