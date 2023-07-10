package go_pool

/// 手动实现协程池
/// 参考 https://iswbm.com/291.html

type GoPool struct {
	work chan func()
	// 数量
	num chan struct{}
}

// poolSize int 该值必须大于 0！
func NewGoPool(poolSize int) *GoPool {
	return &GoPool{
		work: make(chan func()),
		num:  make(chan struct{}, poolSize),
	}
}

func (p *GoPool) NewTask(task func()) {
	select {
	case p.work <- task:
	case p.num <- struct{}{}:
		go p.worker(task)
	}
}

func (p *GoPool) worker(task func()) {
	defer func() {
		<-p.num
	}()
	for {
		task()
		task = <- p.work
	}
}
