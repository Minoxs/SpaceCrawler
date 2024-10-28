package DiskExplorer

import (
	"fmt"
	"os"
	"strings"
)

// DiskInfo struct contains the information a node in a file tree
type DiskInfo struct {
	Path     string
	Name     string
	IsDir    bool
	Children []DiskInfo
	Mode     os.FileMode

	size uint64
}

// Prefix returns a string prefix depending on the node type
func (d *DiskInfo) Prefix() string {
	if !d.IsDir {
		return "F"
	}

	if d.Explored() {
		return "E"
	} else {
		return "N"
	}
}

// String returns a string representation of the node
func (d *DiskInfo) String() string {
	return fmt.Sprintf("%s %s %d %s", d.Prefix(), d.Mode, d.Size(), d.Name)
}

// Render will return a string with the file tree representation
func (d *DiskInfo) Render() string {
	return d.stringDepth(1)
}

// stringDepth will recursively build the file tree representation
func (d *DiskInfo) stringDepth(depth int) (str string) {
	str = d.String() + "\n"

	for _, child := range d.Children {
		str += strings.Repeat("\t", depth) + child.stringDepth(depth+1)
	}

	return
}

// Explored will recursively check if the node is explored
func (d *DiskInfo) Explored() bool {
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
