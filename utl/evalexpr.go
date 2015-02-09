// Copyright 2012 Dorival de Moraes Pedroso. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utl

import (
	"go/ast"
	"go/parser"
	"go/token"
	"math"
	"strconv"
)

// some constants and mathematical functions
var (
	myvars = map[string]float64{
		"pi": math.Pi,
	}
	myfuncs = map[string]interface{}{
		"sin": math.Sin,
		"cos": math.Cos,
		"tan": math.Tan,
		"pow": math.Pow,
		"exp": math.Exp,
		"log": math.Log,
	}
)

// ParseE will parse an expression.
//  Example: expr = "((1+a)/(1-b))*2*sin(pi)"
func ParseE(expr string) (ast.Expr, error) {
	return parser.ParseExpr(expr)
}

// EvalE evaluates an expression for given parsed structure in 'e'
func EvalE(aex ast.Expr, vars map[string]float64, funcs map[string]interface{}) (res float64, err error) {

	// Note: Based on Rémy Oudompheng' code: http://play.golang.org/p/HA8-8lxdtV
	// recursively evaluate tree
	var errors string
	var eval_tree func(ast.Expr) float64
	eval_tree = func(e ast.Expr) float64 {
		switch o := e.(type) {

		// identifier == variable
		case *ast.Ident:
			v, ok := vars[o.Name]
			if !ok {
				v, ok = myvars[o.Name]
				if !ok {
					errors += Sf(_evalexpr_err1, o.Name)
					return 0
				}
			}
			return v

		// literal == float64
		case *ast.BasicLit:
			v, parseerror := strconv.ParseFloat(o.Value, 64)
			if parseerror != nil {
				errors += Sf(_evalexpr_err2, o.Value, parseerror.Error())
				return 0
			}
			return v

		// unary expression
		case *ast.UnaryExpr:
			x := eval_tree(o.X)
			switch o.Op {
			case token.SUB:
				return -1.0 * x
			}
			errors += Sf(_evalexpr_err3, o.Op)
			return 0

		// binary expression
		case *ast.BinaryExpr:
			x, y := eval_tree(o.X), eval_tree(o.Y)
			switch o.Op {
			case token.ADD:
				return x + y
			case token.SUB:
				return x - y
			case token.MUL:
				return x * y
			case token.QUO:
				return x / y
			}
			errors += Sf(_evalexpr_err4, o.Op)
			return 0

		// function
		case *ast.CallExpr:
			fcn := o.Fun.(*ast.Ident)
			f, ok := funcs[fcn.Name]
			if !ok {
				f, ok = myfuncs[fcn.Name]
				if !ok {
					errors += Sf(_evalexpr_err5, fcn.Name)
					return 0
				}
			}
			switch len(o.Args) {
			case 1:
				F := f.(func(float64) float64)
				return F(eval_tree(o.Args[0]))
			case 2:
				F := f.(func(a, b float64) float64)
				return F(eval_tree(o.Args[0]), eval_tree(o.Args[1]))
			}
			errors += Sf(_evalexpr_err6, fcn.Name)
			return 0

		// parentheses
		case *ast.ParenExpr:
			return eval_tree(o.X)
		}

		// error
		errors += Sf(_evalexpr_err7, e)
		return 0
	}

	// do evaluate
	res = eval_tree(aex)

	// check errors
	if errors != "" {
		err = Err(errors)
	}
	return
}

// Eval evaluates an expression with variables in vars and functions in funcs
func Eval(expr string, vars map[string]float64, funcs map[string]interface{}) float64 {
	aex, err := ParseE(expr)
	if err != nil {
		Panic(_evalexpr_err8, err)
	}
	res, err := EvalE(aex, vars, funcs)
	if err != nil {
		Panic(_evalexpr_err9, err)
	}
	return res
}

// error messages
var (
	_evalexpr_err1 = "evalexpr.go: EvalE: variable named '%s' is not available in 'vars' map"
	_evalexpr_err2 = "evalexpr.go: EvalE: cannot convert '%s' to float64 number\n%v"
	_evalexpr_err3 = "evalexpr.go: EvalE: token '%v' cannot be handled in unary expression"
	_evalexpr_err4 = "evalexpr.go: EvalE: token '%v' cannot be handled in binary expression"
	_evalexpr_err5 = "evalexpr.go: EvalE: function named '%s' is not available in 'funcs' map"
	_evalexpr_err6 = "evalexpr.go: EvalE: cannot parse function due to incorrect number of arguments: '%#v'"
	_evalexpr_err7 = "evalexpr.go: EvalE: type '%#v' cannot be parsed"
	_evalexpr_err8 = "evalexpr.go: Eval: cannot parse expression:\n%v"
	_evalexpr_err9 = "evalexpr.go: Eval: cannot evaluate expression:\n%v"
)
