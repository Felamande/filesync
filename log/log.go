package log

import (
	"database/sql"
	"fmt"
	"time"

	"code.google.com/p/log4go"
	_ "github.com/mattn/go-sqlite3"
	//"io"
)

type Level int

const (
	FINEST Level = iota
	FINE
	DEBUG
	TRACE
	INFO
	WARNING
	ERROR
	CRITICAL
	PANIC
)

type Logger interface {
	Info(source, message string)
	Debug(source, message string)
	Warn(source, message string)
	Error(source, message string)
	Critical(source, message string)
	Panic(source, message string)
	Close() error
}

type DBLogger struct {
	db *sql.DB
}

func NewDBLogger(name string) *DBLogger {

	db, err := sql.Open("sqlite3", name)
	if err != nil {
		return nil
	}

	db.Exec(`
		create table log(time datetime not null,level tinyint not null,source text, message text)
	`)

	return &DBLogger{db}
}

func (l *DBLogger) write(lv Level, source, message string) {
	tx, err := l.db.Begin()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	stmt, err := tx.Prepare("insert into log(time,level,source,message) values(?,?,?,?)")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(time.Now().Format("2006-01-02 15:04:05"), int(lv), source, message)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	tx.Commit()
}

func (l *DBLogger) Close() error {
	return l.Close()
}

func (l *DBLogger) Debug(source, message string) {
	l.write(DEBUG, source, message)
	fmt.Printf("%v : %v", source, message)
}
func (l *DBLogger) Info(source, message string) {
	l.write(INFO, source, message)
}

func (l *DBLogger) Warn(source, message string) {
	l.write(WARNING, source, message)
}

func (l *DBLogger) Critical(source, message string) {
	l.write(CRITICAL, source, message)
}

func (l *DBLogger) Error(source, message string) {
	l.write(ERROR, source, message)
}

func (l *DBLogger) Panic(source, message string) {
	l.write(PANIC, source, message)
}

type FileLogger struct {
	logger4go *log4go.FileLogWriter
}

func NewFileLogger(file string) *FileLogger {
	l := log4go.NewFileLogWriter(file, true)
	l.SetRotateDaily(true)
	return &FileLogger{l}
}

func (l *FileLogger) Close() error {
	l.logger4go.Close()
	return nil
}

func (l *FileLogger) Debug(source, message string) {
	l.logger4go.LogWrite(&log4go.LogRecord{log4go.DEBUG, time.Now(), source, message})
}
func (l *FileLogger) Info(source, message string) {
	l.logger4go.LogWrite(&log4go.LogRecord{log4go.INFO, time.Now(), source, message})
}

func (l *FileLogger) Warn(source, message string) {
	l.logger4go.LogWrite(&log4go.LogRecord{log4go.WARNING, time.Now(), source, message})
}

func (l *FileLogger) Critical(source, message string) {
	l.logger4go.LogWrite(&log4go.LogRecord{log4go.CRITICAL, time.Now(), source, message})
}

func (l *FileLogger) Error(source, message string) {
	l.logger4go.LogWrite(&log4go.LogRecord{log4go.ERROR, time.Now(), source, message})
}

func (l *FileLogger) Panic(source, message string) {
	l.logger4go.LogWrite(&log4go.LogRecord{log4go.CRITICAL, time.Now(), source, message})
}
