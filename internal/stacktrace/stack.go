// Package stacktrace provides support for gathering stack traces efficiently.
package stacktrace

import (
	"runtime"

	"github.com/rabbit-rm/xgo/internal/buffer"
	"github.com/rabbit-rm/xgo/internal/pool"
)

var bufferPool = buffer.NewPool()

var _stackPool = pool.New(func() *Stack {
	return &Stack{
		storage: make([]uintptr, 64),
	}
})

// Depth 指定应捕获堆栈的深度
type Depth int

const (
	// First 捕获第一帧
	First Depth = iota

	// Full 捕获所有帧
	Full
)

// Stack 捕获堆栈信息
type Stack struct {
	pcs     []uintptr
	frames  *runtime.Frames
	storage []uintptr
}

// Free releases resources associated with this stacktrace
// and returns it back to the pool.
func (st *Stack) Free() {
	st.frames = nil
	st.pcs = nil
	_stackPool.Put(st)
}

// Count reports the total number of frames in this stacktrace.
// Count DOES NOT change as Next is called.
func (st *Stack) Count() int {
	return len(st.pcs)
}

// Next returns the next frame in the stack trace,
// and a boolean indicating whether there are more after it.
func (st *Stack) Next() (_ runtime.Frame, more bool) {
	return st.frames.Next()
}

// Capture 捕获指定 Depth 的堆栈跟踪，跳过提供的帧数。
// skip = 0 标识捕获的调用方
func Capture(skip int, depth Depth) *Stack {
	stack := _stackPool.Get()

	switch depth {
	case First:
		stack.pcs = stack.storage[:1]
	case Full:
		stack.pcs = stack.storage
	}

	// Unlike other "skip"-based APIs, skip=0 identifies runtime.Callers
	// itself. +2 to skip captureStacktrace and runtime.Callers.
	numFrames := runtime.Callers(
		skip+2,
		stack.pcs,
	)

	// runtime.Callers truncates the recorded stacktrace if there is no
	// room in the provided slice. For the full stack trace, keep expanding
	// storage until there are fewer frames than there is room.
	if depth == Full {
		pcs := stack.pcs
		for numFrames == len(pcs) {
			pcs = make([]uintptr, len(pcs)*2)
			numFrames = runtime.Callers(skip+2, pcs)
		}

		// Discard old storage instead of returning it to the pool.
		// This will adjust the pool size over time if stack traces are
		// consistently very deep.
		stack.storage = pcs
		stack.pcs = pcs[:numFrames]
	} else {
		stack.pcs = stack.pcs[:numFrames]
	}

	stack.frames = runtime.CallersFrames(stack.pcs)
	return stack
}

// Take 返回当前 Stack 的字符串表示形式
//
// skip 表示要跳过的帧数，skip=0 标识 跳过 Take
func Take(skip int, depth Depth) string {
	stack := Capture(skip+1, depth)
	defer stack.Free()

	buf := bufferPool.Get()
	defer buf.Free()

	formatter := NewFormatter(buf)
	switch depth {
	case First:
		frame, _ := stack.Next()
		formatter.FormatFrame(frame)
	case Full:
		formatter.FormatStack(stack)

	}
	return buf.String()
}

// Formatter formats a stack trace into a readable string representation.
type Formatter struct {
	b        *buffer.Buffer
	nonEmpty bool // 确保已经写入了一帧
}

// NewFormatter builds a new Formatter.
func NewFormatter(b *buffer.Buffer) Formatter {
	return Formatter{b: b}
}

// FormatStack formats all remaining frames in the provided stacktrace
func (sf *Formatter) FormatStack(stack *Stack) {
	for frame, more := stack.Next(); more; frame, more = stack.Next() {
		sf.FormatFrame(frame)
	}
}

// FormatFrame formats the given frame.
func (sf *Formatter) FormatFrame(frame runtime.Frame) {
	if sf.nonEmpty {
		sf.b.AppendByte('\n')
	}
	sf.nonEmpty = true
	sf.b.AppendString(frame.Function)
	sf.b.AppendByte('\n')
	sf.b.AppendByte('\t')
	sf.b.AppendString(frame.File)
	sf.b.AppendByte(':')
	sf.b.AppendInt(int64(frame.Line))
}
