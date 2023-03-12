package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"io"
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
	To         io.Writer
	TempFile   *os.File
	ConvStr    string
	Conv       struct {
		UpperCase  bool
		LowerCase  bool
		TrimSpaces bool
	}
}

func InitOptionConv(options *Options) {
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
				panic(errors.New("invalid -conv arguments"))
			}
		}
	}
	// Если нужно обрезать пробелы, то нужен вспомогательный файл
	if options.Conv.TrimSpaces {
		var err error
		options.TempFile, err = os.Create(TempFileName)
		check(err)
	}
}

func InitOptionInputOutput(options *Options) {
	var err error
	if options.FileInput != "" {
		options.From, err = os.Open(options.FileInput)
		check(err)
	} else {
		options.From = os.Stdin
	}
	if options.FileOutput != "" {
		if _, err := os.Stat(options.FileOutput); errors.Is(err, os.ErrNotExist) {
			options.To, err = os.Create(options.FileOutput)
			if err != nil {
				fmt.Println(err)
				return
			}
		} else {
			panic(errors.New("file already exists"))
		}
	} else {
		options.To = os.Stdout
	}
}

func optionsInit(options *Options) {
	InitOptionInputOutput(options)
	InitOptionConv(options)
	if options.Offset < 0 {
		panic(errors.New("offset must be positiv"))
	}
	if options.Conv.UpperCase && options.Conv.LowerCase {
		panic(errors.New("lower_case and upper_case cannot be paramentary at the same time"))
	}
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

func check(e error) {
	if e != nil {
		panic(e)
	}
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

func skipLSEP(plp *os.File) int {
	//fsymbol := []byte(" ")
	b := make([]byte, 6)
	_, err := plp.Read(b)
	chr, _ := utf8.DecodeRune(b)
	if !unicode.IsSpace(chr) {
		return 6
	}
	if err == io.EOF {
		return 0
	}
	return 0
}

func FindTextBounds() (int, int) {
	pp, _ := os.Open(TempFileName)
	plp := bufio.NewReader(pp)
	flagg := false
	i := 0
	lbound := 0
	rbound := 0
	for {
		rn, inc, err := plp.ReadRune()
		if rn == '\\' {
			i += skipLSEP(pp)
			for j := 0; j <= 4; j++ {
				_, _, err := plp.ReadRune()
				if err == io.EOF {
					break
				}
			}
			continue
		}
		if !unicode.IsSpace(rn) && !flagg {
			flagg = true
			lbound = i
			rbound = i
		} else if !unicode.IsSpace(rn) && flagg && err != io.EOF {
			rbound = i
		}
		if err != nil && !errors.Is(err, io.EOF) {
			return 0, 0
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
	return lbound, rbound + 1

}

func ToLower(slice []byte) []byte {
	return bytes.ToLower(slice)
}

func ToUpper(slice []byte) []byte {
	return bytes.ToUpper(slice)
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
	if options.Conv.TrimSpaces {
		if _, err := io.WriteString(options.TempFile, string(buff)); err != nil {
			panic(err)
		}
	} else {
		buff = ValidateBytesArray(buff)
		if options.Conv.LowerCase && !options.Conv.TrimSpaces {
			buff = ToLower(buff)
		}
		if options.Conv.UpperCase && !options.Conv.TrimSpaces {
			buff = ToUpper(buff)
		}
		if _, err := io.WriteString(options.To, string(buff)); err != nil {
			panic(err)
		}
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

func RWInit(options *Options) {
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

func PostProcessing(options *Options, l int, r int) {
	options.FileInput = TempFileName
	InitOptionInputOutput(options)
	options.Offset = l
	options.Limit = r - l
	options.Conv.TrimSpaces = false
	skipOffset(options)
	RWInit(options)
}

func main() {
	opts, err := ParseFlags()
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "can not parse flags:", err)
		os.Exit(1)
	}

	optionsInit(opts)
	skipOffset(opts)
	RWInit(opts)
	if opts.Conv.TrimSpaces {
		l, r := FindTextBounds()
		PostProcessing(opts, l, r)
	}
	if len(MiniBuff) > 0 {
		if _, err := io.WriteString(opts.To, string(MiniBuff)); err != nil {
			panic(err)
		}
	}
}
