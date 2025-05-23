package DiskExplorer

import (
	"context"
	"os"
	"path/filepath"
	"sync"
)

// Map returns the disk info for a particular directory
// This will list out all the contents of the given directory and no more
// Call DiskInfo.Deepen to map out the lower layers of the tree
// Use DiskInfo.Expand to expand specific nodes
func Map(path string) (d DiskInfo) {
	// Set Info name
	var abs, _ = filepath.Abs(path)
	d = DiskInfo{
		Path: abs,
		Name: filepath.Base(abs),
	}

	// Grab metadata from disk
	var info, err = os.Stat(abs)
	if err != nil {
		d.denied = true
		return
	}

	const InvalidFile = os.ModeSymlink | os.ModeNamedPipe | os.ModeSocket | os.ModeDevice | os.ModeCharDevice | os.ModeIrregular
	d.IsDir = info.IsDir()
	d.Mode = info.Mode()
	d.denied = d.Mode.Type()&InvalidFile != 0
	if d.denied {
		return
	}

	d.explore()
	return
}

// explore will iterate over the directory and list out subfolders and files
func (d *DiskInfo) explore() {
	d.Children = []DiskInfo{}

	var files, err = os.ReadDir(d.Path)
	if err != nil {
		d.denied = true
		return
	}

	for _, file := range files {
		var info, _ = file.Info()

		var child = DiskInfo{
			Path:     filepath.Join(d.Path, file.Name()),
			Name:     file.Name(),
			IsDir:    file.IsDir(),
			Children: nil,
			Mode:     info.Mode(),

			size: uint64(info.Size()),
		}

		d.addChild(child)
	}
}

// addChild adds appends a new child to the end of the tree
func (d *DiskInfo) addChild(child DiskInfo) {
	d.Children = append(d.Children, child)
}

// Expand will explore a given node and populate its children
// Returns whether it changed
func (d *DiskInfo) Expand() bool {
	if d.Expanded() {
		return false
	}

	d.explore()
	return true
}

// Exhaust will expand every single node without stopping
// A context can be passed into this function to cancel out of it early
func (d *DiskInfo) Exhaust(ctx context.Context) {
	if ctx.Err() != nil {
		return
	}

	d.Expand()

	var wg = sync.WaitGroup{}
	for i := 0; i < len(d.Children); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			d.Children[i].Exhaust(ctx)
		}(i)
	}
	wg.Wait()
}

// Deepen will deepen the exploration by 1 layer
func (d *DiskInfo) Deepen() {
	if !d.Expand() {
		var wg = sync.WaitGroup{}
		for i := 0; i < len(d.Children); i++ {
			wg.Add(1)
			go func(i int) {
				defer wg.Done()
				d.Children[i].Deepen()
			}(i)
		}
		wg.Wait()
	}
}
