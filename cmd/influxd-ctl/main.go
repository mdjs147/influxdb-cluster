package main

import (
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"

	"github.com/influxdata/influxdb/cmd"
	"github.com/influxdata/influxdb/cmd/influxd-ctl/add_data"
	"github.com/influxdata/influxdb/cmd/influxd-ctl/add_meta"
	"github.com/influxdata/influxdb/cmd/influxd-ctl/common"
	"github.com/influxdata/influxdb/cmd/influxd-ctl/help"
	"github.com/influxdata/influxdb/cmd/influxd-ctl/remove_data"
	"github.com/influxdata/influxdb/cmd/influxd-ctl/remove_meta"
	"github.com/influxdata/influxdb/cmd/influxd-ctl/show"
	"github.com/influxdata/influxdb/cmd/influxd-ctl/update_data"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	m := NewMain()
	if err := m.Run(os.Args[1:]...); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Main represents the program execution.
type Main struct {
	Stdin  io.Reader
	Stdout io.Writer
	Stderr io.Writer
}

// NewMain return a new instance of Main.
func NewMain() *Main {
	return &Main{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}
}

// Run determines and runs the command specified by the CLI args.
func (m *Main) Run(args ...string) error {
	cOpts, args, err := m.parseFlags(args)
	if err == flag.ErrHelp {
		return nil
	} else if err != nil {
		return err
	}
	name, args := cmd.ParseCommandName(args)

	// Extract name from args.
	switch name {
	case "", "help":
		if err := help.NewCommand().Run(args...); err != nil {
			return fmt.Errorf("help: %s", err)
		}
	case "add-data":
		cmd := add_data.NewCommand(cOpts)
		if err := cmd.Run(args...); err != nil {
			return fmt.Errorf("add-data: %s", err)
		}
	case "add-meta":
		cmd := add_meta.NewCommand(cOpts)
		if err := cmd.Run(args...); err != nil {
			return fmt.Errorf("add-meta: %s", err)
		}
	case "remove-data":
		cmd := remove_data.NewCommand(cOpts)
		if err := cmd.Run(args...); err != nil {
			return fmt.Errorf("remove-data: %s", err)
		}
	case "remove-meta":
		cmd := remove_meta.NewCommand(cOpts)
		if err := cmd.Run(args...); err != nil {
			return fmt.Errorf("remove-meta: %s", err)
		}
	case "show":
		cmd := show.NewCommand(cOpts)
		if err := cmd.Run(args...); err != nil {
			return fmt.Errorf("show: %s", err)
		}
	case "update-data":
		cmd := update_data.NewCommand(cOpts)
		if err := cmd.Run(args...); err != nil {
			return fmt.Errorf("update-data: %s", err)
		}
	default:
		return fmt.Errorf(`unknown command "%s"`+"\n"+`Run 'influxd-ctl help' for usage`+"\n\n", name)
	}

	return nil
}

func (m *Main) parseFlags(args []string) (*common.Options, []string, error) {
	options := &common.Options{}
	fs := flag.NewFlagSet("", flag.ContinueOnError)
	fs.StringVar(&options.BindAddr, "bind", "localhost:8091", "Bind HTTP address of a meta node")
	fs.BoolVar(&options.BindTLS, "bind-tls", false, "Use TLS")
	fs.BoolVar(&options.SkipTLS, "k", false, "Skip certificate verification (ignored without -bind-tls)")
	fs.StringVar(&options.ConfigPath, "config", "", "Config file path")
	fs.StringVar(&options.AuthType, "auth-type", "none", "Type of authentication to use (none, basic, jwt)")
	fs.StringVar(&options.Username, "user", "", "User name (ignored without -auth-type basic | jwt)")
	fs.StringVar(&options.Password, "pwd", "", "Password (ignored without -auth-type jwt)")
	fs.StringVar(&options.Secret, "secret", "", "JWT shared secret (ignored without -auth-type jwt)")
	fs.Usage = func() { help.NewCommand().Run(args...) }
	if err := fs.Parse(args); err != nil {
		return options, args, err
	}
	return options, fs.Args(), nil
}
