package cmuxrpc

import (
	"context"
	"errors"
	"io"
	"sync"

	"fmt"
	"time"

	"github.com/calvindc/Web3RpcHub/cmuxrpc/codec"
)

var closeSinkerTimeout = time.Second * 10

// ByteSinker stream缓存
type ByteSinker interface {
	// WriteCloser datas interface(1：Write + 2：Closer=CloseWithError)
	io.WriteCloser
	// CloseWithError 在查询耗尽前关闭一个查询讲返回一个err
	CloseWithError(error) error

	SetEncoding(re RequestEncoding)
}

var _ ByteSinker = (*ByteSink)(nil)

type ByteSink struct {
	w         *codec.Writer
	closedMu  sync.Mutex //不同于codec.Writer的mutex
	closed    error      //
	streamCtx context.Context
	pkt       codec.Packet //包的完整形式
}

func newByteSink(ctx context.Context, w *codec.Writer) *ByteSink {
	return &ByteSink{
		w:         w,
		streamCtx: ctx,
		pkt:       codec.Packet{},
	}
}

func (bs *ByteSink) Write(b []byte) (int, error) {
	bs.closedMu.Lock()
	defer bs.closedMu.Unlock()

	if bs.closed != nil {
		return 0, bs.closed
	}

	//  func Stream(ctx context.Context, out chan<- Value) error {
	//  	for {
	//  		v, err := DoSomething(ctx)
	//  		if err != nil {
	//  			return err
	//  		}
	//  		select {
	//  		case <-ctx.Done():
	//  			return ctx.Err()
	//  		case out <- v:
	//  		}
	//  	}
	//  }
	// 检查上次写入后接收器是否关闭,否则中断write
	select {
	case <-bs.streamCtx.Done():
		bs.closed = bs.streamCtx.Err()
		return 0, bs.closed
	default:
		//继续写入
	}

	if bs.pkt.Req == 0 {
		return -1, fmt.Errorf("req ID not set (Flag: %s)", bs.pkt.Flag)
	}

	bs.pkt.Body = b
	err := bs.w.WritePacket(bs.pkt)
	if err != nil {
		bs.closed = err
		return -1, err
	}
	return len(b), nil
}

func (bs *ByteSink) Close() error {
	return bs.CloseWithError(io.EOF)
}

func (bs *ByteSink) CloseWithError(err error) error {
	bs.closedMu.Lock()
	defer bs.closedMu.Unlock()

	if bs.closed != nil {
		return bs.closed
	}

	var closePkt codec.Packet
	var isStream = bs.pkt.Flag.Get(codec.FlagStream)
	if err == io.EOF || err == nil {
		closePkt = newEndOkayPacket(bs.pkt.Req, isStream)
	} else {
		var epkt error
		closePkt, epkt = newEndErrPacket(bs.pkt.Req, isStream, err)
		if epkt != nil {
			return fmt.Errorf("close bytesink: error building error packet for %s: %w", err, epkt)
		}
		bs.closed = err
	}

	// 容许等待Write过程（等待closed packets）
	var errc = make(chan error)
	go func() {
		errc <- bs.w.WritePacket(closePkt)
	}()

	select {
	case werr := <-errc:
		if werr != nil {
			bs.closed = werr
		}
		return werr
	case <-time.After(closeSinkerTimeout):
		bs.closed = errors.New("muxrpc: close timeout exceeded")
		return bs.closed
	}

	bs.closed = err
	return nil
}

// SetEncoding request's data-type
func (bs *ByteSink) SetEncoding(re RequestEncoding) {
	bs.closedMu.Lock()
	defer bs.closedMu.Unlock()
	encFlag, err := re.asCodecFlag()
	if err != nil {
		//todo for  invalid request encoding? no others
		panic(err)
	}
	if re == TypeBinary {
		bs.pkt.Flag = bs.pkt.Flag.Clear(codec.FlagJSON)
	}
	bs.pkt.Flag = bs.pkt.Flag.Set(encFlag)
}
