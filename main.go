package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/rs/xid"
)

type outFormat int // Output format

const (
	outFormatHex outFormat = iota
	outFormatBinary
)

var outFormatNames = [...]string{
	"hex",
	"binary",
}

func (f outFormat) String() string {
	return outFormatNames[f]
}

func (f *outFormat) Set(value string) error {
	for i, name := range outFormatNames {
		if value == name {
			*f = outFormat(i)
			return nil
		}
	}

	return fmt.Errorf("Invalid output format: %s", value)
}

func (f *outFormat) UnmarshalText(text []byte) error {
	return f.Set(string(text))
}

func (f outFormat) MarshalText() ([]byte, error) {
	return []byte(f.String()), nil
}

func main() {
	var (
		verbose bool

		decode   string
		validate string
		count    int

		outFmt outFormat

		outFile string

		// machineID string
		// processID int
		// counter   int
	)

	flag.BoolVar(&verbose, "v", false, "Verbose output")

	flag.StringVar(&decode, "decode", "", "Decode xid")
	flag.StringVar(&validate, "validate", "", "Validate xid")
	flag.IntVar(&count, "n", 1, "Generate n xid")

	flag.TextVar(&outFmt, "format", outFormatHex, "Output format [hex, binary]")

	flag.StringVar(&outFile, "o", "", "Output file")

	// flag.StringVar(&machineID, "machine-id", "", "Machine ID")
	// flag.IntVar(&processID, "process-id", 0, "Process ID")
	// flag.IntVar(&counter, "counter", 0, "Counter")

	flag.Parse()

	out := os.Stdout

	if outFile != "" {
		f, err := os.OpenFile(outFile, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Open file error:", err)
			os.Exit(1)
			return
		}

		defer f.Close()

		out = f
	}

	// Decode
	if decode != "" {
		if decode == "-" {
			scanner := bufio.NewScanner(os.Stdin)
			var i int
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line == "" {
					continue
				}

				if err := decodeXID(out, line); err != nil {
					fmt.Fprintln(os.Stderr, "Decode error:", err)
					os.Exit(1)
				}

				i++
			}

			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "Read stdin error:", err)
				os.Exit(1)
				return
			}

			return
		}

		if err := decodeXID(out, decode); err != nil {
			fmt.Fprintln(os.Stderr, "Decode error:", err)
			os.Exit(1)
			return
		}

		return
	}

	// Validate
	if validate != "" {
		if validate == "-" {
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line == "" {
					continue
				}

				err := validateXID(line)
				if errors.Is(err, xid.ErrInvalidID) {
					if verbose {
						fmt.Fprintf(out, "Invalid %q\n", line)
					} else {
						fmt.Fprintf(os.Stderr, "Invalid %q\n", line)
					}

					os.Exit(1)
					return
				}

				if err != nil {
					if verbose {
						fmt.Fprintf(out, "%q %v\n", line, err)
					} else {
						fmt.Fprintf(os.Stderr, "%q %v\n", line, err)
					}
					os.Exit(1)
					return
				}
			}

			if err := scanner.Err(); err != nil {
				fmt.Fprintln(os.Stderr, "Read stdin error:", err)
				os.Exit(1)
				return
			}

			return
		}

		err := validateXID(validate)
		if errors.Is(err, xid.ErrInvalidID) {
			if verbose {
				fmt.Fprintln(out, "Invalid ID")
			} else {
				fmt.Fprintln(os.Stderr, "Invalid ID")
			}

			os.Exit(1)
			return
		}

		if err != nil {
			if verbose {
				fmt.Fprintln(out, err)
			} else {
				fmt.Fprintln(os.Stderr, err)
			}
			os.Exit(1)
			return
		}

		return
	}

	// Generate
	for i := 0; i < count; i++ {
		if verbose {
			// TODO uudashr: we need a way to set machine ID, process ID, and counter
			if i > 0 {
				fmt.Fprintln(out)
			}

			generateXIDVerbose(out)
		} else {
			// TODO uudashr: we need a way to set machine ID, process ID, and counter
			generateXID(out)
		}
	}
}

func decodeXID(w io.Writer, hex string) error {
	id, err := xid.FromString(hex)
	if err != nil {
		return err
	}

	fmt.Fprintf(w, "Timestamp:   %s\n", id.Time().Format(time.RFC3339))
	fmt.Fprintf(w, "Machine ID:  %x\n", id.Machine())
	fmt.Fprintf(w, "Process ID:  %d\n", id.Pid())
	fmt.Fprintf(w, "Counter:     %d\n", id.Counter())
	return nil
}

func validateXID(hex string) error {
	_, err := xid.FromString(hex)
	return err
}

func generateXID(w io.Writer) {
	fmt.Fprintln(w, xid.New().String())
}

func generateXIDVerbose(w io.Writer) {
	id := xid.New()
	fmt.Fprintf(w, "XID:         %s\n", id.String())
	fmt.Fprintf(w, "Timestamp:   %s\n", id.Time().Format(time.RFC3339))
	fmt.Fprintf(w, "Machine ID:  %x\n", id.Machine())
	fmt.Fprintf(w, "Process ID:  %d\n", id.Pid())
	fmt.Fprintf(w, "Counter:     %d\n", id.Counter())
}
