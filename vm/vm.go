package vm

import (
	"fmt"
	"monkey/code"
	"monkey/compiler"
	"monkey/object"
)

const StackSize = 4096

type VM struct {
	constants    []object.Object
	instructions code.Instructions

	stack []object.Object
	sp    int
}

func New(bytecode *compiler.Bytecode) *VM {
	return &VM{
		instructions: bytecode.Instructions,
		constants:    bytecode.Constants,
		stack:        make([]object.Object, StackSize),
		sp:           0,
	}
}

func (vm *VM) StackTop() object.Object {
	if vm.sp == 0 {
		return nil
	}
	return vm.stack[vm.sp-1]
}

func (vm *VM) Run() error {
	for ip := 0; ip < len(vm.instructions); ip++ {
		op := code.Opcode(vm.instructions[ip])

		switch op {
		case code.OpConstant:
			constIdx := code.ReadUint16(vm.instructions[ip+1:])
			ip += 2

			if err := vm.push(vm.constants[constIdx]); err != nil {
				return err
			}

		case code.OpAdd:
			r := vm.pop()
			l := vm.pop()

			lValue := l.(*object.Integer).Value
			rValue := r.(*object.Integer).Value

			res := lValue + rValue
			vm.push(&object.Integer{Value: res})
		}
	}

	return nil
}

func (vm *VM) push(o object.Object) error {
	if vm.sp >= StackSize {
		return fmt.Errorf("stack overflow")
	}

	vm.stack[vm.sp] = o
	vm.sp++

	return nil
}

func (vm *VM) pop() object.Object {
	o := vm.stack[vm.sp-1]
	vm.sp--

	return o
}
