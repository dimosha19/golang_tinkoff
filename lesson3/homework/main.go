package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"math"
	"os"
	"strings"
	"unicode"
	"unicode/utf8"
)

const DefaultBlockSize int = 1024
const DefaultOffset int = 0
const DefaultLimit int = -1
const TempFileName string = "temp.txt"

var MiniBuff []byte

type Options struct {
	FileInput  string
	FileOutput string
	Offset     int
	Limit      int
	BlockSize  int
	From       io.Reader
	To         OffsetWriter
	TempFile   *os.File
	ConvStr    string
	Conv       struct {
		UpperCase  bool
		LowerCase  bool
		TrimSpaces bool
	}
}

func InitOptionConv(options *Options, create bool) error {
	var err error
	options.Conv.LowerCase = false
	options.Conv.UpperCase = false
	options.Conv.TrimSpaces = false
	if options.ConvStr != "" {
		s := strings.Split(options.ConvStr, ",")
		for i := 0; i < len(s); i++ {
			switch s[i] {
			case "upper_case":
				options.Conv.UpperCase = true
			case "lower_case":
				options.Conv.LowerCase = true
			case "trim_spaces":
				options.Conv.TrimSpaces = true
			default:
				return errors.New("invalid -conv arguments")
			}
		}
	}
	// Если нужно обрезать пробелы, то нужен вспомогательный файл
	if options.Conv.TrimSpaces && create {
		options.TempFile, err = os.Create(TempFileName)
	}
	return err
}

func InitOptionInputOutput(options *Options) error {
	var err error
	if options.FileInput != "" {
		options.From, err = os.Open(options.FileInput)
		if err != nil {
			return err
		}
	} else {
		options.From = os.Stdin
	}
	if options.FileOutput != "" {
		if _, err := os.Stat(options.FileOutput); errors.Is(err, os.ErrNotExist) {
			outputStream, err2 := os.Create(options.FileOutput)
			if err2 != nil {
				return err2
			}
			options.To = OffsetWriter{outputStream, int64(options.Offset)}
		} else {
			return errors.New("file already exists")
		}
	} else {
		options.To = OffsetWriter{os.Stdout, int64(options.Offset)}
	}
	return err
}

type OffsetWriter struct {
	W io.Writer
	N int64
}

func (p *OffsetWriter) Write(b []byte) (int, error) {
	m := int64(math.Min(float64(len(b)), float64(p.N)))
	p.N -= m
	return io.WriteString(p.W, string(b[m:]))
}

func (p *OffsetWriter) GetN() int64 {
	return p.N
}

func optionsInit(options *Options, create bool) error {
	var err error
	err = InitOptionInputOutput(options)
	if err != nil {
		return err
	}
	err = InitOptionConv(options, create)
	if err != nil {
		return err
	}
	if options.Limit != -1 {
		options.From = io.LimitReader(options.From, int64(options.Limit+options.Offset))
	}
	if options.Offset < 0 {
		return errors.New("offset must be positiv")
	}
	if options.Conv.UpperCase && options.Conv.LowerCase {
		return errors.New("lower_case and upper_case cannot be paramentary at the same time")
	}
	return err
}

func ParseFlags() (*Options, error) {
	var opts Options

	flag.StringVar(&opts.FileInput, "from", "", "file to read. by default - stdin")
	flag.StringVar(&opts.FileOutput, "to", "", "file to write. by default - stdout")
	flag.IntVar(&opts.Offset, "offset", DefaultOffset, "number of bytes to skip when copying")
	flag.IntVar(&opts.Limit, "limit", DefaultLimit, "the maximum number of bytes to read")
	flag.IntVar(&opts.BlockSize, "block-size", DefaultBlockSize, "the size of one block in bytes when reading and writing")
	flag.StringVar(&opts.ConvStr, "conv", "", "UpperCase / LowerCase / TrimSpaces")

	flag.Parse()

	return &opts, nil
}

