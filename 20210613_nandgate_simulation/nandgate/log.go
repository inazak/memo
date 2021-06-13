package nandgate

import (
  "log"
  "os"
  "strings"
)

var Logger *log.Logger

func init() {
  Logger = log.New(os.Stderr, "## ", 0)
}

type trace struct {
  name   string
  format string
  lines  []*Line
}
var traces []*trace

func Trace(name, format string, lines ...*Line) {
  t := &trace{name: name, format: format, lines: lines}
  traces = append(traces, t)
}

func Log(prefix string) {
  for _, t := range traces {
    if strings.HasPrefix(t.name, prefix) {
      var arg []interface{}
      for _, line := range t.lines {
        arg = append(arg, line.State)
      }
      Logger.Printf(t.name + ": " + t.format, arg...)
    }
  }
}

