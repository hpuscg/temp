/*
#Time      :  2020/12/3 10:02 上午
#Author    :  chuangangshen@deepglint.com
#File      :  zeroLog.go
#Software  :  GoLand
*/
package main

import (
	"github.com/rs/zerolog"
	"gopkg.in/natefinch/lumberjack.v2"
	"time"
)

var (
	Root *zerolog.Logger
)

type RotateLogger struct {
	lumberjack.Logger
	Level zerolog.Level
}

func NewRotateLogger() *RotateLogger {
	return &RotateLogger{}
}

func (this *RotateLogger) WriteLevel(level zerolog.Level, p []byte) (int, error) {
	if level >= this.Level {
		return this.Write(p)
	} else {
		return len(p), nil
	}
}


func main() {
	Init()
	Root.Error().Msg("yes")
}

func Init() {
	/*core := lumberjack.Logger{
		Filename:   "log",
		MaxSize:    30,
		MaxAge:     5,
		MaxBackups: 1,
		LocalTime:  false,
		Compress:   true,
	}*/


	writer := NewRotateLogger()
	writer.Level, _ = zerolog.ParseLevel("info")
	writer.Filename = "log"
	writer.MaxBackups = 30
	writer.MaxSize = 16
	writer.MaxAge = 14
	writer.Compress = true
	zlog := zerolog.New(writer).With().Timestamp().Logger()
	zlog = zlog.With().Caller().Logger()
	// zlog.ErrorStackMarshaler = pkgerrors.MarshalStack
	zlog = zlog.Output(zerolog.ConsoleWriter{
		Out: writer,
		TimeFormat: time.RFC3339,
	})

	Root = &zlog
}





