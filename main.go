package main

import (
	"fmt"
	"log"
	"os"

	"github.com/jessevdk/go-flags"
	"textpad.com/app"
)

var (
	// Version is the version number or commit hash
	// These variables should be set by the linker when compiling
	Version = "0.0.0-unknown"
	// CommitHash is the git hash of last commit
	CommitHash = "Unknown"
	// CompileDate is the date of build
	CompileDate = "Unknown"
)

type Opts struct {
	Host         string `long:"host" env:"PAD_HOST" default:"0.0.0.0" description:"listening address"`
	Port         int    `long:"port" env:"PAD_PORT" default:"8080" description:"listening port"`
	DatabasePath string `long:"db" env:"PAD_DB_PATH" default:"db" description:"path to database files"`
	Verbose      bool   `long:"verbose" description:"verbose logging"`
	Version      bool   `short:"v" long:"version" description:"show the version number"`
}

func main() {
	var opts Opts

	p := flags.NewParser(&opts, flags.Default)
	if _, err := p.ParseArgs(os.Args[1:]); err != nil {
		os.Exit(1)
	}

	if opts.Version {
		fmt.Printf("Version: %s\nCommit hash: %s\nCompile date: %s\n", Version, CommitHash, CompileDate)
		os.Exit(0)
	}

	log.Printf("[DEBUG] opts: %+v", opts)
	if err := app.Execute("localhost", 8888, "db/bolt.db"); err != nil {
		log.Fatal(err)
	}

}
