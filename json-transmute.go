package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	transmute "github.com/fre-sch/go-libtransmute"
)

func loadJSONFile(path string, target interface{}) (err error) {
	var data []byte
	if data, err = ioutil.ReadFile(path); err != nil {
		fmt.Fprintf(os.Stderr, "could not read expr file: %s\n", path)
		return
	}
	if err = json.Unmarshal(data, target); err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse JSON: %s\n", path)
		fmt.Fprintf(os.Stderr, "%v", err)
		return
	}
	return
}

type inputArgs struct {
	exprFilePath   *string
	dataFilePath   *string
	inputFilePath  *string
	exprStrPath    *string
	contextStrPath *string
}

func main() {
	splitCmd := flag.NewFlagSet("split", flag.ExitOnError)
	singleCmd := flag.NewFlagSet("single", flag.ExitOnError)

	args := inputArgs{
		exprFilePath: splitCmd.String("expr", "", "required, file path to JSON expression"),
		dataFilePath: splitCmd.String("data", "", "required, file path to JSON data"),

		inputFilePath:  singleCmd.String("input", "", "required, path to JSON file containing expression and context data"),
		exprStrPath:    singleCmd.String("expr", "", "required, JSON-Path to expression"),
		contextStrPath: singleCmd.String("data", "", "required, JSON-Path to context data"),
	}

	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "split\n")
		splitCmd.PrintDefaults()
		fmt.Fprintf(os.Stderr, "single\n")
		singleCmd.PrintDefaults()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "split":
		splitCmd.Parse(os.Args[2:])
		handleSplitCommand(args)
	case "single":
		singleCmd.Parse(os.Args[2:])
		handleSingleCommand(args)
	}

}

func handleSplitCommand(args inputArgs) {
	if *(args.exprFilePath) == "" || *(args.dataFilePath) == "" {
		os.Exit(1)
	}

	var expr interface{}
	var context interface{}

	if err := loadJSONFile(*args.exprFilePath, &expr); err != nil {
		os.Exit(1)
	}
	if err := loadJSONFile(*args.dataFilePath, &context); err != nil {
		os.Exit(1)
	}

	var result interface{}
	var encodeResult []byte
	var err error

	if result, err = transmute.Transmute(expr, context); err != nil {
		fmt.Fprintf(os.Stderr, "error transmuting: %#+v", err)
		os.Exit(1)
	}

	if encodeResult, err = json.MarshalIndent(result, "", "  "); err != nil {
		fmt.Fprintf(os.Stderr, "error encoding to JSON: %#+v", err)
		os.Exit(1)
	}

	fmt.Println(string(encodeResult))
}

func handleSingleCommand(args inputArgs) {
	var jsonData interface{}
	var contextData interface{}
	var exprData interface{}
	var result interface{}
	var encodeResult []byte
	var err error

	if err = loadJSONFile(*args.inputFilePath, &jsonData); err != nil {
		fmt.Fprintf(os.Stderr, "fai")
		os.Exit(1)
	}

	contextData, err = transmute.Transmute(*args.contextStrPath, jsonData)
	if err != nil {
		os.Exit(1)
	}
	exprData, err = transmute.Transmute(*args.exprStrPath, jsonData)
	if err != nil {
		os.Exit(1)
	}

	if result, err = transmute.Transmute(exprData, contextData); err != nil {
		fmt.Fprintf(os.Stderr, "error transmuting: %#+v", err)
		os.Exit(1)
	}

	if encodeResult, err = json.MarshalIndent(result, "", "  "); err != nil {
		fmt.Fprintf(os.Stderr, "error encoding to JSON: %#+v", err)
		os.Exit(1)
	}

	fmt.Println(string(encodeResult))
}
