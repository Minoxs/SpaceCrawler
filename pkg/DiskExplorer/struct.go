package DiskExplorer

import (
	"fmt"
	"os"
	"strings"
)

type DiskInfo struct {
	Name       string
	IsDir      bool
	IsExplored bool
	Size       uint64
	Children   []DiskInfo
	Mode       os.FileMode
}

func (d *DiskInfo) Prefix() string {
	var explored = " N"
	if d.IsExplored {
		explored = " E"
	}

	if d.IsDir {
		return "[DIR]" + explored
	} else {
		return "[FLE]" + explored
	}
}

func (d *DiskInfo) stringDepth(depth int) (str string) {
	str = fmt.Sprintf("%s %s %d %s\n", d.Prefix(), d.Mode, d.Size, d.Name)

	for _, child := range d.Children {
		str += strings.Repeat("\t", depth) + child.stringDepth(depth+1)
	}

	return
}

func (d *DiskInfo) String() (str string) {
	return d.stringDepth(1)
}
