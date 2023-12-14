package errors

import (
	"bytes"
	"fmt"
	"runtime"
	"sync"
)

var (
	stacktracePool = sync.Pool{
		New: func() interface{} {
			return &stacktrace{
				storage: make([]uintptr, 64),
			}
		},
	}
)

type stacktrace struct {
	pcs    []uintptr // program counters; always a subslice of storage
	frames *runtime.Frames

	// The size of pcs varies depending on requirements:
	// it will be one if the only the first frame was requested,
	// and otherwise it will reflect the depth of the call stack.
	//
	// storage decouples the slice we need (pcs) from the slice we pool.
	// We will always allocate a reasonably large storage, but we'll use
	// only as much of it as we need.
	storage []uintptr
}

// StacktraceDepth specifies how deep of a stack trace should be captured.
type StacktraceDepth int

const (
	// StacktraceFirst captures only the first frame.
	StacktraceFirst StacktraceDepth = iota

	// StacktraceFull captures the entire call stack, allocating more
	// storage for it if needed.
	StacktraceFull
)

// captureStacktrace captures a stack trace of the specified depth, skipping
// the provided number of frames. skip=0 identifies the caller of
// captureStacktrace.
//
// The caller must call Free on the returned stacktrace after using it.
func captureStacktrace(skip int, depth StacktraceDepth) *stacktrace {
	stack := stacktracePool.Get().(*stacktrace) // nolint: forcetypeassert

	switch depth {
	case StacktraceFirst:
		stack.pcs = stack.storage[:1]
	case StacktraceFull:
		stack.pcs = stack.storage
	default:
		stack.pcs = stack.storage[:depth]
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
	if depth == StacktraceFull {
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

// Free releases resources associated with this stacktrace
// and returns it back to the pool.
func (st *stacktrace) Free() {
	st.frames = nil
	st.pcs = nil
	stacktracePool.Put(st)
}

// Count reports the total number of frames in this stacktrace.
// Count DOES NOT change as Next is called.
func (st *stacktrace) Count() int {
	return len(st.pcs)
}

// Next returns the next frame in the stack trace,
// and a boolean indicating whether there are more after it.
func (st *stacktrace) Next() (_ runtime.Frame, more bool) {
	return st.frames.Next()
}

func takeStacktrace(skip int) string {
	return TakeStacktraceDepth(skip+1, StacktraceFull)
}

// TakeStacktraceDepth is used to get stacktrace as string.
func TakeStacktraceDepth(skip int, depth StacktraceDepth) string {
	stack := captureStacktrace(skip+1, depth)
	defer stack.Free()

	buffer := &bytes.Buffer{}

	stackfmt := newStackFormatter(buffer)
	stackfmt.FormatStack(stack)

	return buffer.String()
}

// stackFormatter formats a stack trace into a readable string representation.
type stackFormatter struct {
	b        *bytes.Buffer
	nonEmpty bool // whehther we've written at least one frame already
}

// newStackFormatter builds a new stackFormatter.
func newStackFormatter(b *bytes.Buffer) stackFormatter {
	return stackFormatter{b: b}
}

// FormatStack formats all remaining frames in the provided stacktrace -- minus
// the final runtime.main/runtime.goexit frame.
func (sf *stackFormatter) FormatStack(stack *stacktrace) {
	// Note: On the last iteration, frames.Next() returns false, with a valid
	// frame, but we ignore this frame. The last frame is a runtime frame which
	// adds noise, since it's only either runtime.main or runtime.goexit.
	for frame, more := stack.Next(); more; frame, more = stack.Next() {
		sf.FormatFrame(frame)
	}
}

// FormatFrame formats the given frame.
func (sf *stackFormatter) FormatFrame(frame runtime.Frame) {
	if sf.nonEmpty {
		sf.b.WriteByte('\n')
	}

	sf.nonEmpty = true
	sf.b.WriteString(frame.Function)
	sf.b.WriteByte('\n')
	sf.b.WriteByte('\t')
	sf.b.WriteString(frame.File)
	sf.b.WriteByte(':')
	fmt.Fprint(sf.b, frame.Line)
}
