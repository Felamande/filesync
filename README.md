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
* Install as a system service.
* Run in the system tray.
* An easy-to-use dashboard.

## TODO
* Support different protocols. By now, it can only sync file locally.

## Known issues
* On windows, some program cannot save files that are occupied by the filesync program.
* if some temp files are created in the left dir and then be removed, those ones synchronized to the right directory won't.

## Thanks
* [martini](https://github.com/go-martini/martini)
* [fsnotify](https://gopkg.in/fsnotify.v1)
* [log4go](https://code.google.com/p/log4go)





