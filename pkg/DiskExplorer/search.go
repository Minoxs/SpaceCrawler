package DiskExplorer

import "context"

// BreadthSearch will iterate over the DiskInfo tree breadth first.
// This will spawn up to 5000 concurrent go threads to do this search.
// In other words, this might be a very resource hungry call.
//
// This function is NOT thread-safe. Blocks until completion.
// Cancel the given context to stop the search early.
func (d *DiskInfo) BreadthSearch(ctx context.Context) {
	var sem = newSemaphore(5000)

	sem.Acquire()
	go d.breadthSearch(ctx, sem)

	sem.Wait()
}

// breadthSearch is an inner helper function to run the search in.
func (d *DiskInfo) breadthSearch(ctx context.Context, sem *semaphore) {
	defer sem.Release()
	if ctx.Err() != nil {
		return
	}

	d.Expand()
	for i := 0; i < d.Breadth(); i++ {
		if !d.Children[i].IsDir {
			continue
		}
		sem.Acquire()
		go d.Children[i].breadthSearch(ctx, sem)
	}
}
