package fixture

import "C"

// IntEQ mimics CGo LLVM enum constants such as llvm.IntEQ = C.LLVMIntEQ.
type IntPredicate int

const IntEQ IntPredicate = C.LLVMIntEQ

func CreateICmp(pred IntPredicate) {}

func useInvalidCConst() {
	CreateICmp(IntEQ)
}
