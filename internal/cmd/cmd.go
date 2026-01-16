package cmd

import (
	"flag"
	"io"
	"strings"
	"time"

	"basanos/internal/executor"
	"basanos/internal/fs"
	"basanos/internal/runner"
	"basanos/internal/sink"
	"basanos/internal/tree"
)

type stringSlice []string

func (s *stringSlice) String() string {
	return strings.Join(*s, ",")
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

type Config struct {
	SpecDir     string
	Outputs     []string
	Filter      string
	ShowHelp    bool
	ShowVersion bool
}

type RunOptions struct {
	Config     *Config
	FileSystem fs.FileSystem
	Executor   executor.Executor
	Stdout     io.Writer
	OutputFS   sink.WritableFS
}

func Run(opts RunOptions) error {
	if opts.FileSystem == nil {
		return nil
	}
	specTree, err := tree.LoadSpecTree(opts.FileSystem, opts.Config.SpecDir)
	if err != nil {
		return err
	}
	runID := time.Now().Format("2006-01-02_150405")
	var sinks []sink.Sink
	for _, output := range opts.Config.Outputs {
		if strings.HasPrefix(output, "json") {
			sinks = append(sinks, sink.NewJsonStreamSink(opts.Stdout))
		}
		if strings.HasPrefix(output, "files") {
			path := "runs"
			if _, after, found := strings.Cut(output, ":"); found {
				path = after
			}
			var writableFS sink.WritableFS
			if opts.OutputFS != nil {
				writableFS = opts.OutputFS
			} else {
				writableFS = fs.NewOSWritableFS(path)
			}
			sinks = append(sinks, sink.NewFileSink(writableFS, runID))
		}
		if strings.HasPrefix(output, "junit") {
			sinks = append(sinks, sink.NewJunitSink(opts.Stdout))
		}
	}
	r := runner.NewRunner(opts.Executor, sinks...)
	r.Filter = opts.Config.Filter
	return r.RunWithID(runID, specTree)
}

func ParseArgs(args []string) (*Config, error) {
	config := &Config{}

	var outputs stringSlice
	fs := flag.NewFlagSet("basanos", flag.ContinueOnError)
	fs.StringVar(&config.SpecDir, "s", "spec", "spec directory")
	fs.StringVar(&config.SpecDir, "spec", "spec", "spec directory")
	fs.Var(&outputs, "o", "output sink")
	fs.Var(&outputs, "output", "output sink")
	fs.StringVar(&config.Filter, "f", "", "filter pattern")
	fs.StringVar(&config.Filter, "filter", "", "filter pattern")
	fs.BoolVar(&config.ShowHelp, "h", false, "show help")
	fs.BoolVar(&config.ShowHelp, "help", false, "show help")
	fs.BoolVar(&config.ShowVersion, "v", false, "show version")
	fs.BoolVar(&config.ShowVersion, "version", false, "show version")
	fs.Parse(args)

	if len(outputs) == 0 {
		config.Outputs = []string{"files"}
	} else {
		config.Outputs = outputs
	}

	return config, nil
}
