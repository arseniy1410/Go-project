package main

import (
	"fmt"
	"strconv"
)

// Values

type Kind int

const (
	ValueInt  Kind = 0
	ValueBool Kind = 1
	Undefined Kind = 2
)

type Val struct {
	flag Kind
	valI int
	valB bool
}

func mkInt(x int) Val {
	return Val{flag: ValueInt, valI: x}
}
func mkBool(x bool) Val {
	return Val{flag: ValueBool, valB: x}
}
func mkUndefined() Val {
	return Val{flag: Undefined}
}

func showVal(v Val) string {
	var s string
	switch {
	case v.flag == ValueInt:
		s = Num(v.valI).pretty()
	case v.flag == ValueBool:
		s = Bool(v.valB).pretty()
	case v.flag == Undefined:
		s = "Undefined"
	}
	return s
}

// Types

type Type int

const (
	TyIllTyped Type = 0
	TyInt      Type = 1
	TyBool     Type = 2
)

func showType(t Type) string {
	var s string
	switch {
	case t == TyInt:
		s = "Int"
	case t == TyBool:
		s = "Bool"
	case t == TyIllTyped:
		s = "Illtyped"
	}
	return s
}

// Value State is a mapping from variable names to values
type ValState map[string]Val

// Value State is a mapping from variable names to types
type TyState map[string]Type

// Interface

type Exp interface {
	pretty() string
	eval(s ValState) Val
	infer(t TyState) Type
}

type Stmt interface {
	pretty() string
	eval(s ValState)
	check(t TyState) bool
}

// Statement cases (incomplete)

type Seq [2]Stmt

type Decl struct {
	lhs string
	rhs Exp
}
type IfThenElse struct {
	cond     Exp
	thenStmt Stmt
	elseStmt Stmt
}

type While struct {
	cond Exp
	stmt Stmt
}

type Assign struct {
	lhs string
	// lhs Exp
	rhs Exp
}

type Print [1]Exp

// Expression cases (incomplete)

type Bool bool
type Num int
type Negate [1]Exp
type Mult [2]Exp
type Plus [2]Exp
type And [2]Exp
type Or [2]Exp
type Equals [2]Exp
type Lesser [2]Exp
type Group [1]Exp
type Var string

//-----------------------------------
// Stmt instances

// pretty print

func (stmt Seq) pretty() string {
	return stmt[0].pretty() + "; " + stmt[1].pretty()
}

func (decl Decl) pretty() string {
	return decl.lhs + " := " + decl.rhs.pretty()
}

func (assgn Assign) pretty() string {
	return assgn.lhs + " = " + assgn.rhs.pretty()
}

func (ifStmnt IfThenElse) pretty() string {
	return "if" + ifStmnt.cond.pretty() + "{" + ifStmnt.thenStmt.pretty() + "} else {" + ifStmnt.elseStmt.pretty() + "}"
}

func (whl While) pretty() string {
	return "while" + whl.cond.pretty() + "{" + whl.stmt.pretty() + "}"
}

func (pt Print) pretty() string {
	return pt.pretty()
}

// eval

func (stmt Seq) eval(s ValState) {
	stmt[0].eval(s)
	stmt[1].eval(s)
}

func (ite IfThenElse) eval(s ValState) {
	v := ite.cond.eval(s)
	if v.flag == ValueBool {
		switch {
		case v.valB:
			ite.thenStmt.eval(s)
		case !v.valB:
			ite.elseStmt.eval(s)
		}

	} else {
		fmt.Printf("if-then-else eval fail")
	}

}

func (whl While) eval(s ValState) {
	v := whl.cond.eval(s)
	if v.flag == ValueBool {
		for {
			whl.stmt.eval(s)
			if !v.valB {
				break
			}
		}
	} else {
		fmt.Printf("while eval fail")
	}
}

func (prnt Print) eval(s ValState) {
	prnt[0].eval(s)
}

