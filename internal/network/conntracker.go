package network

import (
	"context"
	"log"
	"net"
	"sync"
	"time"

	"github.com/calvindc/Web3RpcHub/netwrap"
	"github.com/calvindc/Web3RpcHub/secretstream"
)

type connEntry struct {
	c       net.Conn
	started time.Time //接入时间
	done    chan struct{}
	cancel  context.CancelFunc
}

// connLookupMap is feed - connEntry 只认最新conn的feed，比如存在账号异地登录
type connLookupMap map[[32]byte]connEntry

// connTracker tracks连接的active状态
type connTracker struct {
	activeLock sync.Mutex
	active     connLookupMap
}

// NewConnTracker new conn-tracker
func NewConnTracker() ConnTracker {
	return &connTracker{
		active: make(connLookupMap),
	}
}

// OnAccept 新建一个connect to tracker
func (ct *connTracker) OnAccept(ctx context.Context, newConn net.Conn) (bool, context.Context) {
	ct.activeLock.Lock()
	k := toActive(newConn.RemoteAddr())
	oldConn, ok := ct.active[k]
	ct.activeLock.Unlock()
	if ok {
		//closes the connection. Any blocked Read or Write operations will be unblocked and return errors.
		oldConn.c.Close()
		// tells other operation to abandon its work.
		oldConn.cancel()
		select {
		case <-oldConn.done:
			//清理
		case <-time.After(10 * time.Second):
			log.Println("[ConnTracker/lastWins] warning: not accepted, would ghost connection:", oldConn.c.RemoteAddr().String(), time.Since(oldConn.started))
			return false, nil
		}
	}

	ct.activeLock.Lock()
	ctx, cancel := context.WithCancel(ctx)
	ct.active[k] = connEntry{
		c:       newConn,
		started: time.Now(),
		done:    make(chan struct{}),
		cancel:  cancel,
	}
	ct.activeLock.Unlock()

	return true, ctx
}

// Active 对ConnTracker的Active, 当前是否处于活跃
func (ct *connTracker) Active(a net.Addr) (bool, time.Duration) {
	ct.activeLock.Lock()
	defer ct.activeLock.Unlock()

	k := toActive(a)
	l, ok := ct.active[k]
	if !ok {
		return false, 0
	}
	return true, time.Since(l.started)
}

// OnClose 通知tracker一个连接已经断开
func (ct *connTracker) OnClose(conn net.Conn) time.Duration {
	ct.activeLock.Lock()
	defer ct.activeLock.Unlock()

	k := toActive(conn.RemoteAddr())
	who, ok := ct.active[k]
	if !ok {
		return 0
	}
	close(who.done)
	delete(ct.active, k)
	return time.Since(who.started)
}

// CloseAll 对ConnTracker的CloseAll
func (ct *connTracker) CloseAll() {
	ct.activeLock.Lock()
	defer ct.activeLock.Unlock()

	for k, c := range ct.active {
		err := c.c.Close()
		if err != nil {
			log.Printf("failed to close conn for %x: %v\n", k[:5], err)
		}
		// stop immediately
		c.cancel()

		<-c.done
		delete(ct.active, k)
	}
}

// Count active状态的conn数量
func (ct *connTracker) Count() uint {
	ct.activeLock.Lock()
	defer ct.activeLock.Unlock()

	return uint(len(ct.active))
}

func toActive(a net.Addr) [32]byte {
	var pk [32]byte
	shs, ok := netwrap.GetAddr(a, "shs-bs").(secretstream.Addr)
	if !ok {
		panic("not an SHS connection")
	}
	copy(pk[:], shs.PubKey)
	return pk
}

type trackerLastWins struct {
	connTracker
}

/*// NewLastWinsTracker 强制释放以前的连接，然后让新的连接加入
func NewLastWinsTracker() ConnTracker {
	return &trackerLastWins{connTracker{active: make(connLookupMap)}}
}*/
