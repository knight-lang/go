# Knight v3.0: Go Edition
This is a [Knight 3.0](https://github.com/knight-lang/knight-lang) implementation in Go. More details about Knight, its license, and specifications can be found in the [knight-lang](https://github.com/knight-lang/knight-lang) repo.

# Compiling
Simply run `go build .` to build it. You can then execute it via `./go (-e 'expr' | -f filename)`.

# Exemplar
This implementation is an "exemplar" implementation, so that people can get an idea of how Knight implementations might look. It has no fancy tricks or optimizations, and is thoroughly documented. If you don't know how to get started writing a Knight program, take a look at this one!

It does, however, provide some extensions that aren't required by the knight spec, as a convenience for people using it to execute Knight programs:

- Integers are 64 bit, instead of the required 32.
- UTF-8 is fully supported throughout, instead of just the required ASCII-subset that Knight requires
- Some forms of undefined behaviour are handled by returning `error`s (such as syntax errors, or type errors like `+ TRUE 1`). Not _all_ forms are handled, and things like integer wraparound are just ignored.
