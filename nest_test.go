package nest

import (
	"bytes"
	"io/ioutil"
	"sync"
	"testing"
)

func TestNew(t *testing.T) {
	type want struct {
		depth    uint8
		children int
	}

	n := New()
	n2_0 := WithParent(n)
	n2_1 := WithParent(n)
	n3 := WithParent(n2_0)

	tests := map[string]struct {
		nest *Writer
		want want
	}{
		"Writer with no children and two level of depth": {
			nest: n3,
			want: want{
				depth:    2,
				children: 0,
			},
		},
		"Writer with one children and one level of depth": {
			nest: n2_0,
			want: want{
				depth:    1,
				children: 1,
			},
		},
		"Writer with no children and one level of depth": {
			nest: n2_1,
			want: want{
				depth:    1,
				children: 0,
			},
		},
		"Writer with two children and no level of depth": {
			nest: n,
			want: want{
				depth:    0,
				children: 2,
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			if test.nest.Buf == nil {
				t.Error("could not find buffer in the Writer")
			}

			if test.nest.Depth != test.want.depth {
				t.Error("could not match depth in the Writer")
				t.Errorf("got: %d", test.nest.Depth)
				t.Errorf("want: %d", test.want.depth)
			}

			if len(test.nest.Children) != test.want.children {
				t.Error("could not match children in the Writer")
				t.Errorf("got: %d", len(test.nest.Children))
				t.Errorf("want: %d", test.want.children)
			}
		})
	}
}

func TestWriter_Write(t *testing.T) {
	tests := map[string]struct {
		nest *Writer
		p    []byte
		want []byte
	}{
		"one line string with a zero depth Writer": {
			nest: New(),
			p:    []byte("a long string!"),
			want: []byte("a long string!\n"),
		},
		"one line string with a one depth Writer": {
			nest: WithParent(New()),
			p:    []byte("a long string!"),
			want: []byte("    a long string!\n"),
		},
		"string with new lines with a zero depth Writer": {
			nest: New(),
			p:    []byte("a long string!\nwith a new line!"),
			want: []byte("a long string!\nwith a new line!\n"),
		},
		"string with new lines with a two depth Writer": {
			nest: WithParent(WithParent(New())),
			p:    []byte("a long string!\nwith a new line!"),
			want: []byte("        a long string!\n        with a new line!\n"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			i, err := test.nest.Write(test.p)
			if i != len(test.want) {
				t.Error("could not match bytes written")
				t.Errorf("got: %d", i)
				t.Errorf("want: %d", len(test.p))
			}
			if err != nil {
				t.Errorf("could not write string to Writer: %s", err)
			}
			if !bytes.Equal(test.want, test.nest.Buf.Bytes()) {
				t.Error("could not match string written")
				t.Errorf("got: %s", test.nest.Buf.String())
				t.Errorf("want: %s", test.want)
			}
		})
	}
}

func TestWriter_WriteString(t *testing.T) {
	tests := map[string]struct {
		nest *Writer
		s    string
		want string
	}{
		"one line string with a zero depth Writer": {
			nest: New(),
			s:    "a long string!",
			want: "a long string!\n",
		},
		"one line string with a one depth Writer": {
			nest: WithParent(New()),
			s:    "a long string!",
			want: "    a long string!\n",
		},
		"string with new lines with a zero depth Writer": {
			nest: New(),
			s:    "a long string!\nwith a new line!",
			want: "a long string!\nwith a new line!\n",
		},
		"string with new lines with a two depth Writer": {
			nest: WithParent(WithParent(New())),
			s:    "a long string!\nwith a new line!",
			want: "        a long string!\n        with a new line!\n",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			i, err := test.nest.WriteString(test.s)
			if i != len(test.want) {
				t.Error("could not match bytes written")
				t.Errorf("got: %d", i)
				t.Errorf("want: %d", len(test.s))
			}
			if err != nil {
				t.Errorf("could not write string to Writer: %s", err)
			}
			if test.want != test.nest.Buf.String() {
				t.Error("could not match string written")
				t.Errorf("got: %s", test.nest.Buf.String())
				t.Errorf("want: %s", test.want)
			}
		})
	}
}

func TestWriter_WriteTo(t *testing.T) {
	tests := map[string]struct {
		i    int64
		nest *Writer
	}{
		"single simple": {
			41,
			&Writer{
				Buf: bytes.NewBuffer([]byte("this is the content\n that I want to print")),
			},
		},
		"simple with two depth": {
			26,
			&Writer{
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
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			w := &bytes.Buffer{}
			i, err := test.nest.WriteTo(w)
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

func TestWriterRaceConditions(t *testing.T) {
	if !raceEnabled {
		t.Skip("race detector is not enabled")
	}

	const pool = 10_000
	base := New()
	wg := sync.WaitGroup{}
	wg.Add(pool)
	for i := 1; i <= pool; i++ {
		go func(i int) {
			_,_ = base.WriteString("hello")
			_,_ = base.WriteTo(ioutil.Discard)
			wg.Done()
		}(i)
	}
	wg.Wait()
}