func ValidateBytesArray(bytes []byte) []byte {
	var res []byte
	s := append(MiniBuff, bytes...)
	for i := 1; i <= len(s); i++ {
		a, _ := utf8.DecodeRune(s[:i])
		if a != utf8.RuneError {
			MiniBuff = s[i:]
			res = s[:i]
			if i == len(s) {
				return res
			}
		}
	}
	MiniBuff = append(MiniBuff, bytes...)
	return res
}

func FindTextBounds() (int, int, error) {
	pp, _ := os.Open(TempFileName)
	plp := bufio.NewReader(pp)
	flagg := false
	i := 0
	lbound := 0
	rbound := 0
	for {
		rn, inc, err := plp.ReadRune()
		if !unicode.IsSpace(rn) && !flagg {
			flagg = true
			lbound = i
			rbound = i
		} else if !unicode.IsSpace(rn) && flagg && err != io.EOF {
			rbound = i
		}
		if err != nil && !errors.Is(err, io.EOF) {
			return 0, 0, nil
		}
		if err == io.EOF {
			break
		}
		i += inc
	}
	defer func(pp *os.File) {
		err := pp.Close()
		if err != nil {
			panic(err)
		}
	}(pp)
	return lbound, rbound + 1, nil

}

func ToLower(slice []byte) []byte {
	return bytes.ToLower(slice)
}

func ToUpper(slice []byte) []byte {
	return bytes.ToUpper(slice)
}

func readBytes(reader io.Reader, n int) ([]byte, error) {
	res := make([]byte, n)
	bt, err := io.ReadAtLeast(reader, res, 1)
	if bt < n {
		return res[:bt], err
	}
	return res, err
}

func ConvertCase(b []byte, lower bool, upper bool) []byte {
	if !lower && !upper {
		return b
	} else if lower {
		return ToLower(b)
	} else {
		return ToUpper(b)
	}
}

func writeBytes(options *Options, buff []byte) error {
	if options.Conv.TrimSpaces {
		if _, err := io.WriteString(&OffsetWriter{options.TempFile, int64(options.Offset)}, string(buff)); err != nil {
			return err
		}
	} else {
		buff = ValidateBytesArray(buff)
		buff = ConvertCase(buff, options.Conv.LowerCase, options.Conv.UpperCase)
		if _, err := io.WriteString(&options.To, string(buff)); err != nil {
			return err
		}
	}
	return nil
}

func RWInit(options *Options) error {
	var err error
	var res []byte
	bflag := false
	for err != io.EOF && len(res) != 0 || !bflag {
		bflag = true
		res, err = readBytes(options.From, options.BlockSize)
		if err != io.EOF || err == nil {
			err = writeBytes(options, res)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return err
}

func PostProcessing(options *Options, l int, r int) error {
	options.FileInput = TempFileName
	options.Offset = l
	options.Limit = r - l
	if err := optionsInit(options, false); err != nil {
		return err
	}
	options.Conv.TrimSpaces = false
	if err := RWInit(options); err != nil {
		return err
	}
	return nil
}

func main() {
	l := log.New(os.Stderr, "", 0)
	opts, err := ParseFlags()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not parse flags:", err)
		os.Exit(1)
	}

	err = optionsInit(opts, true)
	if err != nil {
		l.Println(err)
		panic(err)
	}

	err = RWInit(opts)
	if err != io.EOF {
		l.Println(err)
		return
	}

	if opts.Conv.TrimSpaces {
		lb, rb, err := FindTextBounds()
		if err != nil && err != io.EOF {
			l.Println(err)
		}
		err = PostProcessing(opts, lb, rb)
		if err != nil && err != io.EOF {
			l.Println(err)
		}
	}
	if len(MiniBuff) > 0 {
		if _, err := io.WriteString(&opts.To, string(MiniBuff)); err != nil {
			panic(err)
		}
	}
	if opts.To.N != 0 {
		panic(errors.New("offset > input"))
	}
}
