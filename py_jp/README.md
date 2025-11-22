# Python_jp

Python framework for Japanese language programming with custom syntax and preprocessor.

---

## Overview

Framework enabling Python programming using Japanese keywords and syntax. Includes a preprocessor that translates Japanese-style Python code into standard Python, along with CLI tools for running .jppy files.

---

## Building

Install from source:
```bash
python setup.py install
```

Or with pip:
```bash
pip install .
```

---

## Running

Execute .jppy files:
```bash
python_jp run examples/hello_world.jppy
```

Or using Python module:
```bash
python -m Python_jp.cli run examples/hello_world.jppy
```

---

## Features

- Japanese keyword support
- Custom syntax preprocessing
- CLI for running .jppy scripts
- Translation between Python and Python_jp
- Custom exception handling
- Utility functions for Japanese text processing

---

## Project Structure

```
py_jp/
├── Python_jp/
│   ├── cli.py           # Command-line interface
│   ├── keywords.py      # Japanese keyword definitions
│   ├── preprocessor.py  # Code preprocessing
│   ├── translator.py    # Translation utilities
│   ├── runner.py        # Script execution
│   ├── exceptions.py    # Custom exceptions
│   └── utils.py         # Helper functions
├── examples/
│   ├── hello_world.jppy # Hello world example
│   └── test.jppy        # Test script
├── bin/
│   └── python_jp        # Executable script
└── setup.py             # Installation script
```

---

## Examples

Check the `examples/` directory for sample .jppy files demonstrating the framework's capabilities.

---

## License

MIT License
