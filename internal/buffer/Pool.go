package buffer

import (
	"github.com/rabbit-rm/xgo/internal/pool"
)

type Pool struct {
	p *pool.Pool[*Buffer]
}

func NewPool() Pool {
	return Pool{
		p: pool.New(func() *Buffer {
			return &Buffer{
				bs: make([]byte, 0, _size),
			}
		}),
	}
}

func (p Pool) Get() *Buffer {
	buf := p.p.Get()
	buf.Reset()
	buf.pool = p
	return buf
}

func (p Pool) put(buf *Buffer) {
	p.p.Put(buf)
}
