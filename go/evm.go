// Package evm is an **incomplete** implementation of the Ethereum Virtual
// Machine for the "EVM From Scratch" course:
// https://github.com/w1nt3r-eth/evm-from-scratch
//
// To work on EVM From Scratch In Go:
//
// - Install Golang: https://golang.org/doc/install
// - Go to the `go` directory: `cd go`
// - Edit `evm.go` (this file!), see TODO below
// - Run `go test ./...` to run the tests
// package evm

// import (
// 	// "fmt"
// 	"math/big"
// )

// // Run runs the EVM code and returns the stack and a success indicator.
// func Evm(code []byte) ([]*big.Int, bool) {
// 	var stack []*big.Int
// 	pc := 0

// 	for pc < len(code) {
// 		op := code[pc]
// 		pc++

// 		if op >= 0x60 && op <= 0x7f { // PUSH1...PUSH32 opcodes
// 			numBytes := int(op) - 0x5f // Calculate the number of bytes to read
// 			if pc+numBytes <= len(code) {
// 				value := big.NewInt(0)
// 				for i := 0; i < numBytes; i++ {
// 					value.Lsh(value, 8)                                   // Shift existing value by 8 bits
// 					value.Or(value, big.NewInt(int64(code[pc+i]))) // Bitwise OR with the next byte
// 				}
// 				stack = append([]*big.Int{value}, stack...)
// 				pc += numBytes // Increment program counter by the number of bytes read
// 			} else {
// 				// Invalid PUSH instruction, code ends unexpectedly
// 				return nil, false
// 			}
// 		} else {
// 			switch op {
// 			case 0x00: // STOP opcode
// 				return stack, true
// 			case 0x5f: // PUSH0 opcode
// 				value := big.NewInt(0)
// 				stack = append(stack, value)
// 				pc++
// 			case 0x60: // PUSH2 opcode
//                 value := big.NewInt(int64(code[pc])<<8 | int64(code[pc+1]))
//                 stack = append(stack, value)
//                 pc += 2
// 			case 0x50:
// 				return stack[1:], true
// 			case 0x01:

// 				a:=stack[0]
// 				b:=stack[1]
// 				stack = stack[2:]
// 				c := big.NewInt(0)
// 				c = c.Add(a,b)
// 				stack = append(stack, c)

// 			default:
// 				// Unsupported opcode for now
// 				return nil, false
// 			}
// 		}

// 		// TODO: Implement the EVM here!

// 			 // delete this; it's only here to make the compiler think you're already using `op`
// 	}

//		return stack, true
//	}
package evm

import (
	"fmt"
	"math/big"

)
const UINT256MAX = 0xFFFFFFFFFFFFFFFF


// handleSTOP handles the STOP opcode.
func handleSTOP(stack []*big.Int) ([]*big.Int, bool) {
    return stack, true
}

// handlePUSH0 handles the PUSH0 opcode.
func handlePUSH0(stack []*big.Int, pc *int) []*big.Int {
    value := big.NewInt(0)
    stack = append(stack, value)
    *pc++
    return stack
}

// handlePUSHN handles the PUSHN (where N ranges from 1 to 32) opcodes.
func handlePUSHN(stack []*big.Int, code []byte, pc *int, numBytes int) []*big.Int {
    value := big.NewInt(0)
    for i := 0; i < numBytes; i++ {
        value.Lsh(value, 8)                                   // Shift existing value by 8 bits
        value.Or(value, big.NewInt(int64(code[*pc+i]))) // Bitwise OR with the next byte
    }
    stack = append([]*big.Int{value}, stack...)
    *pc += numBytes // Increment program counter by the number of bytes read
    return stack
}

// handleADD handles the ADD opcode.
// handleADD handles the ADD opcode.
func handleADD(stack []*big.Int) ([]*big.Int, bool) {
    a := stack[0]
    b := stack[1]
    stack = stack[2:]

    // Define UINT256MAX as a big integer
    UINT256MAX := new(big.Int).SetUint64(0xFFFFFFFFFFFFFFFF)

    // Perform addition
    c := new(big.Int).Add(a, b)

    // Perform modular arithmetic to handle overflow
    c.Mod(c, new(big.Int).Add(UINT256MAX, big.NewInt(1)))

    // Append the result back to the stack
    stack = append(stack, c)

    return stack, true
}


