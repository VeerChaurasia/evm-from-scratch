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
	// "fmt"
	// "encoding/hex"
	"math/big"
	// "golang.org/x/text/cases"
)

const UINT256MAX = 0xFFFFFFFFFFFFFFFF

// handleSTOP handles the STOP opcode.
func handleSTOP(stack []*big.Int) ([]*big.Int, bool) {
	return stack, true
}

type Storage struct {
	data      []byte
	offsetMax int
}

func NewStorage(size int) *Storage {
	return &Storage{
		data:      make([]byte, size),
		offsetMax: 0,
	}
}

func (m *Storage) Store(value []byte) {
	m.data = append(m.data, value...)
}

// Load loads 32 bytes from memory at the specified offset.
func (m *Storage) Load(offset int) []byte {

	return m.data[offset:]
}

type Memory struct {
	data      []byte
	offsetMax int
}

func NewMemory(size int) *Memory {
	return &Memory{
		data:      make([]byte, size),
		offsetMax: 0,
	}
}
func (m *Memory) Store(offset int, value []byte) {
	m.MSIZE(offset)
	copy(m.data[offset:], value)
}

// Load loads 32 bytes from memory at the specified offset.
func (m *Memory) Load(offset int) []byte {
	m.MSIZE(offset)
	return m.data[offset : offset+32]
}

func (m *Memory) LoadforSHA3(offset int, size int) []byte {
	m.MSIZE(offset)
	return m.data[offset : offset+size]
}

func (m *Memory) Store8(offset int, value byte) {
	m.MSIZE(offset - 32)
	m.data[offset] = value
}

func (m *Memory) MSIZE(offset int) int {
	if offset+32 > m.offsetMax {
		m.offsetMax = offset + 32
	}
	return m.offsetMax
}

func (m *Memory) GetOffsetMax() int {
	return m.offsetMax
}


// handlePUSH0 handles the PUSH0 opcode.
func handlePUSH0(stack []*big.Int, pc *int) []*big.Int {
	value := big.NewInt(0)
	stack = append(stack, value)
	*pc++
	return stack
}

// handlePUSHN handles the PUSHN (where N ranges from 1 to 32) opcodes.


// handleADD handles the ADD opcode.
// handleADD handles the ADD opcode.
func handleADD(stack []*big.Int) ([]*big.Int, bool) {
	if len(stack)<2{
		return nil,false
	}
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
	stack = append([]*big.Int{c}, stack...)

	return stack, true
}


