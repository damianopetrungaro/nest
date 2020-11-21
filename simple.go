package nest

import (
	"io"
)

// A SimpleWriter represents an active nestable simplified writer.
// Each write operation makes a single call to
// the Writer's Write method.
// A SimpleWriter can be used simultaneously from multiple goroutines;
// it guarantees to serialize access to the buffer.
type SimpleWriter struct {
	Writer *Writer
}

// NewSimpleWriter creates a new SimpleWriter with a inner Writer.
func NewSimpleWriter() *SimpleWriter {
	return &SimpleWriter{
		Writer: New(),
	}
}

// Child creates a new SimpleWriter from the current one.
func (s *SimpleWriter) Child(str string) *SimpleWriter {
	return &SimpleWriter{
		Writer: WithTitledParent(s.Writer, []byte(str)),
	}
}

// Write wraps a call to a Writer.WriteString.
func (s *SimpleWriter) Write(str string) {
	_, _ = s.Writer.WriteString(str)
}

// WriteTo wraps a call to the inner Writer's WriteTo method.
// Once the data on the SimpleWriter is fully written,
// then the data of each Children is gonna be written
// The return value n is the number of bytes written; it always fits into an
// int, but it is int64 to match the io.WriterTo interface. Any error
// encountered during the write is also returned.
func (s *SimpleWriter) WriteTo(w io.Writer) (n int64, err error) {
	return s.Writer.WriteTo(w)
}