// Maps are represented via points.
// Hence, maps are passed by "reference" and the update is visible for the caller as well.
func (decl Decl) eval(s ValState) {
	v := decl.rhs.eval(s)
	x := (string)(decl.lhs)
	s[x] = v
}

func (asgn Assign) eval(s ValState) {
	v := asgn.rhs.eval(s)
	x := (string)(asgn.lhs)
	s[x] = v
}

// type check

func (stmt Seq) check(t TyState) bool {
	if !stmt[0].check(t) {
		return false
	}
	return stmt[1].check(t)
}

func (decl Decl) check(t TyState) bool {
	ty := decl.rhs.infer(t)
	if ty == TyIllTyped {
		return false
	}

	x := (string)(decl.lhs)
	t[x] = ty
	return true
}

func (a Assign) check(t TyState) bool {
	ty := a.rhs.infer(t)
	if ty == TyIllTyped {
		return false
	}

	x := (string)(a.lhs)
	t[x] = ty
	// if tx == TyIllTyped || ty != tx {
	// 	return false
	// }
	return true
}

func (ite IfThenElse) check(t TyState) bool {
	ty := ite.cond.infer(t)
	if ty == TyIllTyped {
		return false
	}

	th := ite.thenStmt.check(t)
	el := ite.elseStmt.check(t)

	return th && el
}

func (a While) check(t TyState) bool {
	ty := a.cond.infer(t)
	if ty == TyIllTyped {
		return false
	}

	return a.stmt.check(t)
}

//-----------------------------------
// Exp instances

// pretty print

func (x Var) pretty() string {
	return (string)(x)
}

func (x Bool) pretty() string {
	if x {
		return "true"
	} else {
		return "false"
	}
}

func (e Negate) pretty() string {
	var x string
	x = "not "
	x += e[0].pretty()

	return x
}

func (e Equals) pretty() string {
	var x string
	x = "("
	x += e[0].pretty()
	x += "=="
	x += e[1].pretty()
	x += ")"

	return x
}

func (e Lesser) pretty() string {
	var x string
	x = "("
	x += e[0].pretty()
	x += "<"
	x += e[1].pretty()
	x += ")"

	return x
}

func (x Num) pretty() string {
	return strconv.Itoa(int(x))
}

func (e Mult) pretty() string {

	var x string
	x = "("
	x += e[0].pretty()
	x += "*"
	x += e[1].pretty()
	x += ")"

	return x
}

func (e Plus) pretty() string {

	var x string
	x = "("
	x += e[0].pretty()
	x += "+"
	x += e[1].pretty()
	x += ")"

	return x
}

func (e And) pretty() string {

	var x string
	x = "("
	x += e[0].pretty()
	x += "&&"
	x += e[1].pretty()
	x += ")"

	return x
}

func (e Or) pretty() string {

	var x string
	x = "("
	x += e[0].pretty()
	x += "||"
	x += e[1].pretty()
	x += ")"

	return x
}

func (g Group) pretty() string {

	var x string
	x = "("
	x += g.pretty()
	x += ")"

	return x
}

// Evaluator

func (x Bool) eval(s ValState) Val {
	return mkBool((bool)(x))
}

func (e Negate) eval(s ValState) Val {
	b := e[0].eval(s)
	if b.flag == ValueBool {
		return mkBool((bool)(!b.valB))
	}
	return mkUndefined()
}

func (e Equals) eval(s ValState) Val {
	b1 := e[0].eval(s)
	b2 := e[1].eval(s)
	if b1.flag == ValueBool && b2.flag == ValueBool {
		return mkBool(b1.valB == b2.valB)
	} else if b1.flag == ValueInt && b2.flag == ValueInt {
		return mkBool(b1.valI < b2.valI)
	}
	return mkUndefined()
}

