// argparser Project
// Copyright (C) 2021-2022 ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package argparser

type FlagType uint8

// Flag is the options passed along with the commands
// by users. they should send them with prefix "--",
// but we will remove them in the pTools.
type Flag struct {
	name   string
	index  int
	value  interface{}
	fType  FlagType
	emptyT bool
}

type EventArgs struct {
	options    *ParseOptions
	command    string // command without '/' or '!'
	flags      []Flag
	rawData    string
	firstValue string
}

type ParseOptions struct {
	Prefixes []rune
}
