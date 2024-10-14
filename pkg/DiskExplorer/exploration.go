package DiskExplorer

import (
	"os"
)

func Map(dir string) (directory DiskInfo) {
	var files, err = os.ReadDir(dir)
	panicOnError(err)

	directory = DiskInfo{
		Name:       dir,
		IsDir:      true,
		IsExplored: true,
		Size:       0,
	}

	for _, file := range files {
		var info, _ = file.Info()

		var child = DiskInfo{
			Name:       file.Name(),
			IsDir:      file.IsDir(),
			IsExplored: !file.IsDir(),
			Size:       uint64(info.Size()),
			Children:   []DiskInfo{},
		}

		directory.addChild(child)
		directory.Size += child.Size
	}

	return
}

func (d *DiskInfo) addChild(directory DiskInfo) {
	if d.Children == nil {
		d.Children = []DiskInfo{directory}
	} else {
		d.Children = append(d.Children, directory)
	}
}

func (d *DiskInfo) Expand() *DiskInfo {
	for i, child := range d.Children {
		if child.IsExplored {
			if child.IsDir {
				var size = child.Size
				d.Size += child.Expand().Size - size
			}
			continue
		}

		child = Map(d.Name + "/" + child.Name)
		d.Children[i] = child
		d.Size += child.Size
	}
	return d
}

func (d *DiskInfo) childrenExplored() bool {
	for i := 0; i < len(d.Children); i++ {
		if !d.Children[i].FullyExplored() {
			return false
		}
	}
	return true
}

func (d *DiskInfo) FullyExplored() bool {
	// TODO : HAVE THIS AS A PRE-CALCULATED FIELD
	return d.IsExplored && d.childrenExplored()
}
