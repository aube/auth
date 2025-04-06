package main

import (
	"fmt"
)

// This is a demonstration of embedding struct and overriding
// a method of parent struct.
// https://go.dev/play/p/nIdwUV8Qlaj

// Trait is a set of methods that a Base struct
// should implement with default value and a custom
// struct that embedded Base struct should implement
// with a custom value rather than default value inherited
// from Base struct.
type Trait interface {
	Name() string
	MyName() string
}

type Base struct {
	trait Trait
}

func NewBase() *Base {
	return new(Base)
}

// Implement Trait interface
// =========================

func (b *Base) Name() string {
	// panic("NotImplemented. It should be implemented by a struct that embeds Base struct.")

	// Or can return a default value. E.g.
	return "Base"
}

func (b *Base) MyName() string {
	if b.trait != nil {
		return b.trait.Name()
	}

	return b.Name()
}

// Custom1 overrides Base.Name() method.
type Custom1 struct {
	*Base
}

// Custom2 does not override Base.Name() method.
type Custom2 struct {
	*Base
}

func NewCustom1() *Custom1 {
	base := NewBase()
	c := &Custom1{base}
	c.trait = c

	return c
}

func NewCustom2() *Custom2 {
	base := NewBase()
	c := &Custom2{base}
	c.trait = c

	return c
}

// Overriding Name() method.
func (c *Custom1) Name() string {
	return "Custom1"
} // Overriding Name() method.
func (c *Custom2) Name() string {
	return "Custom2"
}

func main() {
	c1 := NewCustom1()
	fmt.Println(c1.MyName()) // "Custom1"

	c2 := NewCustom2()
	fmt.Println(c2.MyName()) // panic
}
