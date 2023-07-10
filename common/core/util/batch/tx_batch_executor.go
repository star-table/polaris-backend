package batch

import (
	"context"
	"time"

	"github.com/star-table/common/core/errors"

	"github.com/star-table/go-common/utils/stack"
	"github.com/star-table/common/core/logger"
	"upper.io/db.v3/lib/sqlbuilder"
)

var log = logger.GetDefaultLogger()

type TxExecFunc func(tx sqlbuilder.Tx) errors.SystemErrorInfo

// TxBatchExecutor TX批量执行器
type TxBatchExecutor struct {
	ctx      context.Context
	cancel   context.CancelFunc
	tx       sqlbuilder.Tx
	parallel int
	queue    chan TxExecFunc
}

func (c *TxBatchExecutor) Init(tx sqlbuilder.Tx, parallel int) {
	ctx, cancel := context.WithCancel(context.Background())
	c.ctx = ctx
	c.cancel = cancel
	c.tx = tx.WithContext(ctx)
	c.parallel = parallel
	c.queue = make(chan TxExecFunc, 100)
}

func (c *TxBatchExecutor) PushJobs(fs []TxExecFunc) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				log.Errorf("[TxBatchExecutor] [PushJobs] panic recovered: %v", r)
				log.Errorf("[TxBatchExecutor] %s", stack.GetStack())
			}
		}()
		defer close(c.queue)

		for i, f := range fs {
			select {
			case <-c.ctx.Done():
				log.Errorf("[TxBatchExecutor] PushJobs cancelled, total: %d, pushed: %d", len(fs), i)
				return
			case c.queue <- f:
				//c.log.Infof("[TxBatchExecutor] [%s] PushJobs pushed org: %d", c.Name, org.Id)
			}
		}
		log.Infof("[TxBatchExecutor] PushJobs finished")
	}()
}

func (c *TxBatchExecutor) WorkerFunc(ctx context.Context, cancel context.CancelFunc, tx sqlbuilder.Tx, i int) errors.SystemErrorInfo {
	//log.Infof("[TxBatchExecutor] [Worker|%d] begin to work.", i)
	for {
		select {
		case <-ctx.Done():
			//log.Infof("[TxBatchExecutor] [Worker|%d] cancelled working.", i)
			return nil
		case f, more := <-c.queue:
			if more {
				if err := f(tx); err != nil {
					return err
				}
			} else {
				//log.Infof("[TxBatchExecutor] [Worker|%d] finished working.", i)
				return nil
			}
		}
	}
}

func (c *TxBatchExecutor) StartAndWaitFinish() errors.SystemErrorInfo {
	//log.Info("[TxBatchExecutor] start.")
	startTime := time.Now()
	ch := make(chan errors.SystemErrorInfo, 1)
	defer close(ch)

	for i := 0; i < c.parallel; i++ {
		go func() {
			defer func() {
				if r := recover(); r != nil {
					log.Errorf("[TxBatchExecutor] [StartAndWaitFinish] %d panic recovered: %v", i, r)
					log.Errorf("[TxBatchExecutor] %s", stack.GetStack())
				}
			}()

			ch <- c.WorkerFunc(c.ctx, c.cancel, c.tx, i)
		}()
	}
	var err errors.SystemErrorInfo
	for i := 0; i < c.parallel; i++ {
		if e := <-ch; e != nil {
			c.cancel()
			err = e
		}
	}
	log.Infof("[TxBatchExecutor] finished, cost time: %v.", time.Now().Sub(startTime))
	return err
}
