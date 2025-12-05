package DiskExplorer

// UpdateCache will update cached values for faster access
// Depth and TotalFiles will only be updated if this is called on the root node
func (d *DiskInfo) UpdateCache() {
	d.Size()
	if d.Parent == nil {
		d.Depth()
		d.TotalFiles()
	}
}

// CachedSize returns the value from the last call of Size
func (d *DiskInfo) CachedSize() uint64 {
	return d.size
}

// CachedDepth returns the value from the last call of Depth
func (d *DiskInfo) CachedDepth() int {
	return d.depth
}

// CachedBreath returns the value from the last call of Breath
func (d *DiskInfo) CachedBreath() int {
	return d.Breadth()
}

// CachedTotalFiles returns the value from the last call of TotalFiles
func (d *DiskInfo) CachedTotalFiles() int {
	return d.total
}

// MaxChildBreadth calculates the maximum bread of a child node
func (d *DiskInfo) MaxChildBreadth() int {
	var breadth = 0
	for i := 0; i < len(d.Children); i++ {
		breadth = max(breadth, d.Children[i].CachedBreath())
	}
	return breadth
}

// MaxChildTotalFiles counts the total number of files under a given node
func (d *DiskInfo) MaxChildTotalFiles() int {
	var total = 0
	for i := 0; i < len(d.Children); i++ {
		total = max(total, d.Children[i].CachedTotalFiles())
	}
	return total
}
