package nest

import (
	"bytes"
	"io/ioutil"
	"sync"
	"testing"
)

func TestNewSimpleWriter(t *testing.T) {
	type want struct {
		depth    uint8
		children int
	}

	n := NewSimpleWriter()
	n2_0 := n.Child("n2_0")
	n2_1 := n.Child("n2_1")
	n3 := n2_0.Child("n3")

	tests := map[string]struct {
		simple *SimpleWriter
		want   want
	}{
		"Writer with no children and two level of depth": {
			simple: n3,
			want: want{
				depth:    2,
				children: 0,
			},
		},
		"Writer with one children and one level of depth": {
			simple: n2_0,
			want: want{
				depth:    1,
				children: 1,
			},
		},
		"Writer with no children and one level of depth": {
			simple: n2_1,
			want: want{
				depth:    1,
				children: 0,
			},
		},
		"Writer with two children and no level of depth": {
			simple: n,
			want: want{
				depth:    0,
				children: 2,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.simple.Writer == nil {
				t.Error("could not find simple in the Writer")
			}

			if test.simple.Writer.Depth != test.want.depth {
				t.Error("could not match depth in the Writer")
				t.Errorf("got: %d", test.simple.Writer.Depth)
				t.Errorf("want: %d", test.want.depth)
			}

			if len(test.simple.Writer.Children) != test.want.children {
				t.Error("could not match children in the Writer")
				t.Errorf("got: %d", len(test.simple.Writer.Children))
				t.Errorf("want: %d", test.want.children)
			}
		})
	}
}

func TestSimple_Write(t *testing.T) {
	tests := map[string]struct {
		simple *SimpleWriter
		p      string
		want   []byte
	}{
		"one line string with a zero depth Writer": {
			simple: NewSimpleWriter(),
			p:      "a long string!",
			want:   []byte("a long string!\n"),
		},
		"string with new lines with a zero depth Writer": {
			simple: NewSimpleWriter(),
			p:      "a long string!\nwith a new line!",
			want:   []byte("a long string!\nwith a new line!\n"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.simple.Write(test.p)
			if !bytes.Equal(test.want, test.simple.Writer.Buf.Bytes()) {
				t.Error("could not match string written")
				t.Errorf("got: %s", test.simple.Writer.Buf.Bytes())
				t.Errorf("want: %s", test.want)
			}
		})
	}
}

func TestSimple_WriteTo(t *testing.T) {
	tests := map[string]struct {
		i      int64
		simple *SimpleWriter
	}{
		"single simple": {
			41,
			&SimpleWriter{&Writer{
				Buf: bytes.NewBuffer([]byte("this is the content\n that I want to print")),
			},
			}},
		"simple with two depth": {
			26,
			&SimpleWriter{&Writer{
				Buf: bytes.NewBuffer([]byte("one\n")),
				Children: []*Writer{
					{
						Depth: 1,
						Buf:   bytes.NewBuffer([]byte("    two\n")),
						Children: []*Writer{
							{
								Depth: 2,
								Buf:   bytes.NewBuffer([]byte("        three\n")),
							},
						},
					},
				},
			},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			w := &bytes.Buffer{}
			i, err := test.simple.WriteTo(w)
			if err != nil {
				t.Errorf("could not write simple content to the writer: %s", err)
			}
			if i != test.i {
				t.Error("could not match written bytes")
				t.Errorf("got: %d", i)
				t.Errorf("want: %d", test.i)
			}
		})
	}
}

func TestSimpleRaceConditions(t *testing.T) {
	if !raceEnabled {
		t.Skip("race detector is not enabled")
	}

	const pool = 10_000
	base := NewSimpleWriter()
	wg := sync.WaitGroup{}
	wg.Add(pool)
	for i := 1; i <= pool; i++ {
		go func(i int) {
			base.Write("hello")
			_, _ = base.WriteTo(ioutil.Discard)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