func (e Lesser) eval(s ValState) Val {
	n1 := e[0].eval(s)
	n2 := e[1].eval(s)
	if n1.flag == ValueInt && n2.flag == ValueInt {
		return mkBool(n1.valI < n2.valI)
	}
	return mkUndefined()
}

func (x Num) eval(s ValState) Val {
	return mkInt((int)(x))
}

func (e Mult) eval(s ValState) Val {
	n1 := e[0].eval(s)
	n2 := e[1].eval(s)
	if n1.flag == ValueInt && n2.flag == ValueInt {
		return mkInt(n1.valI * n2.valI)
	}
	return mkUndefined()
}

func (e Plus) eval(s ValState) Val {
	n1 := e[0].eval(s)
	n2 := e[1].eval(s)
	if n1.flag == ValueInt && n2.flag == ValueInt {
		return mkInt(n1.valI + n2.valI)
	}
	return mkUndefined()
}

func (e And) eval(s ValState) Val {
	b1 := e[0].eval(s)
	b2 := e[1].eval(s)
	switch {
	case b1.flag == ValueBool && b1.valB == false:
		return mkBool(false)
	case b1.flag == ValueBool && b2.flag == ValueBool:
		return mkBool(b1.valB && b2.valB)
	}
	return mkUndefined()
}

func (e Or) eval(s ValState) Val {
	b1 := e[0].eval(s)
	b2 := e[1].eval(s)
	switch {
	case b1.flag == ValueBool && b1.valB == true:
		return mkBool(true)
	case b1.flag == ValueBool && b2.flag == ValueBool:
		return mkBool(b1.valB || b2.valB)
	}
	return mkUndefined()
}

func (g Group) eval(s ValState) Val {
	e := g.eval(s)
	if e.flag == Undefined {
		return mkUndefined()
	}
	return e
}

// Type inferencer/checker

func (x Var) infer(t TyState) Type {
	y := (string)(x)
	ty, ok := t[y]
	if ok {
		return ty
	} else {
		return TyIllTyped // variable does not exist yields illtyped
	}

}

func (x Bool) infer(t TyState) Type {
	return TyBool
}

func (x Num) infer(t TyState) Type {
	return TyInt
}

func (e Negate) infer(t TyState) Type {
	t1 := e[0].infer(t)
	if t1 == TyBool {
		return TyBool
	}
	return TyIllTyped
}

func (e Mult) infer(t TyState) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == TyInt && t2 == TyInt {
		return TyInt
	}
	return TyIllTyped
}

func (e Plus) infer(t TyState) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == TyInt && t2 == TyInt {
		return TyInt
	}
	return TyIllTyped
}

func (e And) infer(t TyState) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 != TyBool {
		return TyIllTyped
	} else if t2 != TyBool {
		return TyIllTyped
	}
	return TyBool
}

func (e Or) infer(t TyState) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)

	if t1 != TyBool {
		return TyIllTyped
	} else if t2 != TyBool {
		return TyIllTyped
	}
	return TyBool
}

func (g Group) infer(t TyState) Type {
	ty := g.infer(t)

	if ty != TyIllTyped {
		return ty
	}
	return TyIllTyped
}

func (e Equals) infer(t TyState) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == TyInt && t2 == TyInt || t1 == TyBool && t2 == TyBool {
		return TyBool
	}
	return TyIllTyped
}

func (e Lesser) infer(t TyState) Type {
	t1 := e[0].infer(t)
	t2 := e[1].infer(t)
	if t1 == TyInt && t2 == TyInt {
		return TyInt
	}
	return TyIllTyped
}

// Helper functions to build ASTs by hand

func number(x int) Exp {
	return Num(x)
}

func boolean(x bool) Exp {
	return Bool(x)
}

func plus(x, y Exp) Exp {
	return (Plus)([2]Exp{x, y})

	// The type Plus is defined as the two element array consisting of Exp elements.
	// Plus and [2]Exp are isomorphic but different types.
	// We first build the AST value [2]Exp{x,y}.
	// Then cast this value (of type [2]Exp) into a value of type Plus.

}

