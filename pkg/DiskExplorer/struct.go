package DiskExplorer

import (
	"fmt"
	"os"
	"strings"
)

type DiskInfo struct {
	Path       string
	Name       string
	IsDir      bool
	IsExplored bool
	Size       uint64
	Children   []DiskInfo
	Mode       os.FileMode
}

func (d *DiskInfo) Prefix() string {
	if !d.IsDir {
		return "F"
	}

	if d.IsExplored {
		return "E"
	} else {
		return "N"
	}
}

func (d *DiskInfo) String() string {
	return fmt.Sprintf("%s %s %d %s", d.Prefix(), d.Mode, d.Size, d.Name)
}

func (d *DiskInfo) Render() string {
	return d.stringDepth(1)
}

func (d *DiskInfo) stringDepth(depth int) (str string) {
	str = d.String() + "\n"

	for _, child := range d.Children {
		str += strings.Repeat("\t", depth) + child.stringDepth(depth+1)
	}

	return
}
