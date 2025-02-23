package main

import (
	"bytes"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"text/template"
)

type EnumData struct {
	SourceName  string
	Type        string
	Value       string
	Description string
	Items       map[string]string
}

func main() {
	pwd, _ := os.Getwd()
	args := os.Args[1:]

	cfg := struct {
		gopackage string
		gofile    string
		pwd       string
		outfile   string
		path      string
	}{
		gopackage: os.Getenv("GOPACKAGE"),
		gofile:    os.Getenv("GOFILE"),
		pwd:       pwd,
		outfile:   "",
		path:      "",
	}
	cfg.path = filepath.Join(cfg.pwd, cfg.gofile)

	f := flag.NewFlagSet("genenum", flag.ExitOnError)

	outVar := f.String("o", "enums.go", "file to output")
	cfg.outfile = filepath.Join(cfg.pwd, *outVar)

	f.Parse(args)

	//----------------------------------------------------------------
	// Parsing maps
	fmt.Printf("enumgen: Reading file %s...\n", cfg.path)

	// Parse the Go file into an AST
	node, err := parser.ParseFile(token.NewFileSet(), cfg.path, nil, parser.ParseComments)
	if err != nil {
		log.Fatalf("enumgen: Failed to parse file: %v", err)
	}

	// Find and parse all map variables ending in "Values"
	maps, err := ParseMapValues(node)
	if err != nil {
		log.Fatalf("enumgen: Failed to parse maps: %v", err)
	}

	// Process each found map
	for _, result := range maps {
		fmt.Printf("enumgen: Found map: %s\n", result.Name)
		for k, v := range result.Values {
			fmt.Printf("  %s: %s\n", k, v)
		}
	}

	// ------------------------------------------------------------------
	// Data

	var buf bytes.Buffer

	type TemplateData struct {
		Package string
		Enums   []EnumData
	}

	data := TemplateData{cfg.gopackage, MapsToEnumData(maps)}

	// ------------------------------------------------------------------
	// Template

	templ, err := template.ParseGlob("../*.gotmpl")
	if err != nil {
		fmt.Printf("enumgen: Error parsing template: %s\n", err.Error())
		os.Exit(1)
	}

	err = templ.ExecuteTemplate(&buf, "enums.gotmpl", data)
	if err != nil {
		fmt.Printf("enumgen: Error executing template: %s.\n", err.Error())
		os.Exit(1)
	}

	// ------------------------------------------------------------------
	// Write

	err = os.WriteFile(cfg.path, buf.Bytes(), 0777)
	if err != nil {
		fmt.Printf("enumgen: Error generating file %s: %s.\n", filepath.Join(cfg.pwd, cfg.outfile), err.Error())
		os.Exit(1)
	}

	fmt.Printf("enumgen: Generated file %s.", filepath.Join(cfg.gopackage, cfg.outfile))
}
