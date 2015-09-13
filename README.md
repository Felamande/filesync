# Filesync

Filesync is a simple tool to sync files between multiple directory pairs.

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

## Features
* Synchronize simultaneously from left directory to right directory.
* Install as system service.
* Run in system tray.
* A web frontend dashboard.

## TODO
* Support different protocols. By now, it can only sync file locally.

## Known issues
* On windows, when some programs create temprorary files of an editing file and then remove them, the filesync program will crash.
* Cannot synchronize file removing properly.

## Thanks
* [martini](https://github.com/go-martini/martini)
* [fsnotify](https://gopkg.in/fsnotify.v1)
* [log4go](https://code.google.com/p/log4go)





