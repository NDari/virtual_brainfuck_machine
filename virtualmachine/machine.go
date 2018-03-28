package main

import (
	"bufio"
)

type Machine struct {
	code []*Instruction
	ip   int

	memory [30000]int32
	dp     int

	input  *bufio.Reader
	output *bufio.Writer

	readBuf []rune
}

func NewMachine(instructions []*Instruction, in *bufio.Reader, out *bufio.Writer) *Machine {
	return &Machine{
		code:    instructions,
		input:   in,
		output:  out,
		readBuf: make([]rune, 1),
	}
}

func (m *Machine) Execute() {
	for m.ip < len(m.code) {
		ins := m.code[m.ip]

		switch ins.Type {
		case Plus:
			m.memory[m.dp] += int32(ins.Argument)
		case Minus:
			m.memory[m.dp] -= int32(ins.Argument)
		case Right:
			m.dp += ins.Argument
		case Left:
			m.dp -= ins.Argument
		case PutChar:
			for i := 0; i < ins.Argument; i++ {
				m.putChar()
			}
		case ReadChar:
			for i := 0; i < ins.Argument; i++ {
				m.readChar()
			}
		case JumpIfZero:
			if m.memory[m.dp] == 0 {
				m.ip = ins.Argument
				continue
			}
		case JumpIfNotZero:
			if m.memory[m.dp] != 0 {
				m.ip = ins.Argument
				continue
			}
		}

		m.ip++
	}
}

func (m *Machine) readChar() {
	n, s, err := m.input.ReadRune()
	if s != 4 {
		panic("Rune was not 4 bytes")
	}
	if err != nil {
		panic(err)
	}
	m.memory[m.dp] = n
}

func (m *Machine) putChar() {
	_, err := m.output.WriteRune(rune(m.memory[m.dp]))
	if err != nil {
		panic(err)
	}
	err = m.output.Flush()
	if err != nil {
		panic(err)
	}
}
