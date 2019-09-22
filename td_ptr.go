// Copyright (c) 2018, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

package testdeep

import (
	"reflect"

	"github.com/maxatome/go-testdeep/internal/ctxerr"
	"github.com/maxatome/go-testdeep/internal/types"
)

type tdPtr struct {
	tdSmugglerBase
}

var _ TestDeep = &tdPtr{}

// Ptr is a smuggler operator. It takes the address of data and
// compares it to "val".
//
// "val" depends on data type. For example, if the compared data is an
// *int, one can have:
//   Ptr(12)
// as well as an other operator:
//   Ptr(Between(3, 4))
//
// TypeBehind method returns the reflect.Type of a pointer on "val",
// except if "val" is a TestDeep operator. In this case, it delegates
// TypeBehind() to the operator and returns the reflect.Type of a
// pointer on the returned value (if non-nil of course).
func Ptr(val interface{}) TestDeep {
	vval := reflect.ValueOf(val)
	if vval.IsValid() {
		p := tdPtr{
			tdSmugglerBase: newSmugglerBase(val),
		}

		if !p.isTestDeeper {
			p.expectedValue = reflect.New(vval.Type())
			p.expectedValue.Elem().Set(vval)
		}
		return &p
	}
	panic("usage: Ptr(NON_NIL_VALUE)")
}

func (p *tdPtr) Match(ctx ctxerr.Context, got reflect.Value) *ctxerr.Error {
	if got.Kind() != reflect.Ptr {
		if ctx.BooleanError {
			return ctxerr.BooleanError
		}
		return ctx.CollectError(&ctxerr.Error{
			Message:  "pointer type mismatch",
			Got:      types.RawString(got.Type().String()),
			Expected: types.RawString(p.String()),
		})
	}

	if p.isTestDeeper {
		return deepValueEqual(ctx.AddPtr(1), got.Elem(), p.expectedValue)
	}
	return deepValueEqual(ctx, got, p.expectedValue)
}

func (p *tdPtr) String() string {
	if p.isTestDeeper {
		return "*<something>"
	}
	return p.expectedValue.Type().String()
}

func (p *tdPtr) TypeBehind() reflect.Type {
	// If the expected value is a TestDeep operator, delegate TypeBehind to it
	if p.isTestDeeper {
		typ := p.expectedValue.Interface().(TestDeep).TypeBehind()
		if typ == nil {
			return nil
		}
		// Add a level of pointer
		return reflect.New(typ).Type()
	}
	return p.expectedValue.Type()
}

type tdPPtr struct {
	tdSmugglerBase
}

var _ TestDeep = &tdPPtr{}

// PPtr is a smuggler operator. It takes the address of the address of
// data and compares it to "val".
//
// "val" depends on data type. For example, if the compared data is an
// **int, one can have:
//   PPtr(12)
// as well as an other operator:
//   PPtr(Between(3, 4))
//
// It is more efficient and shorter to write than:
//   Ptr(Ptr(val))
//
// TypeBehind method returns the reflect.Type of a pointer on a
// pointer on "val", except if "val" is a TestDeep operator. In this
// case, it delegates TypeBehind() to the operator and returns the
// reflect.Type of a pointer on a pointer on the returned value (if
// non-nil of course).
func PPtr(val interface{}) TestDeep {
	vval := reflect.ValueOf(val)
	if vval.IsValid() {
		p := tdPPtr{
			tdSmugglerBase: newSmugglerBase(val),
		}

		if !p.isTestDeeper {
			pVval := reflect.New(vval.Type())
			pVval.Elem().Set(vval)

			p.expectedValue = reflect.New(pVval.Type())
			p.expectedValue.Elem().Set(pVval)
		}
		return &p
	}
	panic("usage: PPtr(NON_NIL_VALUE)")
}

func (p *tdPPtr) Match(ctx ctxerr.Context, got reflect.Value) *ctxerr.Error {
	if got.Kind() != reflect.Ptr || got.Elem().Kind() != reflect.Ptr {
		if ctx.BooleanError {
			return ctxerr.BooleanError
		}
		return ctx.CollectError(&ctxerr.Error{
			Message:  "pointer type mismatch",
			Got:      types.RawString(got.Type().String()),
			Expected: types.RawString(p.String()),
		})
	}

	if p.isTestDeeper {
		return deepValueEqual(ctx.AddPtr(2), got.Elem().Elem(), p.expectedValue)
	}
	return deepValueEqual(ctx, got, p.expectedValue)
}

func (p *tdPPtr) String() string {
	if p.isTestDeeper {
		return "**<something>"
	}
	return p.expectedValue.Type().String()
}

func (p *tdPPtr) TypeBehind() reflect.Type {
	// If the expected value is a TestDeep operator, delegate TypeBehind to it
	if p.isTestDeeper {
		typ := p.expectedValue.Interface().(TestDeep).TypeBehind()
		if typ == nil {
			return nil
		}
		// Add 2 levels of pointer
		return reflect.New(reflect.New(typ).Type()).Type()
	}
	return p.expectedValue.Type()
}
