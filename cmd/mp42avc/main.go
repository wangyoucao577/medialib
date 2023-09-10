package main

import (
	"flag"
	"fmt"

	"github.com/golang/glog"
	"github.com/wangyoucao577/medialib/util/dump"
	"github.com/wangyoucao577/medialib/util/exit"
)

func main() {
	flag.Parse()
	defer glog.Flush()

	// validate and get flags
	if err := validateFlags(); err != nil {
		glog.Error(err)
		exit.Fail()
	}
	contentType, _ := getConentType()

	if flags.outputFilePath == dump.OutputStdout {
		defer fmt.Println() // new line to avoid `%` displayed at the end in Mac shell
	}

	if err := parseMP4(flags.inputFilePath, contentType, flags.outputFilePath); err != nil {
		glog.Error(err)
		exit.Fail()
	}
}
