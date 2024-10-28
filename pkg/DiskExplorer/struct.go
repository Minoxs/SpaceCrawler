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

	size       uint64
	isExplored bool
}

// Prefix returns a string prefix depending on the node type
func (d *DiskInfo) Prefix() string {
	if !d.IsDir {
		return "F"
	}

	if d.IsExplored() {
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

// IsExplored will recursively check if the node is explored
func (d *DiskInfo) IsExplored() bool {
	if !d.isExplored {
		return false
	}

	for i := 0; i < len(d.Children); i++ {
		if !d.Children[i].IsExplored() {
			return false
		}
	}

	return true
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
