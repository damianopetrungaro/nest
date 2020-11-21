package nest

import (
	"bytes"
	"io"
	"sync"
)

// A Writer represents an active nestable writer.
// Each write operation makes a single call to
// the bytes.Buffer's Write method.
// A Writer can be used simultaneously from multiple goroutines;
// it guarantees to serialize access to the buffer.
type Writer struct {
	Buf      *bytes.Buffer
	Depth    uint8
	Children []*Writer

	mutex sync.Mutex
}

// New creates a new Writer with a inner buffer.
func New() *Writer {
	return &Writer{
		Buf: &bytes.Buffer{},
	}
}

// WithParent creates a new Writer from a parent one.
func WithParent(parent *Writer) *Writer {
	return WithTitledParent(parent, nil)
}

// WithTitledParent creates a new Writer with a title.
// The title appears as a fist non-indented line on top of the content.
func WithTitledParent(parent *Writer, t []byte) *Writer {
	child := &Writer{
		Buf:   bytes.NewBuffer(format(t, int(parent.Depth))),
		Depth: parent.Depth + 1,
	}
	parent.mutex.Lock()
	parent.Children = append(parent.Children, child)
	parent.mutex.Unlock()
	return child
}

// Write wraps a call to the inner bytes.Buffer's Write method.
// The content p is formatted and indented depending on the Depth of the Writer.
func (n *Writer) Write(p []byte) (int, error) {
	p = format(p, int(n.Depth))
	n.mutex.Lock()
	defer n.mutex.Unlock()
	return n.Buf.Write(p)
}

// WriteString wraps a call to a Writer.Write.
func (n *Writer) WriteString(s string) (int, error) {
	return n.Write([]byte(s))
}

// WriteTo wraps a call to the inner bytes.Buffer's WriteTo method.
// Once the data on the Writer is fully written,
// then the data of each Children is gonna be written
// The return value i is the number of bytes written; it always fits into an
// int, but it is int64 to match the io.WriterTo interface. Any error
// encountered during the write is also returned.
func (n *Writer) WriteTo(w io.Writer) (i int64, err error) {
	n.mutex.Lock()
	func() {
		defer n.mutex.Unlock()
		ii, writeErr := n.Buf.WriteTo(w)
		i = i + ii
		if writeErr != nil {
			err = writeErr
			return
		}
	}()

	for _, child := range n.Children {
		ii, writeErr := child.WriteTo(w)
		i = i + ii
		if writeErr != nil {
			err = writeErr
			return
		}
	}

	return
}

func format(p []byte, depth int) []byte {
	if len(p) == 0 {
		return p
	}

	//TODO: optimize p2 slice allocation
	var p2, prefix []byte
	prefix = bytes.Repeat([]byte{' ', ' ', ' ', ' '}, depth)

	for i, line := range bytes.Split(p, []byte{'\n'}) {
		if i > 0 {
			p2 = append(p2, '\n')
		}
		p2 = append(p2, prefix...)
		p2 = append(p2, line...)
	}

	p2 = append(p2, '\n')
	return p2
}
