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

// lenLoop10 is a helper function to quickly calculate integer string length
// https://stackoverflow.com/a/68124773
func lenLoop10(i int) int {
	if i >= 1e18 {
		return 19
	}
	x, count := 10, 1
	for x <= i {
		x *= 10
		count++
	}
	return count
}

// FullPrefix returns the Prefix including Depth and Breadth information
func (d *DiskInfo) FullPrefix() (s string) {
	var (
		dWidth = 0
		bWidth = 0
		tWidth = 0
	)
	if d.Parent != nil {
		dWidth = lenLoop10(d.Parent.CachedDepth())
		bWidth = lenLoop10(d.Parent.MaxChildBreadth())
		tWidth = lenLoop10(d.Parent.MaxChildTotalFiles())
	}

	if d.IsDir {
		return d.Prefix() + fmt.Sprintf(
			" D%0*d B%0*d T%0*d",
			dWidth,
			d.CachedDepth(),
			bWidth,
			d.CachedBreath(),
			tWidth,
			d.CachedTotalFiles(),
		)
	} else {
		return d.Prefix() + strings.Repeat(" ", 3+3+dWidth+bWidth+tWidth)
	}
}

// String returns a string representation of the node
func (d *DiskInfo) String() string {
	// Ensure cache is up to date
	d.UpdateCache()
	return fmt.Sprintf("%s %s %s %5.2f%% %s", d.FullPrefix(), d.Mode, d.HumanSize(), d.PercentSize(), d.Name)
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
// Call Size before this to ensure the value is up to date
func (d *DiskInfo) HumanSize() string {
	var s = float64(d.CachedSize())

	var (
		idx  = 0
		unit = []string{"B", "Kb", "Mb", "Gb", "Tb", "Pb"}
	)

	for ; s > 1024 && idx < len(unit)-1; idx++ {
		s /= 1024
	}

	var pad = ""
	if idx == 0 {
		pad = " "
	}

	return fmt.Sprintf("%6.2f%s", s, unit[idx]) + pad
}

// PercentSize calculates the percentage of the current node's size compared to its parent
// Call Size before this to ensure the value is up to date
func (d *DiskInfo) PercentSize() float64 {
	if d.Parent == nil {
		return 100
	}

	return 100 * float64(d.CachedSize()) / float64(d.Parent.CachedSize())
}
