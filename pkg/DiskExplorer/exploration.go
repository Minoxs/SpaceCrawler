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
		directory.IsExplored = directory.IsExplored && child.IsExplored
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
	if d.IsExplored {
		return d
	}
	d.IsExplored = true

	for i, child := range d.Children {
		if !child.IsDir || child.IsExplored {
			continue
		}

		var size = child.Size
		if len(child.Children) == 0 {
			child = Map(d.Name + "/" + child.Name)
		} else {
			child = *child.Expand()
		}
		d.Children[i] = child

		d.IsExplored = d.IsExplored && child.IsExplored
		d.Size += child.Size - size
	}
	return d
}