func handleMultiplication(stack []*big.Int) ([]*big.Int, bool) {
	// UINT256MAX := new(big.Int).SetUint64(0xFFFFFFFFFFFFFFFF)
	// UINT256MAX := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))
	if len(stack)<2{
		return nil,false
	}
	c := big.NewInt(0)
	c = c.Mul(stack[0], stack[1])
	stack = stack[2:]
	c.Mod(c, new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil)) //handling oveflow

	stack = append([]*big.Int{c}, stack...)
	return stack, true
}
func handleSubtraction(stack []*big.Int) ([]*big.Int, bool) {
	UINT256MAX := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))
	c := new(big.Int)
	c = c.Sub(stack[0], stack[1])
	stack = stack[2:]
	if c.Sign() < 0 {
		c.Add(c, new(big.Int).Add(UINT256MAX, big.NewInt(1))) //Handling Overflow
	}
	stack = append([]*big.Int{c}, stack...)
	
	return stack, true
}
func handleDivision(stack []*big.Int) ([]*big.Int, error) {

	// b:=stack[1]
	// d:=b.Uint64()

	c := new(big.Int)
	c = c.Div(stack[0], stack[1])
	stack = stack[2:]
	stack = append(stack, c)
	return stack, nil
}
func handleMod(stack []*big.Int) ([]*big.Int, error) {

	c := new(big.Int)
	c = c.Mod(stack[0], stack[1])

	if stack[1].Cmp(big.NewInt(0)) == 0 {
		stack = stack[2:]
		stack = append(stack, big.NewInt(0))
		return stack, nil

	}
	stack = stack[2:]
	stack = append(stack, c)
	return stack, nil
}
func handleAddMod(stack []*big.Int) ([]*big.Int, bool) {
	c := new(big.Int)
	d := new(big.Int)
	d = d.Add(stack[0], stack[1])
	c = c.Mod(d, stack[2])
	stack = stack[3:]
	stack = append([]*big.Int{c}, stack...)
	
	return stack, true

}
func handleMulMod(stack []*big.Int) ([]*big.Int, bool) {
	c := new(big.Int)
	c = c.Mul(stack[0], stack[1])
	c = c.Mod(c, stack[2])
	stack = stack[3:]
	stack = append(stack, c)
	return stack, true
}
func exp(stack []*big.Int) ([]*big.Int, bool) {
	c := new(big.Int)
	c = c.Exp(stack[0], stack[1], nil)
	stack = stack[2:]
	stack = append(stack, c)
	return stack, true
}
func handleSIGNEXTEND(stack []*big.Int) ([]*big.Int, bool) {
	if len(stack) < 2 {
		return nil, false
	}

	// Pop b and x from the stack
	b := stack[0]
	x := stack[1]
	stack = stack[2:]

	// Convert `b` to an integer and ensure it's within the bounds of our bit width (0-31 for a 256-bit number)
	bInt := int(b.Int64())
	if bInt >= 32 {
		return stack, false
	}

	// Calculate the number of bits to consider
	bits := (bInt + 1) * 8
	signBit := new(big.Int).Lsh(big.NewInt(1), uint(bits-1))

	// Check if the sign bit is set
	if x.Cmp(signBit) >= 0 {
		// If the sign bit is set, extend with 1s
		extended := new(big.Int).Lsh(big.NewInt(1), uint(256-bits))
		extended.Sub(extended, big.NewInt(1))
		extended.Lsh(extended, uint(bits))
		x.Or(x, extended)
	} else {
		// Ensure higher bits are zero
		mask := new(big.Int).Lsh(big.NewInt(1), uint(bits))
		mask.Sub(mask, big.NewInt(1))
		x.And(x, mask)
	}

	// Push the result back onto the stack
	stack = append([]*big.Int{x}, stack...)
	return stack, true
}
func sMod(stack []*big.Int) ([]*big.Int, bool) {
	if len(stack) < 2 {
		return nil, false
	}
	if stack[1].Cmp(big.NewInt(0)) == 0 {
		stack = stack[2:]
		stack = append(stack, big.NewInt(0))

	} else {
		value1 := stack[0].Int64()
		value2 := stack[1].Int64()
		int8Value1 := int8(value1)
		int8Value2 := int8(value2)
		value := int8Value1 % int8Value2
		bits := 8
		if value < 0 {
			value8 := new(big.Int).Add(big.NewInt(int64(256)), big.NewInt(int64(value)))
			extended := new(big.Int).Lsh(big.NewInt(1), uint(256-bits))
			extended.Sub(extended, big.NewInt(1))
			extended.Lsh(extended, uint(bits))
			value8.Or(value8, extended)
			stack = stack[2:]
			stack = append(stack, value8)
		} else {
			stack = stack[2:]
			stack = append(stack, big.NewInt(int64(value)))

		}
	}
	return stack, true
}
func handleLT(stack []*big.Int) ([]*big.Int, bool) {
	if len(stack) < 2 {
		return nil, false

		// Pop a and b from the stack
	}
	a := stack[0]
	b := stack[1]
	if a.Cmp(b) == 0 || a.Cmp(b) == 1 {
		stack = stack[2:]
		stack = append(stack, big.NewInt(0))
	} else {
		stack = stack[2:]
		stack = append(stack, big.NewInt(1))
	}
	return stack, true
}
func handleGT(stack []*big.Int) ([]*big.Int, bool) {
	if len(stack) < 2 {
		return nil, false

	}
	a := stack[0]
	b := stack[1]
	if a.Cmp(b) == 0 || a.Cmp(b) == -1 {
		stack = stack[2:]
		stack = append(stack, big.NewInt(0))
	} else {
		stack = stack[2:]
		stack = append(stack, big.NewInt(1))

	}
	return stack, true
}
func handleSLT(stack []*big.Int) ([]*big.Int, bool) {
	if len(stack) < 2 {
		return nil, false
	}
	value1 := stack[0].Int64()
	int8Value1 := int8(value1)
	value2 := stack[1].Int64()
	int8Value2 := int8(value2)
	if int8Value1 < int8Value2 {
		stack = stack[2:]
	stack = append([]*big.Int{big.NewInt(1)}, stack...)
		

	} else {
		stack = stack[2:]
		stack = append([]*big.Int{big.NewInt(0)}, stack...)
		

	}
	return stack, true
}
func handleGLT(stack []*big.Int) ([]*big.Int, bool) {
	if len(stack) < 2 {
		return nil, false
	}
	value1 := stack[0].Int64()
	int8Value1 := int8(value1)
	value2 := stack[1].Int64()
	int8Value2 := int8(value2)
	if int8Value1 > int8Value2 {
		stack = stack[2:]
		// stack = append(stack, big.NewInt(1))
		stack = append([]*big.Int{big.NewInt(1)}, stack...)


	} else {
		stack = stack[2:]
		// stack = append(stack, big.NewInt(0))
		stack = append([]*big.Int{big.NewInt(0)}, stack...)


	}
	return stack, true
}
func handleEQ(stack []*big.Int) ([]*big.Int, bool) {
	if len(stack) < 2 {
		return nil, false
	}
	if stack[0].Cmp(stack[1]) == 0 {
		stack = stack[2:]
		// stack = append(stack, big.NewInt(1))
		stack = append([]*big.Int{big.NewInt(1)}, stack...)

	} else {
		stack = stack[2:]
		// stack = append(stack, big.NewInt(0))
		stack = append([]*big.Int{big.NewInt(0)}, stack...)

	}
	return stack, true
}
func handleISZERO(stack []*big.Int) ([]*big.Int, bool) {
	// if len(stack)<1{
	// 	return nil, false
	// }
	if stack[0].Cmp(big.NewInt(0)) == 0 {
		stack = stack[1:]
		// stack = append(stack, big.NewInt(1))
		stack = append([]*big.Int{big.NewInt(1)}, stack...)

	} else {
		stack = stack[1:]
		// stack = append(stack, big.NewInt(0))
		stack = append([]*big.Int{big.NewInt(0)}, stack...)


	}
	return stack, true
}
func handleNot(stack []*big.Int) ([]*big.Int, bool) {
	if len(stack) < 1 {
		return nil, false
	}

	value := stack[0]
	stack = stack[1:]

	// Create a bitmask with all bits set to 1 for 256 bits
	bitmask := new(big.Int).Sub(new(big.Int).Lsh(big.NewInt(1), 256), big.NewInt(1))

	value1 := new(big.Int).Xor(value, bitmask)

	stack = append(stack, value1)

	return stack, true
}
func handleAnd(stack []*big.Int) ([]*big.Int, bool) {
	if len(stack) < 2 {
		return nil, false
	}
	value := stack[0]
	value1 := stack[1]
	stack = stack[2:]
	c := new(big.Int).And(value, value1)
	stack = append(stack, c)
	return stack, true
}
func handleOr(stack []*big.Int) ([]*big.Int, bool) {
	if len(stack) < 2 {
		return nil, false
	}
	value := stack[0]
	value1 := stack[1]
	stack = stack[2:]
	c := new(big.Int).Or(value, value1)
	stack = append(stack, c)
	return stack, true
}
func handleXor(stack []*big.Int) ([]*big.Int, bool) {
	if len(stack) < 2 {
		return nil, false
	}
	value := stack[0]
	value1 := stack[1]
	stack = stack[2:]
	c := new(big.Int).Xor(value, value1)
	stack = append(stack, c)
	return stack, true
}
func handleSHL(stack []*big.Int) ([]*big.Int, bool) {
	if len(stack) < 2 {
		return nil, false
	}

	extended := stack[1]
	UINT256Max := new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil), big.NewInt(1))

	if stack[0].Cmp(big.NewInt(255)) > 1 {
		extended = big.NewInt(0)
	} else {
		extended = new(big.Int).Lsh(stack[1], uint(stack[0].Int64()))
		extended.And(extended, UINT256Max)
	}

	stack = stack[2:]
	stack = append(stack, extended)
	return stack, true
}
func handleSHR(stack []*big.Int) ([]*big.Int, bool) {
	if len(stack) < 2 {
		return nil, false
	}

	value := stack[1]
	UINT256MAX := new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil), big.NewInt(1))
	if stack[0].Cmp(UINT256MAX) > 1 {
		value = big.NewInt(0)
		stack = stack[2:]
		stack = append(stack, value)
	} else {
		value = new(big.Int).Rsh(stack[1], uint(stack[0].Int64()))
		value = new(big.Int).And(value, UINT256MAX)
		stack = stack[2:]
		stack = append(stack, value)
	}
	return stack, true
}
func jumpdest(pc int, code []byte, stack []*big.Int) (int, []*big.Int, bool) {
	stack = []*big.Int{}
loop:
	for i := pc; i < len(code); i++ {
		if code[i] == 0x5B {

			switch code[i-1] {
			case 0x60, 0x61, 0x62, 0x63, 0x64, 0x65, 0x66, 0x67, 0x68, 0x69, 0x6A, 0x6B, 0x6C, 0x6D, 0x6E, 0x6F, 0x70, 0x71, 0x72, 0x73, 0x74, 0x75, 0x76, 0x77, 0x78, 0x79, 0x7A, 0x7B, 0x7C, 0x7D, 0x7E, 0x7F:

			default:
				pc = i
				break loop
			}
		}
		if i == len(code)-1 && code[len(code)-1] != 0x5B {
			return pc, stack, false
		}

	}
	return pc, stack, true
}





