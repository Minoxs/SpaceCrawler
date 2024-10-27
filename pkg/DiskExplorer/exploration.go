package DiskExplorer

import (
	"log"
	"os"
	"path/filepath"
)

// Map returns the disk info for a particular directory
// This will list out all the contents of the given directory and no more
// Call DiskInfo.Expand to map out the lower layers of the tree
func Map(dir string) (directory DiskInfo) {
	var path, _ = filepath.Abs(dir)

	directory = DiskInfo{
		Path:       path,
		Name:       dir,
		IsDir:      true,
		IsExplored: true,
		Size:       0,
	}

	var files, err = os.ReadDir(dir)
	if err != nil {
		log.Println(err)
		return
	}

	for _, file := range files {
		var info, _ = file.Info()

		var child = DiskInfo{
			Path:       filepath.Join(path, file.Name()),
			Name:       file.Name(),
			IsDir:      file.IsDir(),
			IsExplored: !file.IsDir(),
			Size:       uint64(info.Size()),
			Children:   []DiskInfo{},
			Mode:       info.Mode(),
		}

		directory.addChild(child)
		directory.Size += child.Size
		directory.IsExplored = directory.IsExplored && child.IsExplored
	}

	return
}

// addChild adds appends a new child to the end of the tree
func (d *DiskInfo) addChild(directory DiskInfo) {
	if d.Children == nil {
		d.Children = []DiskInfo{directory}
	} else {
		d.Children = append(d.Children, directory)
	}
}

// Expand will recursively expand the tree further down
// Only expands 1 layer at a time
func (d *DiskInfo) Expand() *DiskInfo {
	if d.IsExplored {
		return d
	}
	d.IsExplored = true

	if len(d.Children) == 0 {
		*d = Map(d.Path)
		return d
	}

	for i, child := range d.Children {
		if !child.IsDir || child.IsExplored {
			continue
		}

		if len(child.Children) == 0 {
			d.Children[i] = Map(child.Path)
		} else {
			d.Children[i].Expand()
		}

		d.Size += d.Children[i].Size - child.Size
		d.IsExplored = d.IsExplored && d.Children[i].IsExplored
	}
	return d
}
