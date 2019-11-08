package main

import (
	"fmt"
	"github.com/metajar/expect"
	"github.com/sirupsen/logrus"
	"os"
	"regexp"
	"time"
)

func init() {
	// Log as JSON instead of the default ASCII formatter.
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	logrus.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	logrus.SetLevel(logrus.WarnLevel)
}

func main() {
	ssh, err := expect.Spawn("ssh", "10.205.172.2", "-l", "root")
	if err != nil {
		fmt.Println(err)
	}
	ssh.SetTimeout(10 * time.Second)
	ssh.SetLogger(LogRUsLogger())
	const PROMPT = `.*#`

	// Login to the Device
	ssh.Expect(`[Pp]assword:`)
	ssh.SendMasked("root") // SendMasked hides from logging
	ssh.Send("\n")
	ssh.Expect(PROMPT) // Wait for prompt
	ssh.SendLn("term len 0")
	ssh.Expect(PROMPT)
	ssh.SendLn("show interfaces")
	m, _ := ssh.Expect(PROMPT)
	fmt.Println(m)

}

func LogRUsLogger() expect.Logger {
	return &logRusLogger{
		logger: logrus.New(),
	}
}

type logRusLogger struct {
	logger *logrus.Logger
}

func fmtTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.999 -0700 MST")
}

func (logger *logRusLogger) Send(t time.Time, data []byte) {
	logger.logger.WithFields(logrus.Fields{
		"time":      fmtTime(t),
		"direction": "Send",
		"data":      string(data),
	}).Info()
}

func (logger *logRusLogger) SendMasked(t time.Time, _ []byte) {
	logger.logger.WithFields(logrus.Fields{
		"time":      fmtTime(t),
		"direction": "Send",
		"data":      "*** MASKED ***",
	}).Info()
}

func (logger *logRusLogger) Recv(t time.Time, data []byte) {
	logger.logger.WithFields(logrus.Fields{
		"time":      fmtTime(t),
		"direction": "Recv",
		"data":      string(data),
	}).Info()
}

func (logger *logRusLogger) RecvNet(t time.Time, data []byte) {
	// This is likely too verbose.
}

func (logger *logRusLogger) RecvEOF(t time.Time) {
	logger.logger.WithFields(logrus.Fields{
		"time":      fmtTime(t),
		"direction": "Recv",
		"data":      "EOF",
	}).Info()

}

func (logger *logRusLogger) ExpectCall(t time.Time, r *regexp.Regexp) {
	logger.logger.WithFields(logrus.Fields{
		"time":      fmtTime(t),
		"direction": "Recv",
		"data":      fmt.Sprintf("%v", r),
	}).Info()
}

func (logger *logRusLogger) ExpectReturn(t time.Time, m expect.Match, e error) {
	logger.logger.WithFields(logrus.Fields{
		"time":      fmtTime(t),
		"direction": "Recv",
		"data":      fmt.Sprintf(" ExpectReturn %q %v", m, e),
	}).Info()
}

func (logger *logRusLogger) Close(t time.Time) {
	logger.logger.WithFields(logrus.Fields{
		"time":      fmtTime(t),
		"direction": "CLOSE",
		"data":      "EOF",
	}).Info()

}
