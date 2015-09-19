# Filesync

Filesync 是一个简单的工具用来将一个地址的内容同步到另一个地址，并且可以指定多对这样的地址对。

## 从源码安装
```
go get -u github.com/Felamande/filesync
go install github.com/Felamande/filesync
```
或使用[gopm](http://gopm.io) 在天朝获得更快的下载速度。
```
gopm get -g -u github.com/Felamande/filesync 
go install github.com/Felamande/filesync
```
在命令行工具中输入 ```filesync -help``` 获得详细用法。

## 功能
* 在多个地址对之间实时同步内容。
* 安装为系统服务，占用资源极小，完全静默同步。
* 直接在终端中运行.

## 待开发
* 不同的协议的支持。
* 简单易用的控制面板.

## 已知的问题
* 某些程序使用临时文件编辑原文件时，可能无法保存到原文件，因为filesync正在占用该文件。出现的情况较少，只针对特定少数几个程序。
* 某些程序创建临时文件编辑原文件时，filesync会将临时文件也同步到另一边地址里并且可能不会删除这些已经同步的临时文件.
* 第一次使用命令行安装并运行filesync服务时，服务安装成功但不会开始运行，需要手工启动服务或等待下一次重启。Windows用户可以使用Win+R打开service.msc找到名为"FileSync Service"的服务用右键启动。

## 本项目用到的开源项目
* [martini](https://github.com/go-martini/martini)
* [fsnotify](https://gopkg.in/fsnotify.v1)
* [log4go](https://code.google.com/p/log4go)
* [service](https://github.com/kardianos/service)





