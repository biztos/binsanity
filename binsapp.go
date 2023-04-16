// binsapp.go -- binsanity app running code called from the thin binary.

package binsanity

import (
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/urfave/cli/v2"
)

var ExitFunc = os.Exit              // exit function, override for testing
var OutWriter io.Writer = os.Stdout // Standard Output, override for testing
var ErrWriter io.Writer = os.Stderr // Standard Error, override for testing

const Version = "v0.2.0"

const AppDescription = `Convert asset files to Go source with test coverage.

The most common use case is to simply run "binsanity ASSET_DIR" in the Go
source directory where you plan to use the asset functions.  ASSET_DIR is the
directory of data to be included 

The default values will usually work if you have an up-to-date go.mod file in
the current directory or above it.  The files generated in the working dir
will be binsanity.go and binsanity_test.go.

The generated source and text files will be overwritten if they exist.

Hats off to Jim Teeuwen for doing the much more powerful version of this thing
first.  If you aren't too sensitive about testing you should probably use
go-bindata instead: https://pkg.go.dev/github.com/jteeuwen/go-bindata
`

// RunApp uses the cli package to run the app with the args provided, which
// should be os.Args equivalent.  Process is called with the parsed options.
//
// If it encounters any error, it exits with a nonzero value though ExitFunc.
// Standard output and error are written to OutWriter and ErrWriter.
func RunApp(args []string) {

	// cli is pretty good but its error handling is weird (not as weird as
	// docopt's so there's that). We avoid the craziness of things like
	// cli.HandleExitCoder(cli.Exit(err, 1)) by doing our own craziness.
	cfg := &Config{}
	app := &cli.App{
		Name:        "binsanity",
		Usage:       "embed assets with testing",
		UsageText:   "binsanity [options] ASSET_DIR",
		Description: AppDescription,
		Version:     Version,
		Writer:      OutWriter,
		ErrWriter:   ErrWriter,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "output",
				Aliases:     []string{"o"},
				Value:       "binsanity.go",
				Usage:       "output file for generated Go source",
				Destination: &(cfg.File),
				Required:    false,
			},
			&cli.StringFlag{
				Name:        "module",
				Aliases:     []string{"m"},
				Value:       "",
				Usage:       "module identifier (see description)",
				Destination: &(cfg.Module),
				Required:    false,
			},
			&cli.StringFlag{
				Name:        "package",
				Aliases:     []string{"p"},
				Value:       "",
				Usage:       "package name (see description)",
				Destination: &(cfg.Package),
				Required:    false,
			},
		},
		Action: func(cCtx *cli.Context) error {
			// Surprised this isn't built in to the app spec...
			if cCtx.NArg() != 1 {
				return errors.New("Single arg required: ASSET_DIR")
			}
			cfg.Dir = cCtx.Args().Get(0)

			res, err := Process(cfg)
			if err != nil {
				return err
			}

			fmt.Fprintln(OutWriter, res.String())

			return nil
		},
	}

	if err := app.Run(args); err != nil {
		fmt.Fprintln(ErrWriter, err)
		ExitFunc(1)
	}

}
