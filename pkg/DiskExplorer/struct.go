package DiskExplorer

import (
	"os"
)

// DiskInfo struct contains the information a node in a file tree
type DiskInfo struct {
	Path     string
	Name     string
	IsDir    bool
	Children []DiskInfo
	Mode     os.FileMode

	size   uint64
	denied bool
}

// Explored will recursively check if the node is explored
func (d *DiskInfo) Explored() bool {
	// Denied files are explored by default
	if d.denied {
		return true
	}

	// Files are explored by default
	if !d.IsDir {
		return true
	}

	// If the children array is nil then this node has not been explored
	if d.Children == nil {
		return false
	}

	// Check if every child is also explored
	for i := 0; i < len(d.Children); i++ {
		if !d.Children[i].Explored() {
			return false
		}
	}

	return true
}

// Expanded will check if the node has been expanded already
func (d *DiskInfo) Expanded() bool {
	return len(d.Children) > 0 || d.Explored()
}

// Size will recursively calculate the size of the node for directories
// File size is cached
func (d *DiskInfo) Size() uint64 {
	if d.IsDir {
		var size = uint64(0)
		for i := 0; i < len(d.Children); i++ {
			size += d.Children[i].Size()
		}
		return size
	} else {
		return d.size
	}
}

// Depth calculates the maximum depth of a given node
func (d *DiskInfo) Depth() int {
	if !d.IsDir || !d.Expanded() {
		return 0
	}

	var deep = 0
	for i := 0; i < len(d.Children); i++ {
		deep = max(deep, d.Children[i].Depth())
	}
	return 1 + deep
}

// Breadth calculates the breadth of a given node
func (d *DiskInfo) Breadth() int {
	return len(d.Children)
}

// Denied returns if access to the given node has been denied
func (d *DiskInfo) Denied() bool {
	return d.denied
}
