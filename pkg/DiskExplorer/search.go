package DiskExplorer

import "context"

func (d *DiskInfo) BreadthSearch(ctx context.Context) {
	var sem = newSemaphore(5000)

	sem.Acquire()
	go d.breadthSearch(ctx, sem)

	sem.Wait()
}

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
