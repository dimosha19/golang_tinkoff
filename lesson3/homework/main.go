package main

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"math"
	"os"
)

const DEFAULT_BLOCK_SIZE int = 1024
const DEFAULT_OFFSET int = 0
const DEFAULT_LIMIT int = -1

type Options struct {
	FileInput  string
	FileOutput string
	Offset     int
	Limit      int
	BlockSize  int
	From       io.Reader
	To         io.Writer
}

func optionsInit(options *Options) {
	var err error
	if options.FileInput != "" {
		options.From, err = os.Open(options.FileInput)
		check(err)
	} else {
		options.From = os.Stdin
	}
	if options.FileOutput != "" {
		if _, err := os.Stat(options.FileOutput); errors.Is(err, os.ErrNotExist) {
			options.To, err = os.OpenFile(options.FileOutput, os.O_RDONLY|os.O_CREATE, 0644)
			check(err)
		} else {
			panic(errors.New("File already exists"))
		}
	} else {
		options.To = os.Stdout
	}
	if options.Offset < 0 {
		panic(errors.New("offset must be positiv"))
	}
}

func ParseFlags() (*Options, error) {
	var opts Options

	flag.StringVar(&opts.FileInput, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.FileOutput, "to", "", "file to write. by default - stdout")
	flag.IntVar(&opts.Offset, "offset", DEFAULT_OFFSET, "number of bytes to skip when copying")
	flag.IntVar(&opts.Limit, "limit", DEFAULT_LIMIT, "the maximum number of bytes to read")
	flag.IntVar(&opts.BlockSize, "block-size", DEFAULT_BLOCK_SIZE, "the size of one block in bytes when reading and writing")

	flag.Parse()

	return &opts, nil
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readBytes(options *Options, n int) ([]byte, error) {
	res := make([]byte, n)
	bt, err := io.ReadAtLeast(options.From, res, 1)
	if bt < n {
		return res[:bt], err
	}
	return res, err
}

func writeBytes(options *Options, buff []byte) {
	if _, err := io.WriteString(options.To, string(buff)); err != nil {
		panic(err)
	}
}

func TrimLimit(options *Options, len int) int {
	if options.Limit == -1 {
		return len
	}
	if options.Limit >= len {
		options.Limit -= len
	} else {
		remainder := options.Limit
		options.Limit = 0
		return remainder
	}
	return len
}

func skipOffset(options *Options) {
	if options.Offset != 0 {
		for options.Offset != 0 {
			block := int(math.Min(float64(options.Offset), float64(options.BlockSize)))
			res, err := readBytes(options, block)
			if err != nil {
				panic(err)
			}
			options.Offset -= len(res)
		}
	}
}

func readInit(options *Options) {
	res, err := readBytes(options, options.BlockSize)
	res = res[:TrimLimit(options, len(res))]
	for err != io.EOF && len(res) != 0 {
		writeBytes(options, res)
		if options.Limit == 0 {
			return
		}
		res, err = readBytes(options, options.BlockSize)
		res = res[:TrimLimit(options, len(res))]
	}
}

func main() {
	opts, err := ParseFlags()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not parse flags:", err)
		os.Exit(1)
	}

	optionsInit(opts)
	skipOffset(opts)
	readInit(opts)
}
