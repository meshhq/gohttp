#GOMETALINTER

https://stackoverflow.com/questions/49716/what-is-static-code-analysis

# Static code analysis? What is it?

Analyzing code without executing it. Generally used to find bugs or ensure conformance to coding guidelines. The classic example is a compiler which finds lexical, syntactic and even some semantic mistakes.

# Why Use it?

Static analysis tools should be used when they help maintain code quality. If they're used, they should be integrated into the build process, local development environment etc... otherwise they will be ignored.

# Pitfalls

1. Slows down dev because of dump forgetful errors. 
2. Tools take a long time to run. 

# Why

Go lang has a ton of statically checking go source code. Rather than running them all independently, and serially, this tool runs them all concurrently. 

# Installation

go get -u gopkg.in/alecthomas/gometalinter.v2

# Default Supported Linters

# Optional Supported Linters

# First Run 

gometalinter.v2

# Options

## Fast Run 

gometalinter.v2 --fast

## Fast Run 

 --skip=