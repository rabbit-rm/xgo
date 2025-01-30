// Package pool provides internal pool utilities.
package pool

import (
	"sync"
)

// Pool 针对 [sync.Pool] 的封装，提供强类型对象池
type Pool[T any] struct {
	pool sync.Pool
}

// New 返回一个新的存储 T 的 [Pool]，如果 [Pool] 中不存在 T，使用 NewFn 函数创建
func New[T any](fn func() T) *Pool[T] {
	return &Pool[T]{
		pool: sync.Pool{
			New: func() any {
				return fn()
			},
		},
	}
}

// Get 从 [Pool] 中获取一个 T如果池为空，则创建一个新的 T
func (p *Pool[T]) Get() T {
	return p.pool.Get().(T)
}

// Put 将 T 重新放入池中
func (p *Pool[T]) Put(x T) {
	p.pool.Put(x)
}
