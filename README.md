glitch
======

Continuously and automatically, build, test, vet and install your go package.

Status
------

This was written quickly, and without tests. I was also learning
how to use the fsnotify package. There are some fundamental / 
obvious race conditions in the current version, however it's
working well in practice, so I'm not worrying about it for now.

Use at your own risk, but let me know if it's abusing your system somehow.

Installation
------------

```shell
go get -u github.com/levicook/glitch
go install github.com/levicook/glitch
```

Usage
-----

Make sure $GOPATH/bin is on your PATH, then simply:

```shell
cd <your go package>
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

There's nothing generic about this tool. It encapsulates a specific workflow and only
pays attention to .go files. If you want something different, use something totally
generic like guard. You're welcome to fork this, but I am unlikely to merge pull requests
looking for a different and/or pluggable behavior.


Known Issues
------------
Not a glitch issue per se; OSX users frequently see: too many open files
http://superuser.com/questions/433746/is-there-a-fix-for-the-too-many-open-files-in-system-error-on-os-x-10-7-1
