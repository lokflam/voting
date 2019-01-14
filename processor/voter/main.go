package main

import (
	"os"
	"syscall"

	"voting/processor/voter/handler"

	"github.com/hyperledger/sawtooth-sdk-go/logging"
	"github.com/hyperledger/sawtooth-sdk-go/processor"
	flags "github.com/jessevdk/go-flags"
)

// Opts represents options for running the processor
type Opts struct {
	Verbose []bool `short:"v" long:"verbose" description:"Increase verbosity"`
	Connect string `short:"C" long:"connect" description:"Validator component endpoint to connect to" default:"tcp://localhost:4004"`
	Queue   uint   `long:"max-queue-size" description:"Set the maximum queue size before rejecting process requests" default:"100"`
	Threads uint   `long:"worker-thread-count" description:"Set the number of worker threads to use for processing requests in parallel" default:"0"`
}

func main() {
	var opts Opts

	logger := logging.Get()

	parser := flags.NewParser(&opts, flags.Default)
	_, err := parser.Parse()
	if err != nil {
		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			logger.Error("Failed to parse args: ", err)
			os.Exit(2)
		}
	}

	var loggingLevel int
	switch len(opts.Verbose) {
	case 0:
		loggingLevel = logging.WARN
	case 1:
		loggingLevel = logging.INFO
	default:
		loggingLevel = logging.DEBUG
	}
	logger.SetLevel(loggingLevel)

	handler := &handler.AdminHandler{}
	processor := processor.NewTransactionProcessor(opts.Connect)
	processor.SetMaxQueueSize(opts.Queue)
	if opts.Threads > 0 {
		processor.SetThreadCount(opts.Threads)
	}
	processor.AddHandler(handler)
	processor.ShutdownOnSignal(syscall.SIGINT, syscall.SIGTERM)
	err = processor.Start()
	if err != nil {
		logger.Error("Processor stopped: ", err)
	}
}
