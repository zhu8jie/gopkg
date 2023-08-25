package xmulti

import (
	"context"
	"errors"
	"strconv"
	"time"

	"go.uber.org/zap"
)

type InputParam interface{}

type HandlerFunc func(context.Context, []InputParam) error

type Xmulti struct {
	fun              HandlerFunc
	limit            int // 并行个数，也做限流用
	limitChan        chan struct{}
	multiCnt         int           // 一次处理多少条数据
	multiTimeout     time.Duration // 批量处理超时设置，该时间内收到的数据还没到multiCnt数量则不等待，直接执行
	multiChan        chan struct{}
	multiParam       []InputParam
	multiTimeoutChan chan struct{}
	closeChan        chan struct{}
	retryCount       int // 重试次数(不包含默认的第一次)
	name             string
	logger           *zap.SugaredLogger
}

func New(f HandlerFunc, limit, multiCnt, retryCount int, name string, multiTimeout time.Duration, logger *zap.SugaredLogger) *Xmulti {
	if logger == nil {
		l, _ := zap.NewDevelopment()
		logger = l.Sugar()
	}

	if name == "" {
		name = "xMulti_" + strconv.FormatInt(time.Now().Unix(), 10)
		logger.Infof("xMulti name is new:%v", name)
	}
	if limit == 0 {
		// 默认限制100
		limit = 100
	}
	// 限制数最大不能超5000
	if limit > 5000 {
		limit = 5000
	}
	if multiCnt <= 0 {
		// 默认一次执行一个
		multiCnt = 1
	}

	if multiCnt > 1 && multiTimeout <= 0 {
		multiTimeout = time.Millisecond * 500
	}

	if retryCount > 3 {
		retryCount = 3 // 最多只能重试3次
	} else if retryCount < 0 {
		retryCount = 0
	}

	p := &Xmulti{
		fun:              f,
		limit:            limit,
		limitChan:        make(chan struct{}, limit),
		multiCnt:         multiCnt,
		multiTimeout:     multiTimeout,
		multiChan:        make(chan struct{}, multiCnt-1),
		multiParam:       make([]InputParam, 0, multiCnt),
		multiTimeoutChan: make(chan struct{}, 1),
		closeChan:        make(chan struct{}),
		retryCount:       retryCount,
		name:             name,
	}

	return p
}

func (p *Xmulti) Close() {
	close(p.closeChan)
	for {
		time.Sleep(time.Millisecond * 200)
		if len(p.limitChan) == 0 {
			break
		}
	}
}

func (p *Xmulti) Run(ctx context.Context, i InputParam) (err error) {
	select {
	case <-p.closeChan:
		return errors.New("xMulti is closed")
	default:
		p.multiTimeoutChan <- struct{}{}
		defer func() {
			<-p.multiTimeoutChan
		}()

		p.multiParam = append(p.multiParam, i)
		if len(p.multiParam) == 1 {
			p.timeoutRun(ctx)
		}

		select {
		case p.multiChan <- struct{}{}:
		default:
			p.runAction(ctx, p.multiParam)
			p.multiChan = make(chan struct{}, p.multiCnt-1)
			p.multiParam = make([]InputParam, 0, p.multiCnt)
		}
	}
	return nil
}

func (p *Xmulti) timeoutRun(ctx context.Context) {
	if p.multiTimeout > 0 {
		go func() {
			ctxWithTimeout, cancel := context.WithTimeout(context.Background(), p.multiTimeout)
			<-ctxWithTimeout.Done()
			cancel()
			p.runTimeout(ctx)
		}()
	}
}

func (p *Xmulti) runTimeout(ctx context.Context) {
	p.multiTimeoutChan <- struct{}{}
	defer func() {
		<-p.multiTimeoutChan
	}()
	if len(p.multiParam) > 0 {
		p.runAction(ctx, p.multiParam)
		p.multiChan = make(chan struct{}, p.multiCnt-1)
		p.multiParam = make([]InputParam, 0, p.multiCnt)
	}
}

func (p *Xmulti) runAction(ctx context.Context, i []InputParam) {
	p.limitChan <- struct{}{}

	go func() {
		defer func() {
			if err := recover(); err != nil {
				p.logger.Errorf("xMulti||runAction panic||err=%v", p.name, err)
			}
		}()

		err := p.fun(ctx, i)

		if err != nil && p.retryCount > 0 {
			t := 0
			for {
				t++
				if t > p.retryCount {
					break
				}

				// 重试打印日志记录重试次数
				p.logger.Infof("xMulti||retry||name=%v||count=%v", p.name, t)

				// 单位：毫秒，分别有 [20, 80, 180] 3个阶梯的延时重试
				sleepStep := t * t * 20
				if sleepStep > 0 {
					time.Sleep(time.Millisecond * time.Duration(sleepStep))
				}

				err = p.fun(ctx, i)
				if err == nil {
					break
				}
			}
		}

		<-p.limitChan
	}()
}
