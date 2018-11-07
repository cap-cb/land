package main

import (
	"github.com/antlr/antlr4/runtime/Go/antlr"
	"github.com/tzmfreedom/goland/ast"
	"github.com/tzmfreedom/goland/parser"
	"github.com/golang/go/src/cmd/go/testdata/testinternal3"
)

type AstBuilder struct {
	*parser.BaseapexVisitor
	CurrentFile string
}

func (v *AstBuilder) VisitCompilationUnit(ctx *parser.CompilationUnitContext) interface{} {
	return ctx.TypeDeclaration().Accept(v)
}

func (v *AstBuilder) VisitTypeDeclaration(ctx *parser.TypeDeclarationContext) interface{} {
	classOrInterfaceModifiers := ctx.AllClassOrInterfaceModifier()
	modifiers := []ast.Modifier{}
	annotations := []ast.Annotation{}
	for _, classOrInterfaceModifier := range classOrInterfaceModifiers {
		r := classOrInterfaceModifier.Accept(v)
		m, ok := r.(ast.Modifier)
		if ok {
			modifiers = append(modifiers, m)
		}
		a, ok := r.(ast.Annotation)
		if ok {
			annotations = append(annotations, a)
		}
	}

	if ctx.ClassDeclaration() != nil {
		cd := ctx.ClassDeclaration().Accept(v)
		classDeclaration, _ := cd.(ast.ClassDeclaration)
		classDeclaration.Modifiers = modifiers
		classDeclaration.Annotations = annotations
		return classDeclaration
	}
	return nil
}

func (v *AstBuilder) VisitTriggerDeclaration(ctx *parser.TriggerDeclarationContext) interface{} {
	timings := ctx.TriggerTimings().Accept(v)

	name := ctx.ApexIdentifier(0).GetText()
	object := ctx.ApexIdentifier(1).GetText()
	block := ctx.Block().Accept(v)
	return ast.Trigger{
		Name:           name,
		TriggerTimings: timings.([]ast.TriggerTiming),
		Object:         object,
		Statements:     block.([]ast.Node),
		Position:       v.newPosition(ctx),
	}
}

func (v *AstBuilder) VisitTriggerTimings(ctx *parser.TriggerTimingsContext) interface{} {
	allTimings := ctx.AllTriggerTiming()
	timings := make([]ast.Node, len(allTimings))
	for i, timing := range allTimings {
		timings[i] = timing.Accept(v).(ast.Node)
	}
	return timings
}

func (v *AstBuilder) VisitTriggerTiming(ctx *parser.TriggerTimingContext) interface{} {
	return ast.TriggerTiming{
		Timing:   ctx.GetTiming().GetText(),
		Dml:      ctx.GetDml().GetText(),
		Position: v.newPosition(ctx),
	}
}

func (v *AstBuilder) VisitModifier(ctx *parser.ModifierContext) interface{} {
	m := ctx.ClassOrInterfaceModifier()
	if m != nil {
		return m.Accept(v)
	}
	return ast.Modifier{
		Name:     ctx.GetText(),
		Position: v.newPosition(ctx),
	}
}

func (v *AstBuilder) VisitClassOrInterfaceModifier(ctx *parser.ClassOrInterfaceModifierContext) interface{} {
	annotation := ctx.Annotation()
	if annotation != nil {
		return ctx.Annotation().Accept(v)
	}
	return ast.Modifier{
		Name:     ctx.GetText(),
		Position: v.newPosition(ctx),
	}
}

