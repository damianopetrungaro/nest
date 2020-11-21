package nest

import (
	"os"
	"sync"
)

func ExampleWriter() {
	base := New()

	orderedList := WithTitledParent(base, []byte("This is the start of the ordered list"))

	one := WithTitledParent(orderedList, []byte("1. Item one"))
	if _, err := one.Write([]byte("1.1 Written item")); err != nil {
		panic(err)
	}

	two := WithTitledParent(orderedList, []byte("2. Item two"))
	if _, err := two.Write([]byte("2.1 Written item")); err != nil {
		panic(err)
	}
	if _, err := two.Write([]byte("2.1 Written item")); err != nil {
		panic(err)
	}

	three := WithTitledParent(orderedList, []byte("3. Item three"))
	if _, err := three.Write([]byte("3.1 Written item")); err != nil {
		panic(err)
	}

	unorderedList := WithTitledParent(base, []byte("This is the start of the unordered list"))
	if _, err := unorderedList.Write([]byte("- Item one")); err != nil {
		panic(err)
	}
	if _, err := unorderedList.Write([]byte("- Item two")); err != nil {
		panic(err)
	}
	if _, err := unorderedList.Write([]byte("- Item three")); err != nil {
		panic(err)
	}

	if _, err := base.WriteTo(os.Stdout); err != nil {
		panic(err)
	}
	// Output:
	// This is the start of the ordered list
	//     1. Item one
	//         1.1 Written item
	//     2. Item two
	//         2.1 Written item
	//         2.1 Written item
	//     3. Item three
	//         3.1 Written item
	// This is the start of the unordered list
	//     - Item one
	//     - Item two
	//     - Item three
}

func ExampleWriter_concurrent() {
	wg := sync.WaitGroup{}
	base := New()

	orderedList := WithTitledParent(base, []byte("This is the start of the ordered list"))

	wg.Add(1)
	go func() {
		defer wg.Done()

		wg.Add(1)
		go func() {
			defer wg.Done()
			one := WithTitledParent(orderedList, []byte("1. Item one"))
			if _, err := one.Write([]byte("1.1 Written item")); err != nil {
				panic(err)
			}
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			two := WithTitledParent(orderedList, []byte("2. Item two"))
			if _, err := two.Write([]byte("2.1 Written item")); err != nil {
				panic(err)
			}
			if _, err := two.Write([]byte("2.1 Written item")); err != nil {
				panic(err)
			}

		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			three := WithTitledParent(orderedList, []byte("3. Item three"))
			if _, err := three.Write([]byte("3.1 Written item")); err != nil {
				panic(err)
			}
		}()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		unorderedList := WithTitledParent(base, []byte("This is the start of the unordered list"))
		if _, err := unorderedList.Write([]byte("- Item one")); err != nil {
			panic(err)
		}
		if _, err := unorderedList.Write([]byte("- Item two")); err != nil {
			panic(err)
		}
		if _, err := unorderedList.Write([]byte("- Item three")); err != nil {
			panic(err)
		}
	}()

	wg.Wait()
	if _, err := base.WriteTo(os.Stdout); err != nil {
		panic(err)
	}
	// Unordered output:
	// This is the start of the ordered list
	//     1. Item one
	//         1.1 Written item
	//     2. Item two
	//         2.1 Written item
	//         2.1 Written item
	//     3. Item three
	//         3.1 Written item
	// This is the start of the unordered list
	//     - Item one
	//     - Item two
	//     - Item three
}

func ExampleNew() {
	_ = New()
}

func ExampleWithParent() {
	n := New()
	n = WithParent(n)

	if _, err := n.WriteString("indented"); err != nil {
		panic(err)
	}
	if _, err := n.WriteTo(os.Stdout); err != nil {
		panic(err)
	}

	// Output:
	//     indented
}

func ExampleWithTitledParent() {
	n := New()
	n = WithTitledParent(n, []byte("Section 1"))

	if _, err := n.WriteString("indented"); err != nil {
		panic(err)
	}
	if _, err := n.WriteTo(os.Stdout); err != nil {
		panic(err)
	}

	// Output:
	// Section 1
	//     indented
}

func ExampleWriter_Write() {
	n := New()
	if _, err := n.Write([]byte("line one")); err != nil {
		panic(err)
	}
	if _, err := n.WriteTo(os.Stdout); err != nil {
		panic(err)
	}
	// Output:
	// line one
}

func ExampleWriter_Write_multiline() {
	n := New()
	if _, err := n.Write([]byte("line one\nline two")); err != nil {
		panic(err)
	}
	if _, err := n.WriteTo(os.Stdout); err != nil {
		panic(err)
	}
	// Output:
	// line one
	// line two
}

func ExampleWriter_WriteString() {
	n := New()
	if _, err := n.WriteString("line one"); err != nil {
		panic(err)
	}
	if _, err := n.WriteTo(os.Stdout); err != nil {
		panic(err)
	}

	// Output:
	// line one
}

func ExampleWriter_WriteString_multiline() {
	n := New()
	if _, err := n.WriteString("line one\nline two"); err != nil {
		panic(err)
	}
	if _, err := n.WriteTo(os.Stdout); err != nil {
		panic(err)
	}

	// Output:
	// line one
	// line two
}

func ExampleWriter_WriteTo() {
	n := New()

	if _, err := n.WriteTo(os.Stdout); err != nil {
		panic(err)
	}

	// Output:
}
