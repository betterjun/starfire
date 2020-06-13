package wp

// 工作接口
type Worker interface {
	Do()
}

// 工作池
type WorkerPool struct {
	workerCh   chan Worker
	workerSize int
}

var defaultWorkerPool *WorkerPool

func InitDefaultWorkerPool(size int) (wp *WorkerPool) {
	defaultWorkerPool = NewWorkerPool(size)
	defaultWorkerPool.Start()
	return defaultWorkerPool
}

func GetDefaultWorkerPool() (wp *WorkerPool) {
	return defaultWorkerPool
}

func NewWorkerPool(size int) (wp *WorkerPool) {
	wp = &WorkerPool{
		workerSize: size,
		workerCh:   make(chan Worker, size),
	}

	return wp
}

func (wp *WorkerPool) Start() {
	go func(ch chan Worker) {
		for {
			w := <-ch
			if w == nil {
				break
			}
			w.Do()
		}
	}(wp.workerCh)
}

func (wp *WorkerPool) Stop() {
	close(wp.workerCh)
}

func (wp *WorkerPool) Put(w Worker) {
	wp.workerCh <- w
}

func (wp *WorkerPool) Status() (length, capacity int) {
	return len(wp.workerCh), cap(wp.workerCh)
}
