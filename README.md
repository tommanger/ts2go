# ts2go - TypeScript to Go Compiler

A proof-of-concept compiler that transpiles TypeScript code to Go. This project demonstrates the feasibility of converting TypeScript programs with basic types and functions into runnable Go code.

## Features

- **Lexical Analysis**: Tokenizes TypeScript source code
- **Parsing**: Builds an Abstract Syntax Tree (AST) from tokens using recursive descent parsing
- **Type Mapping**: Converts TypeScript types to Go equivalents with best-effort approach
- **Code Generation**: Produces clean, readable Go code

### Supported TypeScript Features

- **Types**: `string`, `number`, `boolean`, `any`, arrays
- **Functions**: Function declarations with typed parameters and return types
- **Variables**: `const`, `let`, `var` declarations with type annotations
- **Control Flow**: `if-else`, `for`, `while` loops
- **Operators**: Arithmetic (`+`, `-`, `*`, `/`, `%`), comparison (`<`, `>`, `<=`, `>=`, `==`, `!=`), logical (`&&`, `||`, `!`)
- **Expressions**: Binary, unary, function calls, member access, assignments
- **Literals**: Numbers, strings, booleans, arrays
- **Built-ins**: `console.log()` → `fmt.Println()`

### Type Mapping

| TypeScript | Go |
|------------|-------|
| `string` | `string` |
| `number` | `float64` |
| `boolean` | `bool` |
| `any` | `interface{}` |
| `type[]` | `[]type` |

## Installation

```bash
# Clone the repository
git clone https://github.com/tmanger/ts2go.git
cd ts2go

# Build the compiler
go build -o ts2go cmd/ts2go/main.go
```

## Usage

```bash
# Basic usage
./ts2go input.ts output.go

# The output filename is optional (defaults to input filename with .go extension)
./ts2go input.ts
```

## Examples

### Hello World

**TypeScript:**
```typescript
function greet(name: string): string {
  return "Hello, " + name + "!";
}

const message: string = greet("World");
console.log(message);
```

**Generated Go:**
```go
package main

import "fmt"

func greet(name string) string {
	return (("Hello, " + name) + "!")
}

func main() {
	var message string = greet("World")
	fmt.Println(message)
}
```

### Calculator

See `examples/calculator.ts` for a more complex example with multiple functions and arithmetic operations.

### Fibonacci

See `examples/fibonacci.ts` for an example demonstrating recursion and loops.

## Project Structure

```
ts2go/
├── cmd/ts2go/           # CLI entry point
│   └── main.go
├── internal/
│   ├── lexer/           # Tokenization
│   │   ├── lexer.go
│   │   └── token.go
│   ├── parser/          # AST construction
│   │   └── parser.go
│   ├── ast/             # AST node definitions
│   │   └── ast.go
│   ├── types/           # Type mapping system
│   │   └── types.go
│   └── codegen/         # Go code generation
│       └── codegen.go
├── examples/            # Example TypeScript files
│   ├── hello.ts
│   ├── calculator.ts
│   ├── fibonacci.ts
│   └── control-flow.ts
└── README.md
```

## Architecture

1. **Lexer**: Scans TypeScript source and produces tokens
2. **Parser**: Consumes tokens and builds an AST using recursive descent parsing
3. **Type Mapper**: Translates TypeScript types to Go types
4. **Code Generator**: Traverses the AST and emits Go code

## Limitations

This is a proof-of-concept compiler focused on basic features. Not supported:

- Classes and interfaces
- Generics
- Decorators
- Async/await
- ES6 imports/exports (beyond simple `export`)
- Advanced type features (union types, type guards, etc.)
- DOM APIs and Node.js built-ins
- Complex object types and destructuring

## Testing

Run the examples to verify the compiler:

```bash
# Build the compiler
go build -o ts2go cmd/ts2go/main.go

# Compile and run examples
./ts2go examples/hello.ts examples/hello.go
go run examples/hello.go

./ts2go examples/calculator.ts examples/calculator.go
go run examples/calculator.go

./ts2go examples/fibonacci.ts examples/fibonacci.go
go run examples/fibonacci.go

./ts2go examples/control-flow.ts examples/control-flow.go
go run examples/control-flow.go
```

## Future Enhancements

- Support for structs and interfaces
- Better error messages with line/column information
- Source maps for debugging
- Optimization passes
- Support for more TypeScript features
- Standard library compatibility layer

## License

MIT

## Contributing

This is a proof-of-concept project. Feel free to fork and experiment!
# ts2go
