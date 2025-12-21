# While Language Interpreter

A complete parser and interpreter for the **While** programming language, implemented in Haskell.

## What is While?

While is a simple imperative language used in programming language theory. It features:

- Variables and assignments (`x := 5`)
- Arithmetic expressions (`+`, `-`, `*`, `/`)
- Boolean expressions and comparisons
- Conditionals (`if-then-else-end`)
- While loops (`while-do-end`)
- Sequential composition (`;`)

## Quick Start

### Build
```bash
stack build
```

### Run a program file
```bash
stack exec while-lang -- examples/factorial.while
```

### Start interactive REPL
```bash
stack exec while-lang
```

## Example Programs

### Factorial (examples/factorial.while)
```
n := 5;
result := 1;
while n > 0 do
    result := result * n;
    n := n - 1
end
// Result: result = 120
```

### Fibonacci (examples/fibonacci.while)
```
n := 10;
a := 0;
b := 1;
while n > 0 do
    temp := a + b;
    a := b;
    b := temp;
    n := n - 1
end;
fib := a
// Result: fib = 55
```

## Language Syntax

### Statements
```
x := <expr>           -- assignment
skip                  -- no operation
<stmt> ; <stmt>       -- sequence
if <bool> then <stmt> else <stmt> end    -- conditional
while <bool> do <stmt> end               -- loop
```

### Arithmetic Expressions
```
42                    -- integer literal
x                     -- variable
a + b, a - b          -- addition, subtraction
a * b, a / b          -- multiplication, division
(a + b) * c           -- parentheses for grouping
```

### Boolean Expressions
```
true, false           -- boolean literals
a = b, a < b, a > b   -- comparisons
a <= b, a >= b        -- comparisons
not b                 -- negation
b1 and b2             -- conjunction
b1 or b2              -- disjunction
```

### Comments
```
// This is a comment
x := 5  // inline comment
```

## REPL Commands

| Command | Description |
|---------|-------------|
| `:help`, `:h` | Show help |
| `:quit`, `:q` | Exit REPL |
| `:state`, `:s` | Show current state |
| `:clear`, `:c` | Clear state |
| `:load <file>` | Load and run a file |

## Project Structure

```
├── src/
│   ├── Main.hs                    -- Entry point
│   └── Language/While/
│       ├── Syntax.hs              -- AST data types
│       ├── Lexer.hs               -- Tokenization
│       ├── Parser.hs              -- Parsing
│       ├── Interpreter.hs         -- Evaluation
│       └── Repl.hs                -- Interactive REPL
├── examples/
│   ├── factorial.while            -- Factorial example
│   └── fibonacci.while            -- Fibonacci example
├── while-lang.cabal               -- Package definition
└── stack.yaml                     -- Stack configuration
```

## Dependencies

- `base` - Haskell base library
- `parsec` - Parser combinator library
- `containers` - Data structures (Map)

## License

BSD-3-Clause
