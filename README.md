# Introduction
dibuilder is a tool to generate a dependancy injection builder at build time. It is 
designed to generate a function that creates the components which will together make
up the running application. The function it generates will generally be a part of
package `main` and will be called as one of the first steps in starting the application.

Note that the functionality described here is mostly aspirational (i.e. it isn't yet
implemented). The basic functionality of dibuilder will be part of the
[v0.1 milestone](https://github.com/sbosnick/dibuilder/milestone/1).

# Simple Example
A simple example will help show what dibuilder is intened to do.

Assume a ficticious project at `github.com/sbosnick/myproject` has the following main.go
file:

```golang
package main

// go:generate dibuilder github.com/sbosnick/myproject/internal/components/...

import (
        "github.com/sbosnick/dibuilder/run"
)

func main() {
        run.BuildAndRun(buildRoot)
}

```

Running `go generate` will generate a `buildroot.go` file that has a `buildRoot`
function with the following signature:

```golang
func buildRoot() rootpkg.RootType
```

where `rootpkg.RootType` is a specific type discovered by dibuilder that implements
the following interface (from package `github.com/sbosnick/dibuilder/run`):

```golang
type Runner interface {
        Run()
}
```

Running `go generate` on this project will cause dibuilder to scan the code in
package `github.com/sbosnick/myproject/internal/components/` and subpackage of
that package looking for constructors (top-level functions whose name starts with
`New`). dibuilder will try to satify the the parameters to these constructors with
the results of calling other such constructors (this is the dependancy injection part)
and will end by returning the result of a final constructor whose return type implements
`Runner`.

# Getting Started
You can get dibuilder by executing

```
go get -u github.com/sbosnick/dibuilder
```

# Concepts
The following concepts will help to explain what dibuilder is doing.

| **Concept**        | **Description**                                                              |
|--------------------|------------------------------------------------------------------------------|
| Component          | the non-error return types of a constructor from a user specified package    |
| Container          | a type used by dibuilder to hold descriptions of Components                  |
| Requirement        | the types of the parameters to the constructor of a Component in a Container |
| Provided Component | the Components that can be constructed by a Container                        |
| Rooted Container   | a Container that has a specifc Component that is identified as the root      |
| Complete Container | a Container for which all of its Requirements are Provided by the Container  |

dibuilder scans the code in a user specified set of pacakages looking for constructors which it then
adds to a Container (notionally along with the constructor's Component). The Components and
Requirements of a Container will then form a directed graph. If the Container is a Rooted, Complete
Container then dibuilder will use an ordering of Components that puts required Components before the
Components that require them. dibuilder takes this ordering and generates code in `buildRoot` to call
each constructor with its required parameters and then return the root component.
