package main

import (
	"bufio"
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
			for scanner.Scan() {
				line := strings.TrimSpace(scanner.Text())
				if line == "" {
					continue
				}

				if err := decodeXID(out, line); err != nil {
					fmt.Fprintln(os.Stderr, "Decode error:", err)
					os.Exit(1)
				}
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
		if _, err := xid.FromString(validate); err != nil {
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
			id := xid.New()
			fmt.Fprintf(out, "XID:         %s\n", id.String())
			fmt.Fprintf(out, "Timestamp:   %s\n", id.Time().Format(time.RFC3339))
			fmt.Fprintf(out, "Machine ID:  %x\n", id.Machine())
			fmt.Fprintf(out, "Process ID:  %d\n", id.Pid())
			fmt.Fprintf(out, "Counter:     %d\n", id.Counter())
			fmt.Fprintln(out)
		} else {
			// TODO uudashr: we need a way to set machine ID, process ID, and counter
			fmt.Fprintln(out, xid.New().String())
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
