stack
=====

Stack is a tiny programming language for scripts. It is stack based,
interpreted, dynamically but strongly typed, homoiconic and has first class
functions. The standard library contains functions for file management, regular
expressions and networking. So far it looks to be write-only.

Source code is made of space delimited words. If a word is a value, it is pushed into the stack. If it's a function, it's executed and given free access to the stack, usually popping the required arguments and pushing the result.

The interpreter can run source code files or used as an interactive interpreter (REPL). The parser, interpreter and libraries are written in Go and the default map implementation (`%` or `map`) is parallel.

Examples
--------

    1 2 3
    # Result: 1 2 3
    
Pushes three number on the stack and terminates the program. The values left in the stack are printed.

    1 1 +
    # Result: 2
    
Pushes two 1's to the stack and runs a sum operation, popping the two operands and pushing the result. 


    1 1 + 2 *
    # Result: 4
    
Pushes two 1's, add them, push a 2 and multiply.

    [1 1 +] eval
    # Result: 2
    
Pushes the list `[1 1 +]`, evaluates it and pushes the result.

    # If format: condition if_true if_false ?
    1 2 < ["1 is less than 2"] ["1 is more than 2"] ?
    # Result: "1 is less than 2"
    
Runs the comparison function `<` with 1 and 2 and conditionally evaluates one of the lists depending on the result.

    ["stored" "values"] "foobar" $
    "foobar" @
    # Result: "stored" "values"
    
Binds (`$`) the string `foobar` to two values. When called (`@`) it pushes the bound values. The linebreak is optional.

    [1 +] "increment" $
    5 "increment" @
    # Result: 6
    
Notice it binds the values inside the list, and not the list itself. This allows binding of partial applications of a function, effectively declaring a new function.
    
    1 . 2 .
    # Result: 1 1 2 2
    
The dot operator pushes two duplicates of the popped value.
    
    [. 5 < [1 + "up5" @] [] ?] "up5" $
    1 "up5" @
    # Result: 5
    
Bound values are resolved at runtime, allowing recursive functions or redeclarations.


Syntactic sugar
---------------

- `:` starts a string until the next space (`:myvar "myvar" =`).
- Strings can be delimited with single quotes.
- A linebreak or tab can be used anywhere that a space is acceptable.
- Spaces between words are collapsed.
- Built-in functions with special characters are available as names too (`1 dup` is `1 .`).
- Comments start with a `#` and go until the end of the line.


Library
-------

This is a toy language but it has a few useful functions in its library. They are mainly about regular expressions, fetching HTTP pages and manipulating files:

    "abbbc" "ab+c" matches
    # Result: true
    
    "abc! abb acd" "(ab.)(!?)" findall
    # Result: [["abc" "!"] ["abb" ""]]
    
    "google.com/robots.txt" get "Disallow: /search" contains
    # Result: true
    
    "google.com/robots.txt" get "Disallow" "Forbidden" replace "robots.txt" write
    # Result: creates a file with the edited robots.txt content
    
    "robots.txt" . read "Forbidden" "Disallow" replace & write
    # Result: runs a find and replace in the robots.txt file saved
    
    "robots.txt" delete
