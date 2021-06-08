// argparser Project
// Copyright (C) 2021 wotoTeam, ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package interfaces

// QString is a rich interface for doing string operations.
// it has usefull and (of course) safe method to do special
// operations on a string value.
type QString interface {
	Length() int
	IsEmpty() bool
	GetValue() string
	GetIndexV(int) rune
	IsEqual(QString) bool
	Split(...QString) []QString
	SplitN(int, ...QString) []QString
	SplitFirst(qs ...QString) []QString
	SplitStr(...string) []QString
	SplitStrN(int, ...string) []QString
	SplitStrFirst(...string) []QString
	Contains(...QString) bool
	ContainsStr(...string) bool
	ContainsAll(...QString) bool
	ContainsStrAll(...string) bool
	TrimPrefix(...QString) QString
	TrimPrefixStr(...string) QString
	TrimSuffix(...QString) QString
	TrimSuffixStr(...string) QString
	Trim(qs ...QString) QString
	TrimStr(...string) QString
	Replace(qs, newS QString) QString
	ReplaceStr(qs, newS string) QString
	LockSpecial()
	UnlockSpecial()
}
