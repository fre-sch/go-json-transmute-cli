# go-json-transmute-cli

Command line utility to transfom JSON using JSON-Transmute https://github.com/fre-sch/go-libtransmute

## Usage

```
json-transmute [split|single]

split
  -data string
        required, file path to JSON data
  -expr string
        required, file path to JSON expression
single
  -data string
        required, JSON-Path to context data
  -expr string
        required, JSON-Path to expression
  -input string
        required, path to JSON file containing expression and context data
```

### split

Expression JSON and context JSON are in separate files.

```
json-transmute split -expr path/to/expr.json -data path/to/data.json
```

### single

Expression JSON and context JSON are in same files. Use JSON-Path to specify
expression and data.

```
json-transmute single -input path/to/single.json -expr "$.expression" -data "$.data"
```
