// argparser Project
// Copyright (C) 2021-2022 ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package argparser

import (
	"strconv"
	"strings"

	ws "github.com/ALiwoto/ssg/ssg"
)

//---------------------------------------------------------

func (f *Flag) setName(name string) {
	name = strings.TrimSpace(name)
	f.name = name
}

func (f *Flag) setNameQ(q ws.QString) {
	f.setName(q.GetValue())
}

func (f *Flag) setNilValue() {
	f.value = nil
	f.fType = NoneFlagType
}

func (f *Flag) setRealValue(value string) {
	tmp := strings.TrimSpace(value)
	if len(tmp) == ws.BaseIndex {
		f.setAsBool(true)
		f.setAsEmpty()
		return
	}

	myI, err := strconv.ParseInt(value, ws.BaseTen, ws.Base64Bit)
	if err == nil {
		f.setAsInt(&myI)
		return
	}

	if tmp[ws.BaseIndex] == ws.CHAR_STR {
		myInt := len(tmp) - ws.BaseOneIndex
		if tmp[myInt] == ws.CHAR_STR {
			tmp = strings.TrimPrefix(tmp, ws.STR_SIGN)
			tmp = strings.TrimSuffix(tmp, ws.STR_SIGN)
			f.setAsString(&tmp)
			return
		}
	}

	v, isBool := ToBoolType(value)
	if isBool {
		f.setAsBool(v)
		return
	}

	f.setAsString(&value)
}

func (f *Flag) setAsBool(value bool) {
	f.value = value
	f.fType = BoolFlagType
}

func (f *Flag) setAsEmpty() {
	f.emptyT = true
}

func (f *Flag) setAsString(value *string) {
	f.value = *value
	f.fType = StringFlagType
}

func (f *Flag) setAsInt(value *int64) {
	f.value = *value
	f.fType = Int64FlagType
}

// GetValue will return you the value of this flag.
// remember that value can be
func (f *Flag) GetValue() interface{} {
	return f.value
}

// GetType will return you the type of
// the value this flag. it's an enum.
func (f *Flag) GetType() FlagType {
	return f.fType
}

// GetIndex will return you the index of this flag.
// index of a flag is a unique int.
// even if the name of two flags are the same, their
// index will not be the same.
func (f *Flag) GetIndex() int {
	return f.index
}

// GetName will give you the name of this flag.
// it's not unique actually.
// for example:
//   `/command --test "hello" --test = "HI!"
func (f *Flag) GetName() string {
	return f.name
}

// GetValueAndType returns both value and type of this
// flag. originally it has internal usage, but maybe
// you want to use it in your own package, so I made it
// public!
func (f *Flag) GetValueAndType() (interface{}, FlagType) {
	return f.GetValue(), f.GetType()
}

// GetAsString will give you the value as an string.
// for example if value is `true` (with flag type of
// `bool`), then the string will be "true".
// or if it's 10(an integer), it will give you: "10".
func (f *Flag) GetAsString() string {
	v, t := f.GetValueAndType()
	if t == StringFlagType {
		vS, ok := v.(string)
		if !ok {
			return NoneTypeStr
		}

		return vS
	}

	if t == BoolFlagType {
		vB, ok := v.(bool)
		if !ok {
			return NoneTypeStr
		}

		if vB {
			return TrueHlc
		} else {
			return FalseHlc
		}
	}

	if t.IsInteger() {
		vI, ok := v.(int64)
		if !ok {
			return NoneTypeStr
		}

		return strconv.FormatInt(vI, ws.BaseTen)
	}

	return NoneTypeStr
}

func (f *Flag) isEmpty() bool {
	return f.emptyT || f.fType == NoneFlagType
}

func (f *Flag) GetAsInteger() (vI int64, ok bool) {
	v, t := f.GetValueAndType()
	if t == StringFlagType {
		vS, ok := v.(string)
		if !ok {
			return ws.BaseIndex, false
		}

		vS = strings.ReplaceAll(vS, ws.SPACE_VALUE, ws.EMPTY)
		vS = strings.ReplaceAll(vS, ws.STR_SIGN, ws.EMPTY)
		vI, err := strconv.ParseInt(vS, ws.BaseTen, ws.Base64Bit)
		if err != nil {
			return ws.BaseIndex, false
		}

		return vI, true
	}

	if t == BoolFlagType {
		vB, ok := v.(bool)
		if !ok {
			return ws.BaseIndex, false
		}

		if vB {
			return ws.BaseOneIndex, true
		} else {
			return ws.BaseIndex, true
		}
	}

	if t.IsInteger() {
		vI, ok := v.(int64)
		if !ok {
			return ws.BaseIndex, false
		}

		return vI, true
	}

	return ws.BaseIndex, false
}

func (f *Flag) GetAsBool() bool {
	v, t := f.GetValueAndType()
	if t == StringFlagType {
		vS, ok := v.(string)
		if !ok {
			return false
		}

		vS = strings.ReplaceAll(vS, ws.SPACE_VALUE, ws.EMPTY)
		vS = strings.ReplaceAll(vS, ws.STR_SIGN, ws.EMPTY)
		vB, ok := ToBoolType(vS)
		if !ok {
			return false
		}

		return vB
	}

	if t == BoolFlagType {
		vB, ok := v.(bool)
		if !ok {
			return false
		}

		return vB
	}

	if t.IsInteger() {
		vI, ok := v.(int64)
		if !ok {
			return false
		}

		return vI != ws.BaseIndex
	}

	return false
}

//---------------------------------------------------------
func (e *EventArgs) GetCommand() string {
	return e.command
}

