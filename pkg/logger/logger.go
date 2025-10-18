package logger

import (
    "log"
    "os"
)

var std = log.New(os.Stdout, "", log.LstdFlags|log.Lmicroseconds)

func Info(msg string)  { std.Printf("INFO %s", msg) }
func Warn(msg string)  { std.Printf("WARN %s", msg) }
func Error(msg string) { std.Printf("ERRO %s", msg) }
