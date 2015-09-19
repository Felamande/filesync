# Filesync

Filesync is a simple tool to sync files between multiple uri pairs.

## Install from source
```
go get -u github.com/Felamande/filesync
go install github.com/Felamande/filesync
```
or use [gopm](http://gopm.io) instead to cross the GFW.
```
gopm get -g -u github.com/Felamande/filesync 
go install github.com/Felamande/filesync
```
run ```filesync -help``` to get help.

## Features
* Synchronize simultaneously from left directory to right directory.
* Install as a system service.
* Run in the shell.

## TODO
* Support different protocols. By now, it can only sync file locally.
* An easy-to-use dashboard.

## Known issues
* On windows, some program cannot save files that are occupied by the filesync program.
* if some temp files are created in the left dir and then be removed, those ones synchronized to the right directory won't.
* The service won't start when first installed using the command line, a manual control or a system restart is required.
* The service will crash when stopped.

## Thanks
* [martini](https://github.com/go-martini/martini)
* [fsnotify](https://gopkg.in/fsnotify.v1)
* [log4go](https://code.google.com/p/log4go)
* [service](https://github.com/kardianos/service)