func mult(x, y Exp) Exp {
	return (Mult)([2]Exp{x, y})
}

func and(x, y Exp) Exp {
	return (And)([2]Exp{x, y})
}

func or(x, y Exp) Exp {
	return (Or)([2]Exp{x, y})
}

func group(x Exp) Exp {
	return (Group)([1]Exp{x})
}

func negate(x Exp) Exp {
	// return Negate([1]Exp{x})
	neg := Negate([1]Exp{x})
	return neg
}

func equals(x, y Exp) Exp {
	return (Equals)([2]Exp{x, y})
}

func lesser(x, y Exp) Exp {
	return (Lesser)([2]Exp{x, y})
}

func decl(lhs string, rhs Exp) Stmt {
	return (Decl)(Decl{lhs, rhs})
}

func assig(lhs string, rhs Exp) Stmt {
	return (Assign)(Assign{lhs, rhs})
}

func whil(e Exp, s Stmt) Stmt {
	return (While)(While{e, s})
}

// func print(x Exp) Exp {
// 	return (Print)([1]Exp{x})
// }

func ite(cond Exp, stmt1, stmt2 Stmt) Stmt {
	return (IfThenElse)(IfThenElse{cond, stmt1, stmt2})
}

func seq(s1, s2 Stmt) Stmt {
	return (Seq)(Seq{s1, s2})
}

// Examples

func runExp(e Exp) {
	s := make(map[string]Val)
	t := make(map[string]Type)
	fmt.Printf("\n ******* ")
	fmt.Printf("\n %s", e.pretty())
	fmt.Printf("\n %s", showVal(e.eval(s)))
	fmt.Printf("\n %s", showType(e.infer(t)))
}

func runStmt(stmt Stmt) {
	s := make(map[string]Val)
	t := make(map[string]Type)
	fmt.Printf("\n ******* ")
	fmt.Printf("\n %s", stmt.pretty())
	// fmt.Printf("\n %s", showVal(e.eval(s)))
	stmt.eval(s)
	fmt.Printf("\n %t", stmt.check(t))
}

func ex1() {
	ast := plus(mult(number(1), number(2)), number(0))

	runExp(ast)
}

func ex2() {
	ast := and(boolean(false), number(0))
	runExp(ast)
}

func ex3() {
	ast := or(boolean(false), number(0))
	runExp(ast)
}

func ex4() {
	ast := negate(number(2))
	runExp(ast)
}

func ex5() {
	ast := lesser(number(2), number(4))
	runExp(ast)
}

func ex6() {
	// ast := decl("x", boolean(true))
	ast := assig("z", boolean(false))
	ast1 := assig("z", number(4))
	// ast := assig(decl("x", boolean(true)), number(4))
	runStmt(ast)
	runStmt(ast1)
}

func ex7() {
	ast := seq(decl("x", mult(number(2), number(2))), assig("x", boolean(true)))

	runStmt(ast)
}

func ex8() {
	ast := ite(lesser(number(5), number(4)), decl("x", boolean(true)), decl("y", number(3)))

	runStmt(ast)
}

// func ex9() {
// 	ast := whil(lesser(number(6), number(4)), seq(decl("x", mult(number(2), number(2))), assig("x", boolean(true))))
// 	seq(decl("x", number(2)), whil(lesser(x, number(4)), seq(decl("x", mult(number(2), number(2))), assig("x", boolean(true)))))

// 	runStmt(ast)
// }

func ex10() {
	ast := group(lesser(number(6), number(4)))
	runExp(ast)
}

func main() {

	fmt.Printf("\n")

	ex1()
	ex2()
	ex3()
	ex4()
	ex5()
	ex6()
	ex7()
	ex8()
	// ex9()
	// ex10() // ???????????????????????????????????
}
