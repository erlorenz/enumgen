package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type EnumData struct {
	SourceName  string
	Type        string
	Value       string
	Description string
	Items       map[string]string
}

//go:embed templates
var templates embed.FS

func main() {
	pwd, _ := os.Getwd()
	args := os.Args[1:]

	cfg := struct {
		gopackage  string // The GOPACKAGE env var
		gofile     string // The GOFILE envvar
		pwd        string // The pwd when running
		outFile    string // The destination file (if not split)
		sourceFile string // The path of the source file
		splitFiles bool   // Split output into a file per map
	}{
		gopackage:  os.Getenv("GOPACKAGE"),
		gofile:     os.Getenv("GOFILE"),
		pwd:        pwd,
		outFile:    "",
		sourceFile: "",
	}
	cfg.sourceFile = filepath.Join(cfg.pwd, cfg.gofile)

	f := flag.NewFlagSet("genenum", flag.ExitOnError)
	outVar := f.String("o", "enums", "file to output")
	f.BoolVar(&cfg.splitFiles, "split", false, "split files")
	f.Parse(args)

	cfg.outFile = filepath.Join(cfg.pwd, *outVar+"_gen.go")

	// ------------------------------------------------------------------
	// Parse Template

	templ, err := template.ParseFS(templates, "*/**.gotmpl")
	if err != nil {
		fmt.Printf("enumgen: Error parsing template: %s\n", err.Error())
		os.Exit(1)
	}

	//----------------------------------------------------------------
	// Parsing maps

	fmt.Printf("enumgen: Reading file %s...\n", cfg.sourceFile)

	// Parse the Go file into an AST
	node, err := parser.ParseFile(token.NewFileSet(), cfg.sourceFile, nil, parser.ParseComments)
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

	type TemplateData struct {
		Package string
		Enums   []EnumData
	}

	data := TemplateData{cfg.gopackage, MapsToEnumData(maps)}

	// ------------------------------------------------------------------
	// Write single

	if !cfg.splitFiles || len(data.Enums) < 2 {

		fmt.Printf("enumgen: Writing to single file %s.\n", cfg.outFile)

		var buf bytes.Buffer

		err = templ.ExecuteTemplate(&buf, "enums.gotmpl", data)
		if err != nil {
			fmt.Printf("enumgen: Error executing template: %s.\n", err.Error())
			os.Exit(1)
		}

		err = os.WriteFile(cfg.outFile, buf.Bytes(), 0777)
		if err != nil {
			fmt.Printf("enumgen: Error generating file %s: %s.\n", cfg.outFile, err.Error())
			os.Exit(1)
		}

		fmt.Printf("enumgen: Generated file %s.", filepath.Join(cfg.gopackage, cfg.outFile))

		return
	}

	// ------------------------------------------------------------------
	// Write multiple
	fmt.Println("enumgen: Writing to multiple files...")

	for _, enum := range data.Enums {
		outbase := strings.ToLower(enum.Type) + "_gen.go"
		outfile := filepath.Join(cfg.pwd, outbase)

		// shadow data
		data := TemplateData{
			Package: data.Package,
			Enums:   []EnumData{enum},
		}

		var buf bytes.Buffer

		err = templ.ExecuteTemplate(&buf, "enums.gotmpl", data)
		if err != nil {
			fmt.Printf("enumgen: Error executing template: %s.\n", err.Error())
			os.Exit(1)
		}

		fmt.Printf("enumgen: Writing file %s.\n", outfile)

		err = os.WriteFile(outfile, buf.Bytes(), 0777)
		if err != nil {
			fmt.Printf("enumgen: Error generating file %s: %s.\n", cfg.outFile, err.Error())
			os.Exit(1)
		}

		fmt.Printf("enumgen: Generated file %s.", cfg.outFile)
	}

}
