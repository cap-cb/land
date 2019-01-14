package builtin

import "github.com/tzmfreedom/goland/ast"

func init() {
	instanceMethods := NewMethodMap()
	staticMethods := NewMethodMap()
	staticMethods.Set(
		"currentPage",
		[]ast.Node{
			CreateMethod(
				"currentPage",
				[]string{"PageReference"},
				[]ast.Node{},
				func(this interface{}, params []interface{}, extra map[string]interface{}) interface{} {
					return extra["current_page"]
				},
			),
		},
	)

	classType := CreateClass(
		"ApexPages",
		[]*ast.ConstructorDeclaration{},
		instanceMethods,
		staticMethods,
	)

	primitiveClassMap.Set("ApexPages", classType)
}