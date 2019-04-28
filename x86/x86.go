package x86

import (
	"fmt"
	"tgoc/ast"
	"tgoc/utils"
)

// Identifier name: offset from bsp
var offsets map[string]int

// The number of stored identifier to stack
var varCount int

// The number of total identifier
var varNum int

// To assign a unique number to a label.
var labelCount int

//
var returnFlg bool

func initi(n int) {
	offsets = map[string]int{}
	varCount = 1
	varNum = n
	labelCount = 0
}

func genExpr(expr ast.Node) {

	switch expr := expr.(type) {
	case *ast.IntLit:
		fmt.Printf("	push %d\n", expr.Val)
	case *ast.Boolean:
		if expr.Val {
			fmt.Printf("	push 1\n")
		} else {
			fmt.Printf("	push 0\n")
		}
	case *ast.LogicalExpr:
		genExpr(expr.Lhs)
		genExpr(expr.Rhs)

		fmt.Printf("	pop rdi\n")
		fmt.Printf("	pop rax\n")

		switch expr.Op {
		case "==":
			fmt.Printf("	cmp rax, rdi\n")
			fmt.Printf("	sete al\n")
			fmt.Printf("	movzx rax, al\n")
		case "!=":
			fmt.Printf("	cmp rax, rdi\n")
			fmt.Printf("	sete al\n")
			fmt.Printf("	movzx rax, al\n")
			// 0000 => 0001, 0001 => 0000
			fmt.Printf("	xor rax, 1\n")
		case "<":
			fmt.Printf("	cmp rax, rdi\n")
			fmt.Printf("	setl al\n")
			fmt.Printf("	movzx rax, al\n")
		case "<=":
			fmt.Printf("	cmp rax, rdi\n")
			fmt.Printf("	setle al\n")
			fmt.Printf("	movzx rax, al\n")
		case ">":
			fmt.Printf("	cmp rax, rdi\n")
			fmt.Printf("	setg al\n")
			fmt.Printf("	movzx rax, al\n")
		case ">=":
			fmt.Printf("	cmp rax, rdi\n")
			fmt.Printf("	setge al\n")
			fmt.Printf("	movzx rax, al\n")
		case "&&":
			fmt.Printf("	and rax, rdi\n")
		case "||":
			fmt.Printf("	or rax, rdi\n")
		}
		fmt.Printf("	push rax\n")

	case *ast.BinaryExpr:
		genExpr(expr.Lhs)
		genExpr(expr.Rhs)

		fmt.Printf("	pop rdi\n")
		fmt.Printf("	pop rax\n")

		switch expr.Op {
		case "+":
			fmt.Printf("	add rax, rdi\n")
		case "-":
			fmt.Printf("	sub rax, rdi\n")
		case "*":
			fmt.Printf("	mul rdi\n")
		case "/":
			fmt.Printf("    xor rdx, rdx\n")
			fmt.Printf("    div rdi\n")
		case "%":
			fmt.Printf("    xor rdx, rdx\n")
			fmt.Printf("    div rdi\n")
			fmt.Printf("	mov rax, rdx\n")
		case "<<":
			// To change the cl value, changed the rcx value.
			// cl is lower 8 bit register of rcx register.
			fmt.Printf("	mov rcx, rdi\n")
			fmt.Printf("	shl rax, cl\n")
		case ">>":
			fmt.Printf("	mov rcx, rdi\n")
			fmt.Printf("	sar rax, cl\n")
		case "&":
			fmt.Printf("	and rax, rdi\n")
		case "|":
			fmt.Printf("	or rax, rdi\n")
		case "^":
			fmt.Printf("	xor rax, rdi\n")
		case "&^":
			fmt.Printf("	xor rdi, rax\n")
			fmt.Printf("	and rax, rdi\n")
		}
		fmt.Printf("	push rax\n")

	case *ast.UnaryExpr:
		genExpr(expr.Expr)
		fmt.Printf("    pop rax\n")

		switch expr.Op {
		case "-":
			fmt.Printf("	neg rax\n")
		case "!":
			fmt.Printf("	xor rax, 1\n")
		}
		fmt.Printf("	push rax \n")

	case *ast.Ident:
		os, ok := offsets[expr.Name]
		utils.Assert(ok, "undefined identifier")
		fmt.Printf("	mov rax, QWORD PTR [rbp - %d]\n", 8*os)
		fmt.Printf("	push rax\n")
	}
}

func genDecl(decl ast.Decl) {
	svd, _ := decl.(*ast.SVDecl)
	genExpr(svd.Val)
	fmt.Printf("	pop rax\n")
	fmt.Printf("	mov QWORD PTR [rbp - %d], rax\n", 8*varCount)
	offsets[svd.Name] = varCount
	varCount++
}