// CompareCommand will compare the command of this event arg
// with the given command. This function is case sensitive,
// if you want to compare it case insensitive, then use
// CheckCommand method.
func (e *EventArgs) CompareCommand(cmd string) bool {
	cmd = strings.TrimLeft(cmd, ws.GET_SLASH)
	return e.command == cmd
}

// CheckCommand will compare the command of this event arg
// with the given command. This function is case insensitive,
// if you want to compare it case sensitive, then use
// CompareCommand method.
func (e *EventArgs) CheckCommand(cmd string) bool {
	cmd = strings.TrimLeft(cmd, ws.GET_SLASH)
	return strings.EqualFold(e.command, cmd)
}

func (e *EventArgs) GetFlags() []Flag {
	return e.flags
}

func (e *EventArgs) HasFlag(names ...string) bool {
	for _, current := range e.flags {
		for _, name := range names {
			if strings.EqualFold(current.name, name) {
				return true
			}
		}
	}

	return false
}

func (e *EventArgs) HasFlags(names ...string) bool {
	for _, current := range e.flags {
		for _, name := range names {
			if !strings.EqualFold(current.name, name) {
				return false
			}
		}
	}

	return true
}

func (e *EventArgs) HasRawData() bool {
	return !ws.IsEmpty(&e.rawData)
}

func (e *EventArgs) GetIndexFlag(index int) *Flag {
	if index < ws.BaseIndex || index >= e.GetLength() {
		return nil
	}

	return &e.flags[index]
}

func (e *EventArgs) GetFlag(names ...string) *Flag {
	for _, current := range e.flags {
		for _, name := range names {
			if strings.EqualFold(current.name, name) {
				return &current
			}
		}
	}

	return nil
}

func (e *EventArgs) GetLength() int {
	return len(e.flags)
}

func (e *EventArgs) IsEmpty() bool {
	return len(e.flags) == ws.BaseIndex
}

func (e *EventArgs) IsEmptyOrRaw() bool {
	return e.IsEmpty() && len(e.rawData) == ws.BaseIndex
}

// GetAsString will give you the string value of the flag
// with the specified name.
// if there is no flag with this name, or there is an
// error on our path, it will return you an empty string.
func (e *EventArgs) GetAsString(name ...string) string {
	f := e.GetFlag(name...)
	if f == nil {
		return ws.EMPTY
	}

	return f.GetAsString()
}

// GetAsStringOrRaw will give you the string value of the flag
// with the specified name.
// if there is no flag with this name, or there is an
// error on our path, it will return you the raw data in the
// EventArgs struct.
func (e *EventArgs) GetAsStringOrRaw(name ...string) string {
	tmp := e.GetAsString(name...)

	if ws.IsEmpty(&tmp) {
		return e.rawData
	}

	return tmp
}

// GetAsString will give you the integer value of the flag
// with the specified name.
// if there is no flag with this name, or there is an
// error on our path, it will return you zero and false.
func (e *EventArgs) GetAsInteger(name ...string) (vI int64, ok bool) {
	f := e.GetFlag(name...)
	if f == nil {
		return ws.BaseIndex, false
	}

	return f.GetAsInteger()
}

// GetAsString will give you the integer value of the flag
// with the specified name.
// if there is no flag with this name, or there is an
// error on our path, it will return you zero and false.
func (e *EventArgs) GetAsIntegerOrRaw(name ...string) (vI int64, ok bool) {
	if e.IsEmpty() || name == nil {
		str := e.GetAsStringOrRaw()
		vI, err := strconv.ParseInt(str, ws.BaseTen, ws.Base64Bit)
		if err != nil {
			return ws.BaseIndex, false
		}

		return vI, true
	} else {
		return e.GetAsInteger(name...)
	}
}

func (e *EventArgs) GetFirstNoneEmptyValue() string {
	for _, current := range e.flags {
		if !current.isEmpty() {
			return current.GetAsString()
		}
	}

	// lets return raw data here; as it is the last resort.
	if e.rawData != "" {
		return e.rawData
	}

	return e.firstValue
}

// GetAsString will give you the boolean value of the flag
// with the specified name.
// if there is no flag with this name, or there is an
// error on our path, it will return you false.
func (e *EventArgs) GetAsBool(name ...string) bool {
	f := e.GetFlag(name...)
	if f == nil {
		return false
	}

	return f.GetAsBool()
}

func (e *EventArgs) setFlags(f []Flag) {
	e.flags = f
}

//---------------------------------------------------------

func (t *FlagType) ToString() string {
	switch *t {
	case NoneFlagType:
		return NoneTypeStr
	case BoolFlagType:
		return BoolTypeStr
	case StringFlagType:
		return StringTypeStr
	case UInt8FlagType:
		return UInt8TypeStr
	case UInt16FlagType:
		return UInt16TypeStr
	case UInt32FlagType:
		return UInt32TypeStr
	case UInt64FlagType:
		return UInt64TypeStr
	case Int8FlagType:
		return Int8TypeStr
	case Int16FlagType:
		return Int16TypeStr
	case Int32FlagType:
		return Int32TypeStr
	case Int64FlagType:
		return Int64TypeStr
	}

	return NoneTypeStr
}

func (t *FlagType) Compare(value *FlagType) bool {
	return false
}

func (t *FlagType) IsInteger() bool {
	return *t == Int8FlagType ||
		*t == Int16FlagType ||
		*t == Int32FlagType ||
		*t == Int64FlagType ||
		*t == UInt8FlagType ||
		*t == UInt16FlagType ||
		*t == UInt32FlagType ||
		*t == UInt64FlagType
}

func (t *FlagType) IsString() bool {
	return *t == StringFlagType
}

//---------------------------------------------------------