func Evm(code []byte,) ([]*big.Int, bool) {
	var stack []*big.Int
	memory := NewMemory(1024)
	pc := 0
	
	// ret := ""

	for pc < len(code) {
		op := code[pc]
		pc++
		switch op {
		default:
			if 0x60 <= op && op <= 0x7f {
				increment := int(op-0x60) + 1

				if pc+increment > len(code) {
					return nil, false

				}
				value := new(big.Int).SetBytes(code[pc : pc+increment])
				stack = append([]*big.Int{value}, stack...)
				pc += increment
			}

			
		case 0x00: // STOP opcode
			return handleSTOP(stack)
		case 0x5F:
			value := big.NewInt(0)
			stack = append([]*big.Int{value}, stack...)
		case 0x01: // ADD opcode

			return handleADD(stack)
		case 0x50:
			// return pop(stack)
			if len(stack) < 1 {
				return nil, false
			}

			stack = stack[1:]
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
		case 0x0B:

			if len(stack) < 2 {
				return nil, false
			}

			b := stack[0]
			x := stack[1]
			stack = stack[2:]

			// Calculate the sign extension mask
			bInt := int(b.Int64())
			if bInt >= 32 {
				return stack, false
			}
			bits := (bInt + 1) * 8
			signBit := new(big.Int).Lsh(big.NewInt(1), uint(bits-1))

			// Check if the sign bit is set
			if x.Cmp(signBit) >= 0 {
				// If the sign bit is set, extend with 1s
				extended := new(big.Int).Lsh(big.NewInt(1), uint(256-bits))
				extended.Sub(extended, big.NewInt(1))
				extended.Lsh(extended, uint(bits))
				x.Or(x, extended)
			} else {
				// Ensure higher bits are zero
				mask := new(big.Int).Lsh(big.NewInt(1), uint(bits))
				mask.Sub(mask, big.NewInt(1))
				x.And(x, mask)
			}
			stack = append([]*big.Int{x}, stack...)

		case 0x05:

			if len(stack) < 2 {
				return nil, false
			}

			if stack[1].Cmp(big.NewInt(0)) == 0 {
				stack = stack[2:]
				stack = append([]*big.Int{big.NewInt(0)}, stack...)
			} else {
				value1 := stack[0].Int64()
				int8Value1 := int8(value1)
				value2 := stack[1].Int64()
				int8Value2 := int8(value2)

				value := int8Value1 / int8Value2

				bits := 8

				// Check if the sign bit is set
				if value < 0 {
					value8 := new(big.Int).Add(big.NewInt(int64(256)), big.NewInt(int64(value)))
					// If the sign bit s set, extend with 1s
					extended := new(big.Int).Lsh(big.NewInt(1), uint(256-bits))
					extended.Sub(extended, big.NewInt(1))
					extended.Lsh(extended, uint(bits))
					value8.Or(value8, extended)
					stack = stack[2:]
					stack = append(stack, value8)
				} else {
					stack = stack[2:]
					stack = append([]*big.Int{big.NewInt(int64(value))}, stack...)
				}

			}
		case 0x07:
			return sMod(stack)
		case 0x10:
			return handleLT(stack)
		case 0x11:
			return handleGT(stack)
		case 0x12:
			return handleSLT(stack)
		case 0x13:
			return handleGLT(stack)
		case 0x14:
			return handleEQ(stack)

		case 0x15:
			return handleISZERO(stack)
		case 0x16:
			return handleAnd(stack)
		case 0x17:
			return handleOr(stack)
		case 0x18:
			return handleXor(stack)
		case 0x19:
			return handleNot(stack)

		case 0x1B:
			return handleSHL(stack)
		case 0x1C:
			return handleSHR(stack)
		case 0x1D:
			if len(stack) < 2 {
				return nil, false
			}
			extended := stack[1]
			UINT256Max := new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil), big.NewInt(1))
			INT256MAX := new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(255), nil), big.NewInt(1))
			Check_val := new(big.Int).Sub(INT256MAX, extended)
			if stack[0].Cmp(big.NewInt(256)) != -1 {
				mask := big.NewInt(1)
				mask = mask.Lsh(mask, 255)
				// Create a mask for the first bit
				firstBitMask := stack[1]
				// Assuming mask is 256 bits

				// Extract the first bit of mask
				firstBit := new(big.Int).And(mask, firstBitMask)

				if firstBit.Cmp(big.NewInt(0)) == 0 {
					// If the first bit is 0
					extended = new(big.Int).Lsh(mask, 1)
				} else {
					// If the first bit is 1
					mask = new(big.Int).Lsh(mask, 1)
					mask = new(big.Int).Sub(mask, big.NewInt(1))
					extended = mask
				}

			} else {
				if Check_val.Cmp(big.NewInt(0)) == -1 {
					extended := stack[1]
					shift := uint(stack[0].Uint64())

					// Create a mask that has ones in the positions that should be filled with ones after the shift
					mask := new(big.Int).Lsh(big.NewInt(1), shift)
					mask.Sub(mask, big.NewInt(1))
					mask.Lsh(mask, 256-shift)

					// Perform the right shift and apply the mask
					extended.Rsh(extended, shift)
					extended.Or(extended, mask)
				} else {
					extended.Rsh(extended, uint(stack[0].Int64()))
				}
			}

			extended.And(extended, UINT256Max)
			stack = stack[2:]

			stack = append([]*big.Int{extended}, stack...)
		case 0x1A:
			extended := stack[1]
			UINT256Max := new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil), big.NewInt(1))
			extended.And(extended, UINT256Max)
			mask := big.NewInt(255)

			shift := (31 - stack[0].Int64()) * 8
			if shift >= 0 && shift <= 256 {
				mask = mask.Lsh(mask, uint(shift))
			}

			extended = extended.And(extended, mask)
			extended = extended.Rsh(extended, uint(shift))
			extended = extended.And(extended, big.NewInt(255))
			stack = stack[2:]
			stack = append([]*big.Int{extended}, stack...)
		case 0x80, 0x81, 0x82, 0x83, 0x84, 0x85, 0x86, 0x87, 0x88, 0x89:
			op2 := op - 0x80
			dup := stack[op2]
			stack = append([]*big.Int{dup}, stack...)
		case 0x90:
			a := stack[0]
			b := stack[1]
			stack = stack[2:]
			stack = append(stack, b)
			stack = append(stack, a)

		case 0x91, 0x92, 0x93, 0x94, 0x95, 0x96, 0x97, 0x98, 0x99, 0x9a, 0x9b, 0x9c, 0x9d, 0x9e, 0x9f:
			// SWAP1 through SWAP16
			depth := int(op-0x90) + 1
			if len(stack) < depth+1 {
				return nil, false // insufficient stack depth
			}

			stack[0], stack[depth] = stack[depth], stack[0]

		case 0xFE:
			return nil, false

		case 0x58:

			counter := 0
			for i := pc - 2; i >= 0; i-- {
				if i < len(code) {
					if code[i] == 60 {
						counter = counter + 2
					} else {
						counter++
					}
				}

			}
			stack = append([]*big.Int{big.NewInt(int64(counter))}, stack...)
		case 0x5A:
			UINT256Max := new(big.Int).Sub(new(big.Int).Exp(big.NewInt(2), big.NewInt(256), nil), big.NewInt(1))
			stack = append([]*big.Int{UINT256Max}, stack...)
		
		case 0x56:
			hello := true
			pc, stack, hello = jumpdest(pc, code, stack)
			if !hello {
				return stack, false
			}
		case 0x57:
			value := stack[1]
			if value.Cmp(big.NewInt(0)) != 0 {
				hello := true
				pc, stack, hello = jumpdest(pc, code, stack)
				if !hello {
					return stack, false
				}
			} else {
				stack = []*big.Int{}

			}
		case 0x52: // MSTORE
			offset := stack[0]
			stack = stack[1:]
			value := stack[0]
			stack = stack[1:]
			offsetInt := int(offset.Int64())
			valueBytes := value.Bytes()

			if len(valueBytes) < 32 {
				padding := make([]byte, 32-len(valueBytes))
				valueBytes = append(padding, valueBytes...)
			}
			memory.Store(offsetInt, valueBytes)
		case 0x51: // MLOAD
			offset := stack[0]
			stack = stack[1:]
			offsetInt := int(offset.Int64())
			value := new(big.Int).SetBytes(memory.Load(offsetInt))
			stack = append([]*big.Int{value}, stack...)
		case 0x53: // MSTORE8
			offset := stack[0]

			stack = stack[1:]

			value := int8(stack[0].Int64())
			stack = stack[1:]
			offsetInt := int(offset.Uint64())
			memory.Store8(offsetInt, byte(value))

		case 0x59:
			value := memory.GetOffsetMax()
			m := 32
			value1 := ((value + m - 1) / m) * m
			stack = append([]*big.Int{big.NewInt(int64(value1))}, stack...)
		
		
	}	}

	return stack, true
}
