stack
=====

Stack is a tiny programming language for scripts. It is stack based,
interpreted, dynamically but strongly typed, homomorphic and has first class
functions. The standard library contains functions for file management, regular
expressions and networking. So far it looks to be write-only.

Source code is made of space delimited words. If a word is a value, it is pushed into the stack. If it's a function, it's executed and given free access to the stack, usually popping the required arguments and pushing the result.

The interpreter can run source code files or used as an interactive interpreter (REPL). The parser, interpreter and libraries are written in Go.

Examples
--------

    1 1 +
    # Result: 2
    
Pushes two 1's to the stack and runs a sum operation, popping the two operands and pushing the result. 


    1 1 + 2 *
    # Result: 4
    
Pushes two 1's, add them, push a 2 and multiply them.
