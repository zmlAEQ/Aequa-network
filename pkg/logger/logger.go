package logger

import (
    "encoding/json"
    "log"
    "os"
    "time"
)

var std = log.New(os.Stdout, "", 0)

func Info(msg string)  { std.Printf("%s INFO %s\n", time.Now().Format(time.RFC3339Nano), msg) }
func Warn(msg string)  { std.Printf("%s WARN %s\n", time.Now().Format(time.RFC3339Nano), msg) }
func Error(msg string) { std.Printf("%s ERRO %s\n", time.Now().Format(time.RFC3339Nano), msg) }

func InfoJ(msg string, fields map[string]any)  { logJ("info", msg, fields) }
func WarnJ(msg string, fields map[string]any)  { logJ("warn", msg, fields) }
func ErrorJ(msg string, fields map[string]any) { logJ("error", msg, fields) }

func logJ(level, msg string, fields map[string]any) {
    if fields == nil { fields = map[string]any{} }
    fields["level"] = level
    fields["ts"] = time.Now().Format(time.RFC3339Nano)
    fields["msg"] = msg
    b, _ := json.Marshal(fields)
    std.Println(string(b))
}
