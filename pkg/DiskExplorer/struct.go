package DiskExplorer

import (
	"fmt"
	"strings"
)

type DiskInfo struct {
	Name       string
	IsDir      bool
	IsExplored bool
	Size       uint64
	Children   []DiskInfo

	expanded bool
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
	str = fmt.Sprintf("%s %d %s\n", d.Prefix(), d.Size, d.Name)

	for _, child := range d.Children {
		str += strings.Repeat("\t", depth) + child.stringDepth(depth+1)
	}

	return
}

func (d *DiskInfo) String() (str string) {
	return d.stringDepth(1)
}
