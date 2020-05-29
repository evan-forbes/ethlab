bind was forked from go-ethereum's identically named package. Modification are limited to using a different template for generating bindings.

TODO: 
- Delete the code for generating java and objC bindings.
- Some serious general house cleaning 
1) change the input from slices to a slice of structs contain the same data
2) break up into smaller more specific functions to improve readability.