glitch
======

Automatically, build, test, vet and install your go package.

Status
------

This was written quickly, and without tests. I was also learning
how to use the fsnotify package. There are some fundamental / 
obvious race conditions in the current version, however it's
working well in practice, so I'm not worrying about it for now.

Use at your own risk, but let me know if it's abusing your system somehow.

Usage
-----

```shell
go install github.com/levicook/glitch
cd <your go package root>
glitch
```

Behavior
--------

glitch will `go build`, `go vet ./...`, `go test ./...` and `go install` your go package.
If any one of these steps fail, it stops on that step, and waits for you to fix the issue.

When things go well, your output should look like this: 

```shell
2013/09/15 20:38:07 glitch: building
2013/09/15 20:38:07 glitch: build OK - vetting
2013/09/15 20:38:07 glitch: vet OK - testing
?       github.com/levicook/glitch      [no test files]
2013/09/15 20:38:07 glitch: test OK - installing
2013/09/15 20:38:08 glitch: install OK - wating for next build event
```

When something fails, like `go build`, you'll see less output. eg:

```shell
2013/09/15 20:42:57 glitch: building
# github.com/levicook/glitch
./main.go:83: syntax error: unexpected semicolon or newline, expecting )
```
