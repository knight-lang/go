# Knight v2.0.1: Go Edition 
This is a [Knight 2.0.1](https://github.com/knight-lang/knight-lang) implementation in Go. More details about Knight, its license, and specifications can be found in the [knight-lang](https://github.com/knight-lang/knight-lang) repo.

# Compiling
Simply run `go build .` to build it. You can then execute it via `./go (-e 'expr' | -f filename)`.

## Exemplar
This implementation is meant to be one of the "exemplar" implementations, very clean and easy to understand. We catch a fair amount of errors, but none of them are really required, and are just a convenience for users who want to write knight code.

## UTF-8
Knight only requires a subset of ASCII to be supported, but this project supports arbitrary unicode characters as a fun extension.
