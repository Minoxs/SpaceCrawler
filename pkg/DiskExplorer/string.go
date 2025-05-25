package DiskExplorer

import (
	"fmt"
	"io"
	"log"
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

// Render renders the disk info into the given writer
func (d *DiskInfo) Render(writer io.Writer) {
	d.renderDepth(writer, 0)
}

// renderDepth will recursively build the file tree representation
func (d *DiskInfo) renderDepth(writer io.Writer, depth int) {
	var rep = strings.Repeat("\t", depth) + d.String() + "\n"
	var _, err = writer.Write([]byte(rep))
	if err != nil {
		log.Fatal(err)
	}

	for _, child := range d.Children {
		child.renderDepth(writer, depth+1)
	}
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
