# Nest

Nest is a `io.Writer` implementation which allow you to nested content,
taking care of indenting the content and adding titles to each section.

An example of an output is this
```
This is the start of the ordered list
    1. Item one
        1.1 Written item
    2. Item two
        2.1 Written item
        2.1 Written item
    3. Item three
        3.1 Written item
This is the start of the unordered list
    - Item one
    - Item two
    - Item three
```  

To interact with simplified APIs, but not compliant with the `io.Writer` there is a `SimpleWriter` which allows an even simpler usage. 

### Examples

For the writer take a look at the `nest_doc_test.go` file, for the simple writer take a look at the `simple_doc_test.go` file
