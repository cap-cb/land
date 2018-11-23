package main

import (
	"fmt"
	"os"

	"flag"
	"strings"

	"github.com/tzmfreedom/goland/ast"
	"github.com/tzmfreedom/goland/builtin"
	"github.com/tzmfreedom/goland/compiler"
	"github.com/tzmfreedom/goland/interpreter"
	"github.com/tzmfreedom/goland/visitor"
)

var classMap = builtin.NewClassMap()
var preprocessors = []ast.PreProcessor{
	func(src string) string {
		return strings.Replace(src, "// #debugger", "Debugger.debug();", -1)
	},
}

func main() {
	flg := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	f := flg.String("f", "", "file")
	_ = flg.String("d", "", "directory")

	cmd := os.Args[1]

	err := flg.Parse(os.Args[2:])
	if err != nil {
		panic(err.Error())
	}

	t, err := ast.ParseFile(*f, preprocessors...)
	if err != nil {
		handleError(err)
	}
	switch cmd {
	case "format":
		tos(t)
	case "run":
		root, err := convert(t)
		if err != nil {
			handleError(err)
		}
		t, err := register(root)
		err = semanticAnalysis(t)
		if err != nil {
			handleError(err)
		}
		err = run(t)
		if err != nil {
			handleError(err)
		}
	case "check":
		root, err := convert(t)
		if err != nil {
			handleError(err)
		}
		t, err := register(root)
		err = semanticAnalysis(t)
		if err != nil {
			handleError(err)
		}
	}
}

func convert(n ast.Node) (ast.Node, error) {
	return n, nil
}

func check(n ast.Node) error {
	checker := &visitor.SoqlChecker{}
	_, err := n.Accept(checker)
	return err
}

func register(n ast.Node) (*builtin.ClassType, error) {
	register := &compiler.ClassRegisterVisitor{}
	t, err := n.Accept(register)
	if err != nil {
		return nil, err
	}
	classType := t.(*builtin.ClassType)
	classMap.Set(classType.Name, classType)
	return classType, nil
}

func semanticAnalysis(t *builtin.ClassType) error {
	typeChecker := compiler.NewTypeChecker()
	typeChecker.Context.ClassTypes = builtin.PrimitiveClassMap()
	typeChecker.Context.ClassTypes.Set(t.Name, t)
	_, err := typeChecker.VisitClassType(t)
	if len(typeChecker.Errors) != 0 {
		for _, e := range typeChecker.Errors {
			fmt.Fprintf(os.Stderr, "%s\n", e.Message)
		}
	}
	return err
}

func run(n *builtin.ClassType) error {
	interpreter := interpreter.NewInterpreter(builtin.PrimitiveClassMap())
	interpreter.Context.ClassTypes.Set(n.Name, n)
	invoke := &ast.MethodInvocation{
		NameOrExpression: &ast.Name{
			Value: []string{n.Name, "action"},
		},
	}
	_, err := invoke.Accept(interpreter)
	return err
}

func tos(n ast.Node) {
	visitor := &ast.TosVisitor{}
	r, _ := n.Accept(visitor)
	fmt.Println(r)
}

func handleError(err error) {
	fmt.Fprintf(os.Stderr, err.Error())
	os.Exit(1)
}
