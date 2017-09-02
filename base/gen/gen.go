package gen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
)

var basePath = filepath.Join(os.Getenv("GOPATH"),
	"src/github.com/ng-vu/go-grpc-sample")

// Method ...
type Method struct {
	Name       string
	InputType  string
	OutputType string
}

// ParseServiceFile ...
func ParseServiceFile(inputPath, interfaceName string) map[string]Method {
	p := newParser()
	return p.parse(inputPath, interfaceName)
}

type sortMethodsType []Method

func (a sortMethodsType) Len() int           { return len(a) }
func (a sortMethodsType) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a sortMethodsType) Less(i, j int) bool { return a[i].Name < a[j].Name }

// SortMethods ...
func SortMethods(methods map[string]Method) []Method {
	sortMethods := make(sortMethodsType, 0, len(methods))
	for _, m := range methods {
		sortMethods = append(sortMethods, m)
	}
	sort.Sort(sortMethods)
	return sortMethods
}

// WriteFile ...
func WriteFile(outputPath string, data []byte) {
	absPath := filepath.Join(basePath, outputPath)
	err := ioutil.WriteFile(absPath, data, os.ModePerm)
	if err != nil {
		fmt.Printf("Unable to write file `%v`.\n  Error: %v\n", absPath, err)
		os.Exit(1)
	}
	fmt.Println("Generated file:", absPath)

	out, err := exec.Command("gofmt", "-w", absPath).Output()
	if err != nil {
		fmt.Printf("Unable to run `gofmt -w %v`.\n  Error: %v\n", absPath, err)
		os.Exit(1)
	}
	fmt.Print(string(out))
}

type parserStruct struct {
	interfaceName string

	fset *token.FileSet
}

func newParser() *parserStruct {
	return &parserStruct{
		fset: token.NewFileSet(),
	}
}

func (p *parserStruct) parse(inputPath, interfaceName string) map[string]Method {
	absPath := filepath.Join(basePath, inputPath)
	f, err := parser.ParseFile(p.fset, absPath, nil, 0)
	if err != nil {
		fmt.Printf("Unable to parse file `%v`.\n  Error: %v\n", absPath, err)
		os.Exit(1)
		return nil
	}

	p.interfaceName = interfaceName
	return p.extract(f)
}

func (p *parserStruct) extract(f *ast.File) (methods map[string]Method) {
	inspectFunc := func(node ast.Node) bool {
		switch node := node.(type) {
		case *ast.TypeSpec:
			if node.Name.Name != p.interfaceName {
				return false
			}

			switch typ := node.Type.(type) {
			case *ast.InterfaceType:
				methods = p.extractMethods(typ)
				return false

			default:
				fmt.Printf("Error: %v must be an interface\n", p.interfaceName)
				os.Exit(1)
			}
		}
		return true
	}

	ast.Inspect(f, inspectFunc)
	if methods == nil {
		fmt.Printf("Error: Unable to find type %v\n", p.interfaceName)
		os.Exit(1)
	}
	return methods
}

func (p *parserStruct) extractMethods(typ *ast.InterfaceType) map[string]Method {
	methods := make(map[string]Method)
	for _, m := range typ.Methods.List {
		name := m.Names[0].Name

		fnTyp := m.Type.(*ast.FuncType)
		params := fnTyp.Params.List
		if len(params) != 2 {
			fmt.Printf("Error: Method %v must have exactly 2 arguments\n", name)
			os.Exit(1)
		}
		inputName := p.extractTypeName(name, params[1].Type)

		results := fnTyp.Results.List
		if len(results) != 2 {
			fmt.Printf("Error: Method %v must have exactly 2 results\n", name)
			os.Exit(1)
		}
		outputName := p.extractTypeName(name, results[0].Type)

		methods[name] = Method{
			Name:       name,
			InputType:  inputName,
			OutputType: outputName,
		}
	}
	return methods
}

func (p *parserStruct) extractTypeName(method string, typ ast.Expr) string {
	s := ""
	if t, ok := typ.(*ast.StarExpr); ok {
		typ = t.X
		s = "*"
	}

	switch typ := typ.(type) {
	case *ast.Ident:
		return s + "pb." + typ.Name

	default:
		fmt.Printf("Error: Method %v must have input and output in format *Foo\n", method)
		os.Exit(1)
		return ""
	}
}
