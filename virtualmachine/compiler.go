package main

import (
	"strconv"
)

type Compiler struct {
	code       string
	codeLength int
	position   int

	instructions []*Instruction
}

func NewCompiler(code string) *Compiler {
	return &Compiler{
		code:         code,
		codeLength:   len(code),
		instructions: []*Instruction{},
	}
}

func (c *Compiler) Compile() []*Instruction {
	loopStack := []int{}

	for c.position < c.codeLength {
		current := c.code[c.position]

		switch current {
		case '[':
			insPos := c.EmitWithArg(JumpIfZero, 0)
			loopStack = append(loopStack, insPos)
		case ']':
			// Pop position of last JumpIfZero ("[") instruction off stack
			openInstruction := loopStack[len(loopStack)-1]
			loopStack = loopStack[:len(loopStack)-1]
			// Emit the new JumpIfNotZero ("]") instruction, with correct position as argument
			closeInstructionPos := c.EmitWithArg(JumpIfNotZero, openInstruction)
			// Patch the old JumpIfZero ("[") instruction with new position
			c.instructions[openInstruction].Argument = closeInstructionPos

		case '+':
			c.CompileFoldableInstruction('+', Plus)
		case '-':
			c.CompileFoldableInstruction('-', Minus)
		case '<':
			c.CompileFoldableInstruction('<', Left)
		case '>':
			c.CompileFoldableInstruction('>', Right)
		case '.':
			c.CompileFoldableInstruction('.', PutChar)
		case ',':
			c.CompileFoldableInstruction(',', ReadChar)
		default:
			if isDigit(current) {
				c.CompileNumberedInstruction()
			}
		}

		c.position++
	}

	return c.instructions
}

func (c *Compiler) CompileFoldableInstruction(char byte, insType InsType) {
	count := 1

	for c.position < c.codeLength-1 && c.code[c.position+1] == char {
		count++
		c.position++
	}

	c.EmitWithArg(insType, count)
}

func (c *Compiler) EmitWithArg(insType InsType, arg int) int {
	ins := &Instruction{Type: insType, Argument: arg}
	c.instructions = append(c.instructions, ins)
	return len(c.instructions) - 1
}

func (c *Compiler) CompileNumberedInstruction() {
	startPos := c.position
	for c.position < c.codeLength-1 && isDigit(c.code[c.position+1]) {
		c.position++
	}
	n, err := strconv.Atoi(c.code[startPos : c.position+1])
	if err != nil {
		panic("failed to read digits")
	}
	var op byte
	if c.position < c.codeLength-1 {
		op = c.code[c.position+1]
	} else {
		panic("digits where not followed by operator")
	}

	switch op {
	case '+':
		c.EmitWithArg(Plus, n)
	case '-':
		c.EmitWithArg(Minus, n)
	case '<':
		c.EmitWithArg(Left, n)
	case '>':
		c.EmitWithArg(Right, n)
	case '.':
		c.EmitWithArg(PutChar, n)
	case ',':
		c.EmitWithArg(ReadChar, n)
	default:
		panic("booo")
	}
}

func isDigit(b byte) bool {
	return b >= '0' && b <= '9'
}
