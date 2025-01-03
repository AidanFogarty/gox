package main

import (
	"bytes"
	"fmt"
	"go/format"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type FieldDefinition struct {
	Name string
	Type string
}

type TypeDefinition struct {
	Name   string
	Fields []FieldDefinition
}

type AstTemplateData struct {
	PackageName string
	BaseName    string
	Types       []TypeDefinition
}

func defineAstTemplate(outputDir string, baseName string, data AstTemplateData) {
	tmpl, _ := template.New("ast.tmpl").Funcs(template.FuncMap{
		"lower": strings.ToLower,
	}).ParseFiles("tools/templates/ast.tmpl")

	var buffer bytes.Buffer
	err := tmpl.Execute(&buffer, data)
	if err != nil {
		fmt.Println("error: unable to execute template", err)
		os.Exit(1)
	}

	file, _ := os.Create(filepath.Join(outputDir, strings.ToLower(baseName)+".go"))
	defer file.Close()

	content, _ := format.Source(buffer.Bytes())

	file.Write(content)
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: generate_ast <output directory>")
		os.Exit(1)
	} else if len(os.Args) == 2 {
		_, err := filepath.Abs(os.Args[1])
		if err != nil {
			fmt.Println("error: unable to get absolute path")
			os.Exit(1)
		}
	}

	outputDir := os.Args[1]

	data := AstTemplateData{
		PackageName: "gox",
		BaseName:    "Expr",
		Types: []TypeDefinition{
			{
				Name: "Binary",
				Fields: []FieldDefinition{
					{Name: "Left", Type: "Expr"},
					{Name: "Operator", Type: "*Token"},
					{Name: "Right", Type: "Expr"},
				},
			},
			{
				Name: "Grouping",
				Fields: []FieldDefinition{
					{Name: "Expression", Type: "Expr"},
				},
			},
			{
				Name: "Literal",
				Fields: []FieldDefinition{
					{Name: "Value", Type: "interface{}"},
				},
			},
			{
				Name: "Unary",
				Fields: []FieldDefinition{
					{Name: "Operator", Type: "*Token"},
					{Name: "Right", Type: "Expr"},
				},
			},
		},
	}

	defineAstTemplate(outputDir, "Expr", data)
}
