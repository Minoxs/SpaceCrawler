package DiskExplorer

// CachedSize returns the value from the last call of Size
func (d *DiskInfo) CachedSize() uint64 {
	return d.size
}
