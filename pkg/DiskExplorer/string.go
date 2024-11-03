package DiskExplorer

import (
	"fmt"
	"strings"
)

// Prefix returns a string prefix depending on the node type
func (d *DiskInfo) Prefix() string {
	if d.denied {
		return "D"
	}

	if !d.IsDir {
		return "F"
	}

	if d.Explored() {
		return "E"
	} else {
		return "N"
	}
}

// FullPrefix returns the Prefix including Depth and Breadth information
func (d *DiskInfo) FullPrefix() (s string) {
	if d.IsDir {
		return d.Prefix() + fmt.Sprintf(" D%d B%d", d.Depth(), d.Breadth())
	} else {
		return d.Prefix()
	}
}

// String returns a string representation of the node
func (d *DiskInfo) String() string {
	return fmt.Sprintf("%s %s %s %s", d.FullPrefix(), d.Mode, d.HumanSize(), d.Name)
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

// HumanSize returns the size in a human-readable format
func (d *DiskInfo) HumanSize() string {
	var s = float64(d.Size())

	var (
		idx  = 0
		unit = []string{"B", "Kb", "Mb", "Gb", "Tb", "Pb"}
	)

	for ; s > 1024 && idx < len(unit); idx++ {
		s /= 1024
	}

	return fmt.Sprintf("%.2f%s", s, unit[idx])
}