func (v *AstBuilder) VisitVariableModifier(ctx *parser.VariableModifierContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitClassDeclaration(ctx *parser.ClassDeclarationContext) interface{} {
	return ast.ClassDeclaration{
		Name:     ctx.ApexIdentifier().GetText(),
		Position: v.newPosition(ctx),
	}
}

func (v *AstBuilder) VisitEnumDeclaration(ctx *parser.EnumDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitEnumConstants(ctx *parser.EnumConstantsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitEnumConstant(ctx *parser.EnumConstantContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitEnumBodyDeclarations(ctx *parser.EnumBodyDeclarationsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitInterfaceDeclaration(ctx *parser.InterfaceDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitTypeList(ctx *parser.TypeListContext) interface{} {
	apexTypes := ctx.AllApexType()
	types := make([]ast.Node, len(apexTypes))
	for i, t := range apexTypes {
		types[i] = t.Accept(v).(ast.Node)
	}
	return types
}

func (v *AstBuilder) VisitClassBody(ctx *parser.ClassBodyContext) interface{} {
	bodyDeclarations := ctx.AllClassBodyDeclaration()
	declarations := make([]ast.Node, len(bodyDeclarations))
	for i, d := range bodyDeclarations {
		declarations[i] = d.Accept(v).(ast.Node)
	}
	return declarations
}

func (v *AstBuilder) VisitInterfaceBody(ctx *parser.InterfaceBodyContext) interface{} {
	bodyDeclarations := ctx.AllInterfaceBodyDeclaration()
	declarations := make([]ast.Node, len(bodyDeclarations))
	for i, d := range bodyDeclarations {
		declarations[i] = d.Accept(v).(ast.Node)
	}
	return declarations
}

func (v *AstBuilder) VisitClassBodyDeclaration(ctx *parser.ClassBodyDeclarationContext) interface{} {
	memberDeclaration := ctx.MemberDeclaration()
	if memberDeclaration != nil {
		declaration := memberDeclaration.Accept(v)

		modifiers := ctx.AllModifier()
		declarationModifiers := make([]ast.Modifier, len(modifiers))
		for i, m := range modifiers {
			declarationModifiers[i] = m.Accept(v).(ast.Modifier)
		}
		switch decl := declaration.(type) {
		case *ast.MethodDeclaration:
			decl.Modifiers = declarationModifiers
			return decl
		case *ast.FieldDeclaration:
			decl.Modifiers = declarationModifiers
			return decl
		case *ast.ConstructorDeclaration:
			decl.Modifiers = declarationModifiers
			return decl
		case *ast.InterfaceDeclaration:
			decl.Modifiers = declarationModifiers
			return decl
		case *ast.ClassDeclaration:
			decl.Modifiers = declarationModifiers
			return decl
		//case *ast.EnumDeclaration:
		//	decl.Modifiers = declarationModifiers
		//	return decl
		case *ast.PropertyDeclaration:
			decl.Modifiers = declarationModifiers
			return decl
		}
	}
	return nil
}

func (v *AstBuilder) VisitMemberDeclaration(ctx *parser.MemberDeclarationContext) interface{} {
	return v.VisitChildren(ctx).([]interface{})[0]
}

func (v *AstBuilder) VisitMethodDeclaration(ctx *parser.MethodDeclarationContext) interface{} {
	methodName := ctx.ApexIdentifier().GetText()
	var returnType ast.Node
	if ctx.ApexType() != nil {
		returnType = ctx.ApexType().Accept(v).(ast.Node)
	} else {
		returnType = ast.VoidType
	}
	parameters := ctx.FormalParameters().Accept(v).([]ast.Parameter)
	var throws []ast.Node
	if ctx.QualifiedNameList() != nil {
		throws = ctx.QualifiedNameList().Accept(v).([]ast.Node)
	} else {
		throws = []ast.Node{}
	}
	var statements []ast.Node
	if ctx.MethodBody() != nil {
		statements = ctx.MethodBody().Accept(v).([]ast.Node)
	} else {
		statements = []ast.Node{}
	}
	return ast.MethodDeclaration{
		Name:       methodName,
		ReturnType: returnType,
		Parameters: parameters,
		Throws:     throws,
		Statements: statements,
	}
}

func (v *AstBuilder) VisitConstructorDeclaration(ctx *parser.ConstructorDeclarationContext) interface{} {
	parameters := ctx.FormalParameters().Accept(v).([]ast.Parameter)
	var throws []ast.Node
	if q := ctx.QualifiedNameList(); q != nil {
		throws = q.Accept(v).([]ast.Node)
	} else {
		throws = []ast.Node{}
	}
	body := ctx.ConstructorBody().Accept(v).([]ast.Node)
	return &ast.ConstructorDeclaration{
		Parameters: parameters,
		Throws:     throws,
		Statements: body,
		Position:   v.newPosition(ctx),
	}
}

func (v *AstBuilder) VisitFieldDeclaration(ctx *parser.FieldDeclarationContext) interface{} {
	t := ctx.ApexType().Accept(v)
	d := ctx.VariableDeclarators().Accept(v).([]ast.Node)
	return ast.FieldDeclaration{
		Type:        t,
		Declarators: d,
	}
}

func (v *AstBuilder) VisitPropertyDeclaration(ctx *parser.PropertyDeclarationContext) interface{} {
	t := ctx.ApexType().Accept(v).(ast.Type)
	d := ctx.VariableDeclaratorId().Accept(v).(string)
	b := ctx.PropertyBodyDeclaration().Accept(v).(ast.Node)
	return ast.PropertyDeclaration{
		Type:          t,
		Identifier:    d,
		GetterSetters: b,
	}
}

func (v *AstBuilder) VisitPropertyBodyDeclaration(ctx *parser.PropertyBodyDeclarationContext) interface{} {
	blocks := ctx.AllPropertyBlock()
	declarations := make([]ast.Block, len(blocks))
	for i, b := range blocks {
		declarations[i] = b.Accept(v).(ast.Block)
	}
	return declarations
}

func (v *AstBuilder) VisitInterfaceBodyDeclaration(ctx *parser.InterfaceBodyDeclarationContext) interface{} {
	d := ctx.InterfaceMemberDeclaration().Accept(v).(ast.Interface)
	modifiers := ctx.AllModifier()
	d.Modifiers = make([]ast.Modifier, len(modifiers)+1)
	for i, m := range modifiers {
		d.Modifiers[i] = m.Accept(v).(ast.Modifier)
	}
	d.Modifiers[len(modifiers)] = ast.Modifier{
		Name:     "public",
		Position: v.newPosition(ctx),
	}
	return d
}

func (v *AstBuilder) VisitInterfaceMemberDeclaration(ctx *parser.InterfaceMemberDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitConstDeclaration(ctx *parser.ConstDeclarationContext) interface{} {
	_ = ctx.ApexType().Accept(v)
	_ = ctx.AllConstantDeclarator()

	// TODO: implement
	return nil
}

func (v *AstBuilder) VisitConstantDeclarator(ctx *parser.ConstantDeclaratorContext) interface{} {
	_ = ctx.ApexIdentifier().Accept(v)
	_ = ctx.VariableInitializer().Accept(v)

	// TODO: implement
	return nil
}

func (v *AstBuilder) VisitInterfaceMethodDeclaration(ctx *parser.InterfaceMethodDeclarationContext) interface{} {
	decl := ast.MethodDeclaration{Position: v.newPosition(ctx)}
	decl.Name = ctx.ApexIdentifier().Accept(v).(string)

	if t := ctx.ApexType(); t != nil {
		decl.ReturnType = t.Accept(v).(ast.Type)
	} else {
		// TODO: implement void
	}
	decl.Parameters = ctx.FormalParameters().Accept(v).([]ast.Parameter)
	if q := ctx.QualifiedNameList(); q != nil {
		decl.Throws = q.Accept(v).([]ast.Node)
	} else {
		decl.Throws = []ast.Node{}
	}
	return decl
}

func (v *AstBuilder) VisitVariableDeclarators(ctx *parser.VariableDeclaratorsContext) interface{} {
	variableDeclarators := ctx.AllVariableDeclarator()
	declarators := make([]ast.VariableDeclarator, len(variableDeclarators))
	for i, d := range variableDeclarators {
		declarators[i] = d.Accept(v).(ast.VariableDeclarator)
	}
	return declarators
}

func (v *AstBuilder) VisitVariableDeclarator(ctx *parser.VariableDeclaratorContext) interface{} {
	decl := ast.VariableDeclarator{Position: v.newPosition(ctx)}
	if init := ctx.VariableInitializer(); init != nil {
		decl.Expression = init.Accept(v)
	} else {
		// TODO: implement NULL
	}
	return decl
}

func (v *AstBuilder) VisitVariableDeclaratorId(ctx *parser.VariableDeclaratorIdContext) interface{} {
	return ctx.ApexIdentifier().GetText()
}

func (v *AstBuilder) VisitVariableInitializer(ctx *parser.VariableInitializerContext) interface{} {
	if init := ctx.ArrayInitializer(); init != nil {
		return init.Accept(v)
	}
	return ctx.Expression().Accept(v)
}

func (v *AstBuilder) VisitArrayInitializer(ctx *parser.ArrayInitializerContext) interface{} {
	if inits := ctx.AllVariableInitializer(); len(inits) != 0 {
		initializers := make([]ast.Node, len(inits))
		for i, init := range inits {
			initializers[i] = init.Accept(v)
		}
		return initializers
	}
	return nil
}

func (v *AstBuilder) VisitEnumConstantName(ctx *parser.EnumConstantNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitApexType(ctx *parser.ApexTypeContext) interface{} {
	if interfaceType := ctx.ClassOrInterfaceType(); interfaceType != nil {
		t := interfaceType.Accept(v).(ast.Type)
		// TODO: implement Array
		return t
	} else if primitiveType := ctx.PrimitiveType(); primitiveType != nil {
		t := primitiveType.Accept(v).(ast.Type)
		// TODO: implement Array
		return t
	}
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitTypedArray(ctx *parser.TypedArrayContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitClassOrInterfaceType(ctx *parser.ClassOrInterfaceTypeContext) interface{} {
	if ident := ctx.AllTypeIdentifier(); len(ident) != 0 {
		t := ast.Type{Position: v.newPosition(ctx)}
		// TODO: implement name
		return t
	}
	t := ast.Type{ Position: v.newPosition(ctx) }
	t.Name = ctx.SET().GetText()
	arguments := ctx.AllTypeArguments()
	t.Parameters = make([]ast.Node, len(arguments))
	for i, argument := range arguments {
		t.Parameters[i] = argument.Accept(v).(ast.Node)
	}
	return t
}

func (v *AstBuilder) VisitPrimitiveType(ctx *parser.PrimitiveTypeContext) interface{} {
	return ast.Type{
		Name: ctx.GetText(),
		Position: v.newPosition(ctx),
	}
}

func (v *AstBuilder) VisitTypeArguments(ctx *parser.TypeArgumentsContext) interface{} {
	arguments := ctx.AllTypeArgument()
	typeArguments := make([]ast.Node, len(arguments))
	for i, a := range arguments {
		typeArguments[i] = a.Accept(v).(ast.Node)
	}
	return typeArguments
}

func (v *AstBuilder) VisitTypeArgument(ctx *parser.TypeArgumentContext) interface{} {
	return ctx.ApexType().Accept(v)
}

func (v *AstBuilder) VisitQualifiedNameList(ctx *parser.QualifiedNameListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitFormalParameters(ctx *parser.FormalParametersContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitFormalParameterList(ctx *parser.FormalParameterListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitFormalParameter(ctx *parser.FormalParameterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitLastFormalParameter(ctx *parser.LastFormalParameterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitMethodBody(ctx *parser.MethodBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitConstructorBody(ctx *parser.ConstructorBodyContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitQualifiedName(ctx *parser.QualifiedNameContext) interface{} {
	allIdentifiers := ctx.AllApexIdentifier()
	identifiers := make([]string, len(allIdentifiers))
	for i, identifier := range allIdentifiers {
		ident := identifier.Accept(v)
		identifiers[i], _ = ident.(string)
	}
	return ast.Name{
		Value:    identifiers,
		Position: v.newPosition(ctx),
	}
}

func (v *AstBuilder) VisitLiteral(ctx *parser.LiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitAnnotation(ctx *parser.AnnotationContext) interface{} {
	name := ctx.AnnotationName().Accept(v)
	annotation := ast.Annotation{}
	annotation.Name, _ = name.(ast.Name)
	annotation.Position = v.newPosition(ctx)
	return annotation
}

func (v *AstBuilder) VisitAnnotationName(ctx *parser.AnnotationNameContext) interface{} {
	return ctx.QualifiedName().Accept(v)
}

func (v *AstBuilder) VisitElementValuePairs(ctx *parser.ElementValuePairsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitElementValuePair(ctx *parser.ElementValuePairContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitElementValue(ctx *parser.ElementValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitElementValueArrayInitializer(ctx *parser.ElementValueArrayInitializerContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitBlock(ctx *parser.BlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitBlockStatement(ctx *parser.BlockStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitLocalVariableDeclarationStatement(ctx *parser.LocalVariableDeclarationStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitLocalVariableDeclaration(ctx *parser.LocalVariableDeclarationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitStatement(ctx *parser.StatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitPropertyBlock(ctx *parser.PropertyBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitGetter(ctx *parser.GetterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitSetter(ctx *parser.SetterContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitCatchClause(ctx *parser.CatchClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitCatchType(ctx *parser.CatchTypeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitFinallyBlock(ctx *parser.FinallyBlockContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitWhenStatements(ctx *parser.WhenStatementsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitWhenStatement(ctx *parser.WhenStatementContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitWhenExpression(ctx *parser.WhenExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitForControl(ctx *parser.ForControlContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitForInit(ctx *parser.ForInitContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitEnhancedForControl(ctx *parser.EnhancedForControlContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitForUpdate(ctx *parser.ForUpdateContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitParExpression(ctx *parser.ParExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitExpressionList(ctx *parser.ExpressionListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitStatementExpression(ctx *parser.StatementExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitConstantExpression(ctx *parser.ConstantExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitApexDbExpressionShort(ctx *parser.ApexDbExpressionShortContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitApexDbExpression(ctx *parser.ApexDbExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitTernalyExpression(ctx *parser.TernalyExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitPreUnaryExpression(ctx *parser.PreUnaryExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitArrayAccess(ctx *parser.ArrayAccessContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitPostUnaryExpression(ctx *parser.PostUnaryExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitPrimaryExpression(ctx *parser.PrimaryExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitOpExpression(ctx *parser.OpExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitNewExpression(ctx *parser.NewObjectExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitUnaryExpression(ctx *parser.UnaryExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitMethodInvocation(ctx *parser.MethodInvocationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitCastExpression(ctx *parser.CastExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitShiftExpression(ctx *parser.ShiftExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitFieldAccess(ctx *parser.FieldAccessContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitPrimary(ctx *parser.PrimaryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitCreator(ctx *parser.CreatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitCreatedName(ctx *parser.CreatedNameContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitInnerCreator(ctx *parser.InnerCreatorContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitArrayCreatorRest(ctx *parser.ArrayCreatorRestContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitMapCreatorRest(ctx *parser.MapCreatorRestContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitSetCreatorRest(ctx *parser.SetCreatorRestContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitClassCreatorRest(ctx *parser.ClassCreatorRestContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitExplicitGenericInvocation(ctx *parser.ExplicitGenericInvocationContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitNonWildcardTypeArguments(ctx *parser.NonWildcardTypeArgumentsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitTypeArgumentsOrDiamond(ctx *parser.TypeArgumentsOrDiamondContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitNonWildcardTypeArgumentsOrDiamond(ctx *parser.NonWildcardTypeArgumentsOrDiamondContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitSuperSuffix(ctx *parser.SuperSuffixContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitExplicitGenericInvocationSuffix(ctx *parser.ExplicitGenericInvocationSuffixContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitArguments(ctx *parser.ArgumentsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitSoqlLiteral(ctx *parser.SoqlLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitQuery(ctx *parser.QueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitSelectClause(ctx *parser.SelectClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitFieldList(ctx *parser.FieldListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitSelectField(ctx *parser.SelectFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitFromClause(ctx *parser.FromClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitFilterScope(ctx *parser.FilterScopeContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitSoqlFieldReference(ctx *parser.SoqlFieldReferenceContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitSoqlFunctionCall(ctx *parser.SoqlFunctionCallContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitSubquery(ctx *parser.SubqueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitWhereClause(ctx *parser.WhereClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitWhereFields(ctx *parser.WhereFieldsContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitWhereField(ctx *parser.WhereFieldContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitLimitClause(ctx *parser.LimitClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitOrderClause(ctx *parser.OrderClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitBindVariable(ctx *parser.BindVariableContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitSoqlValue(ctx *parser.SoqlValueContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitWithClause(ctx *parser.WithClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitSoqlFilteringExpression(ctx *parser.SoqlFilteringExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitGroupClause(ctx *parser.GroupClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitFieldGroupList(ctx *parser.FieldGroupListContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitHavingConditionExpression(ctx *parser.HavingConditionExpressionContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitOffsetClause(ctx *parser.OffsetClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitViewClause(ctx *parser.ViewClauseContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitSoslLiteral(ctx *parser.SoslLiteralContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitSoslQuery(ctx *parser.SoslQueryContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitSoslReturningObject(ctx *parser.SoslReturningObjectContext) interface{} {
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitApexIdentifier(ctx *parser.ApexIdentifierContext) interface{} {
	ident := ctx.Identifier()
	if ident != nil {
		return ident.GetText()
	}
	return v.VisitChildren(ctx)
}

func (v *AstBuilder) VisitTypeIdentifier(ctx *parser.TypeIdentifierContext) interface{} {
	return v.VisitChildren(ctx)
}

type PositionContext interface {
	GetStart() antlr.Token
}

func (v *AstBuilder) newPosition(ctx PositionContext) *ast.Position {
	return &ast.Position{
		FileName: v.CurrentFile,
		Column:   ctx.GetStart().GetColumn(),
		Line:     ctx.GetStart().GetLine(),
	}
}
