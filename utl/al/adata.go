// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package al

// Adata defines an Abstract Data object by wrapping interface{}
type Adata interface{}

// strings ////////////////////////////////////////////////////////////////////////////////////////

// FromString converts string to Adata's interface{}
func FromString(val string) Adata {
	return Adata(val)
}

// ToString converts Adata's interface{} to string
func ToString(val interface{}) string {
	return val.(string)
}

// FromStrings converts []string to Adata's interface{}
func FromStrings(val []string) Adata {
	return Adata(val)
}

// ToStrings converts Adata's interface{} to []string
func ToStrings(val interface{}) []string {
	return val.([]string)
}

// ints ////////////////////////////////////////////////////////////////////////////////////////

// FromInt converts int to Adata's interface{}
func FromInt(val int) Adata {
	return Adata(val)
}

// ToInt converts Adata's interface{} to int
func ToInt(val interface{}) int {
	return val.(int)
}

// FromInts converts []int to Adata's interface{}
func FromInts(val []int) Adata {
	return Adata(val)
}

// ToInts converts Adata's interface{} to []int
func ToInts(val interface{}) []int {
	return val.([]int)
}

// float64s ////////////////////////////////////////////////////////////////////////////////////////

// FromFloat64 converts float64 to Adata's interface{}
func FromFloat64(val float64) Adata {
	return Adata(val)
}

// ToFloat64 converts Adata's interface{} to float64
func ToFloat64(val interface{}) float64 {
	return val.(float64)
}

// FromFloat64s converts []float64 to Adata's interface{}
func FromFloat64s(val []float64) Adata {
	return Adata(val)
}

// ToFloat64s converts Adata's interface{} to []float64
func ToFloat64s(val interface{}) []float64 {
	return val.([]float64)
}

// bools ////////////////////////////////////////////////////////////////////////////////////////

// FromBool converts bool to Adata's interface{}
func FromBool(val bool) Adata {
	return Adata(val)
}

// ToBool converts Adata's interface{} to bool
func ToBool(val interface{}) bool {
	return val.(bool)
}

// FromBools converts []bool to Adata's interface{}
func FromBools(val []bool) Adata {
	return Adata(val)
}

// ToBools converts Adata's interface{} to []bool
func ToBools(val interface{}) []bool {
	return val.([]bool)
}
