# Mist Language

Mist (Mild Rust) is a simple programming language inspired by rust syntax and functional paradigm, written in Go using only Go’s standard libraries.

I wrote this while reading the Writing an [Interpreter in Go](https://interpreterbook.com), with tons of extra features (mentioned below).

# Table Of Contents

- [Features](#features)
- [Getting Started](#getting-started)
- [Usage](#usage)
- [License](#license)

# Features

This language has:

- Recursion
- Scopes and variable shadowing
- Currying
- Method-chaining
- Type system
- if, else if, else conditionals
- Binary Operators
- Annonymous functions
- Lists and Maps
- Functional-ish methods like mapand filter
- Builtin functions like max, len , print and range
- Precise error messages, pointing to the exact character/token that caused the error.
- Implicit returns

All of these features are demonstarted in the [examples](https://github.com/MohamedAbdeen21/Mist-Lang/tree/master/examples) folder. The extension .rs is just for syntax highlighting. Disable LSP temporarily to avoid rust-related error messages.

# Getting Started

One great thing about this language, is that it’s built using only the standard Go libraries. The only external dependency is for printing the AST in the terminal.

To get started:

1. Clone the repo 

    ```console
    git clone https://github.com/MohamedAbdeen21/Mist-Lang.git
    ```

2. Install [Go programming language](https://go.dev/doc/install)

3. Install the dependencies 

    ```cosnole
    cd Mist-Lang && go mod tidy
    ```

# Usage

To run a file

```
go run . examples/fibonacci.rs
```

Feel free to explore the [examples](https://github.com/MohamedAbdeen21/Mist-Lang/tree/master/examples) for sample usages.


# License

This project is licensed under the MIT License - see the LICENSE file for details.