func pop(stack []*big.Int) ([]*big.Int, bool) {
	return stack[1:], true
}
func handleMultiplication(stack []*big.Int) ([]*big.Int, bool){
	// UINT256MAX := new(big.Int).SetUint64(0xFFFFFFFFFFFFFFFF)
	// UINT256MAX := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

	c:=big.NewInt(0)
	c=c.Mul(stack[0],stack[1])
	stack =stack[2:]
	c.Mod(c,new(big.Int).Exp(big.NewInt(2),big.NewInt(256),nil))


	stack = append(stack,c)
	return stack , true
}
func handleSubtraction(stack []*big.Int) ([]*big.Int, bool){
	UINT256MAX := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))
	c:=new(big.Int)
	c=c.Sub(stack[0],stack[1])
	stack = stack[2:]
	if c.Sign()<0{
		c.Add(c,new(big.Int).Add(UINT256MAX,big.NewInt(1)))
	}
	stack=append(stack, c)
	return stack, true
}
func handleDivision(stack []*big.Int) ([]*big.Int,error){
	
	b:=stack[1]
	d:=b.Uint64()

	c:=new(big.Int)
	c=c.Div(stack[0],stack[1])
	stack=stack[2:]
	if d==0{
		return nil,fmt.Errorf("Divided by ZERO")
	}
	stack =append(stack, c)
	return stack,nil
}
func handleMod(stack []*big.Int)([]*big.Int, error){
	
	c:=new(big.Int)
	c=c.Mod(stack[0],stack[1])
	
	if stack[1].Cmp(big.NewInt(0))==0{
		stack = stack[2:]
		stack = append(stack,big.NewInt(0))
		return stack, nil

	}
	stack = stack[2:]
	stack = append(stack, c)
	return stack, nil
}
func handleAddMod(stack []*big.Int)([]*big.Int, bool){
	c:=new(big.Int)
	d:=new(big.Int)
	d=d.Add(stack[0],stack[1])
	c=c.Mod(d,stack[2])
	stack=stack[3:]
	stack = append(stack, c)
	return stack, true

}
func handleMulMod(stack []*big.Int)([]*big.Int, bool){
	c:=new(big.Int)
	c=c.Mul(stack[0],stack[1])
	c=c.Mod(c,stack[2])
	stack=stack[3:]
	stack=append(stack, c)
	return stack,true
}
func exp(stack []*big.Int)([]*big.Int, bool){
	c:=new(big.Int)
	c=c.Exp(stack[0],stack[1],nil)
	stack =stack[2:]
	stack = append(stack, c)
	return stack ,true
}
// Evm runs the EVM code and returns the stack and a success indicator.
func Evm(code []byte) ([]*big.Int, bool) {
    var stack []*big.Int
    pc := 0

    for pc < len(code) {
        op := code[pc]
        pc++

        if op >= 0x60 && op <= 0x7f { // PUSH1...PUSH32 opcodes
            numBytes := int(op) - 0x5f // Calculate the number of bytes to read
            if pc+numBytes <= len(code) {
                stack = handlePUSHN(stack, code, &pc, numBytes)
            } else {
                // Invalid PUSH instruction, code ends unexpectedly
                return nil, false
            }
        } else {
            switch op {
            case 0x00: // STOP opcode
                return handleSTOP(stack)
            case 0x5f: // PUSH0 opcode
                stack = handlePUSH0(stack, &pc)
            case 0x01: // ADD opcode
				
                return handleADD(stack)
			case 0x50:
				return pop(stack)
			case 0x02:
				return handleMultiplication(stack)
			case 0x03:
				return handleSubtraction(stack)

			case 0x04:
				if len(stack) < 2 {
					return nil, false
				}
				if stack[1].Cmp(big.NewInt(0)) == 0 {
					stack = stack[2:]
					stack = append([]*big.Int{big.NewInt(0)}, stack...)
				} else {
					value := new(big.Int).Div(stack[0], stack[1])
					value.Mod(value, new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)) // Modulo operation to handle overflow
					stack = stack[2:]
					stack = append([]*big.Int{value}, stack...)
				}
			case 0x06:
				if len(stack) < 2 {
					return nil, false
				}
				if stack[1].Cmp(big.NewInt(0)) == 0 {
					stack = stack[2:]
					stack = append([]*big.Int{big.NewInt(0)}, stack...)
				} else {
					value := new(big.Int).Mod(stack[0], stack[1])
					// value.Mod(value, new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)) // Modulo operation to handle overflow
					stack = stack[2:]
					stack = append([]*big.Int{value}, stack...)
				}
			case 0x08:
				return handleAddMod(stack)
				
			case 0x09:
				return handleMulMod(stack)	
			case 0x0A:
				return exp(stack)
            default:
                // Unsupported opcode for now
                return nil, false
            }
        }
    }

    return stack, true
}