func genStmt(stmt ast.Stmt) {
	switch stmt := stmt.(type) {
	case *ast.ExprStmt:
		genExpr(stmt.Expr)
		fmt.Printf("	pop rax\n")
	case *ast.DeclStmt:
		genDecl(stmt.Decl)
	case *ast.AssignStmt:
		genExpr(stmt.Val)
		fmt.Printf("	pop rax\n")
		os, ok := offsets[stmt.Name]
		utils.Assert(ok, "undefined identifier")
		fmt.Printf("	mov QWORD PTR [rbp - %d], rax\n", 8*os)
	case *ast.ReturnStmt:
		genExpr(stmt.Expr)
		fmt.Printf("	pop rax\n")

		// printNumStdout1()

		// fmt.Printf("	mov rsp, rbp\n")
		// fmt.Printf("	pop rbp\n")
		// fmt.Printf("	ret\n")

		// printNumStdout2()
		return
	case *ast.IfStmt:
		genExpr(stmt.Cond)
		fmt.Printf("	pop rax\n")
		fmt.Printf("	cmp rax, 0\n")

		lAlt := makeLabel()

		fmt.Printf("	je .L%s\n", lAlt)
		genStmts(stmt.Cons)
		if stmt.Alt != nil {
			lEnd := makeLabel()
			fmt.Printf("	jmp .L%s\n", lEnd)
			fmt.Printf(".L%s:\n", lAlt)
			genStmts(stmt.Alt)
			fmt.Printf(".L%s:\n", lEnd)
		} else {
			fmt.Printf(".L%s:\n", lAlt)
		}
	case *ast.ForSingleStmt:
		loop := makeLabel()
		slipOut := makeLabel()
		fmt.Printf(".LOOP%s:\n", loop)
		genExpr(stmt.Cond)
		fmt.Printf("	pop rax\n")
		fmt.Printf("	cmp rax, 0\n")
		fmt.Printf("	je .L%s\n", slipOut)
		genStmts(stmt.Stmts)
		fmt.Printf("	jmp .LOOP%s\n", loop)
		fmt.Printf(".L%s:\n", slipOut)
	case *ast.ForClauseStmt:
		loop := makeLabel()
		slipOut := makeLabel()
		genStmt(stmt.Init)
		fmt.Printf(".LOOP%s:\n", loop)
		genExpr(stmt.Cond)
		fmt.Printf("	pop rax\n")
		fmt.Printf("	cmp rax, 0\n")
		fmt.Printf("	je .L%s\n", slipOut)
		genStmts(stmt.Stmts)
		genStmt(stmt.Post)
		fmt.Printf("	jmp .LOOP%s\n", loop)
		fmt.Printf(".L%s:\n", slipOut)
	}
}

func genStmts(stmts []ast.Stmt) {
	for _, stmt := range stmts {
		rs, ok := stmt.(*ast.ReturnStmt)
		if !ok {
			genStmt(stmt)
		} else {
			genStmt(rs)
			fmt.Printf("	jmp _end\n")
			return
		}
	}
}

func gen(stmts []ast.Stmt) {

	if varNum > 0 {
		// なぜ一つの変数につき、rspを16下げる？（8ではなく）
		//fmt.Printf("	sub rsp, %d\n", varNum*16)
		fmt.Printf("	sub rsp, %d\n", varNum*8)
	}
	genStmts(stmts)
}

func Gen(stmts []ast.Stmt, varNum int) {
	initi(varNum)

	//fmt.Printf(".section	__TEXT,__text,regular,pure_instructions\n")

	fmt.Printf("	.intel_syntax noprefix\n")
	fmt.Printf("	.globl _main\n")

	fmt.Printf("_main:\n")
	fmt.Printf("	push rbp\n")
	fmt.Printf("	mov rbp, rsp\n")

	gen(stmts)

	//rintNumStdout1()

	fmt.Printf("_end:\n")
	fmt.Printf("	mov rsp, rbp\n")
	fmt.Printf("	pop rbp\n")
	fmt.Printf("	ret\n")

	//printNumStdout2()
}

func makeLabel() string {
	l := fmt.Sprintf("%04d", labelCount)
	labelCount++
	return l
}

func printNumStdout1() {
	fmt.Printf("	lea	rdi, [rip + L_.str]\n")
	fmt.Printf("	mov	rsi, rax\n")
	//fmt.Printf("	movabs rsi, rax\n")
	fmt.Printf("	mov	al, 0\n")
	fmt.Printf("	call	_printf\n")

	// fmt.Printf("	xor	ecx, ecx\n")
	// fmt.Printf("	mov	dword ptr [rbp - 4], eax\n")
	// fmt.Printf("	mov	eax, ecx\n")
}

func printNumStdout2() {
	//fmt.Printf(".section	__TEXT,__cstring,cstring_literals\n")
	fmt.Printf("L_.str:\n")
	fmt.Printf("	.asciz	\"%%ld\\n\"\n")
}
