package nest

import (
	"os"
	"sync"
)

func ExampleSimple() {
	base := NewSimpleWriter()

	orderedList := base.Child("This is the start of the ordered list")

	one := orderedList.Child("1. Item one")
	one.Write("1.1 Written item")

	two := orderedList.Child("2. Item two")
	two.Write("2.1 Written item")
	two.Write("2.1 Written item")

	three := orderedList.Child("3. Item three")
	three.Write("3.1 Written item")

	unorderedList := base.Child("This is the start of the unordered list")
	unorderedList.Write("- Item one")
	unorderedList.Write("- Item two")
	unorderedList.Write("- Item three")

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

func ExampleSimple_concurrent() {
	wg := sync.WaitGroup{}
	base := NewSimpleWriter()

	orderedList := base.Child("This is the start of the ordered list")

	wg.Add(1)
	go func() {
		defer wg.Done()

		wg.Add(1)
		go func() {
			defer wg.Done()
			one := orderedList.Child("1. Item one")
			one.Write("1.1 Written item")
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			two := orderedList.Child("2. Item two")
			two.Write("2.1 Written item")
			two.Write("2.1 Written item")
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			three := orderedList.Child("3. Item three")
			three.Write("3.1 Written item")
		}()
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		unorderedList := base.Child("This is the start of the unordered list")
		unorderedList.Write("- Item one")
		unorderedList.Write("- Item two")
		unorderedList.Write("- Item three")
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

func ExampleNewSimple() {
	_ = NewSimpleWriter()
}

func ExampleSimple_Child_with_content() {
	n := NewSimpleWriter()
	n = n.Child("")

	n.Write("indented")
	if _, err := n.WriteTo(os.Stdout); err != nil {
		panic(err)
	}

	// Output:
	//     indented
}

func ExampleSimple_Child_without_content() {
	n := NewSimpleWriter()
	n = n.Child("Section 1")

	n.Write("indented")
	if _, err := n.WriteTo(os.Stdout); err != nil {
		panic(err)
	}

	// Output:
	// Section 1
	//     indented
}

func ExampleSimple_Write() {
	n := NewSimpleWriter()
	n.Write("line one")
	if _, err := n.WriteTo(os.Stdout); err != nil {
		panic(err)
	}
	// Output:
	// line one
}

func ExampleSimple_Write_multiline() {
	n := NewSimpleWriter()
	n.Write("line one\nline two")
	if _, err := n.WriteTo(os.Stdout); err != nil {
		panic(err)
	}
	// Output:
	// line one
	// line two
}

func ExampleSimple_WriteTo() {
	n := NewSimpleWriter()

	if _, err := n.WriteTo(os.Stdout); err != nil {
		panic(err)
	}

	// Output:
}
