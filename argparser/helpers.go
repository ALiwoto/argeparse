// argparser Project
// Copyright (C) 2021-2022 ALiwoto
// This file is subject to the terms and conditions defined in
// file 'LICENSE', which is part of the source code.

package argparser

import (
	"errors"
	"strings"

	ws "github.com/ALiwoto/ssg/ssg"
)

func ParseArg(text string, prefixes []rune) (e *EventArgs, err error) {
	return ParseArgWithOptions(text, &ParseOptions{
		Prefixes: prefixes,
	})
}

// ParseArgWithOptions will parse the whole text into an EventArg and will return it.
func ParseArgWithOptions(text string, options *ParseOptions) (e *EventArgs, err error) {
	if text == "" {
		return nil, errors.New("ParseArgWithOptions: text cannot be empty")
	}

	if options == nil {
		options = GetDefaultParseOptions()
	}

	if len(options.Prefixes) == 0 {
		options.Prefixes = DefaultPrefixes
	}

	ss := ws.Ss(text)
	if !ss.HasRunePrefix(options.Prefixes...) {
		return nil, errors.New("ParseArgWithOptions: input text is not a command")
	}

	cmdR := ss.SplitStr(ws.SPACE_VALUE)
	if len(cmdR) == ws.BaseIndex {
		return nil, errors.New("ParseArgWithOptions: unable to get the command")
	}

	cmd := cmdR[ws.BaseIndex]
	if cmd.IsEmpty() {
		return nil, errors.New("ParseArgWithOptions: length of the command cannot be zero")
	}

	cmdSs := cmd.TrimStr(toStrArray(options.Prefixes)...)
	if cmdSs.IsEmpty() {
		return nil, errors.New("ParseArgWithOptions: command cannot be only whitespace")
	}

	cmdStr := cmdSs.GetValue()

	e = &EventArgs{
		command: cmdStr,
		options: options,
	}

	// lock the special characters such as "--", ":", "=".
	ss.LockSpecial()

	tmpOSs := ss.SplitStr(ws.FLAG_PREFIX)
	// check if we have any flags or not.
	// I think this is not necessary actually,
	// but I just added it to prevent some cases of
	// panics. and also it will reduce the time order
	// I guess.
	if len(tmpOSs) < ws.BaseTwoIndex {
		// please notice that we should send the original
		// text to this function.
		// because our locked QString contains JA characters
		// and should not be used here.
		lookRaw(&text, e)
		return e, nil
	} else {
		tmpFirstValue := tmpOSs[ws.BaseIndex]
		tmpFirstValue.UnlockSpecial()
		firstValue := tmpFirstValue.GetValue()[ws.BaseOneIndex:]
		firstValue = strings.TrimPrefix(firstValue, e.command)
		firstValue = strings.TrimSpace(firstValue)
		e.firstValue = firstValue
	}

	flagsR := tmpOSs[ws.BaseOneIndex:]
	// it means it has no flags available.
	// so return the e.
	if len(flagsR) == ws.BaseIndex {
		// please notice that we should send the original
		// text to this method.
		// because our locked QString contains JA characters
		// and should not be used here.
		lookRaw(&text, e)
		return e, nil
	}

	myFlags := make([]Flag, ws.BaseIndex)
	tmp := ws.EMPTY
	var tmpFlag Flag
	var tmpArr []ws.QString

	for i, current := range flagsR {
		tmpFlag = Flag{
			index: i,
		}

		tmp = ws.EMPTY
		// let me explain you something here.
		// it really does matter how you pass these constants to this functions.
		// first of all should be equal.
		// and then double dot (':')
		// and in the end, it should be space.
		// please don't forget that you should prioritize them.
		tmpArr = current.SplitStrFirst(ws.EqualStr, ws.DdotSign, ws.SPACE_VALUE)

		tmpFlag.setNameQ(tmpArr[ws.BaseIndex])
		if len(tmpArr) == ws.BaseIndex {
			tmpFlag.setNilValue()
			myFlags = append(myFlags, tmpFlag)
			continue
		}

		for i, ar := range tmpArr {
			if i == ws.BaseIndex {
				// ignore first slice, because it's flag name.
				continue
			}

			ar.UnlockSpecial()
			tmp += ar.GetValue()
		}
		tmpFlag.setRealValue(fixTmpStr(tmp))

		myFlags = append(myFlags, tmpFlag)
	}

	e.setFlags(myFlags)

	return e, nil
}

func ParseArgDefault(text string) (e *EventArgs, err error) {
	return ParseArgWithOptions(text, GetDefaultParseOptions())
}

func GetDefaultParseOptions() *ParseOptions {
	return &ParseOptions{
		Prefixes: DefaultPrefixes,
	}
}

func toStrArray(r []rune) []string {
	var s []string
	for _, v := range r {
		s = append(s, string(v))
	}
	s = append(s, ws.SPACE_VALUE)
	return s
}

func fixTmpStr(tmp string) string {
	tmp = strings.TrimSpace(tmp)
	if strings.HasPrefix(tmp, ws.EqualStr) {
		tmp = strings.TrimPrefix(tmp, ws.EqualStr)
		tmp = strings.TrimSpace(tmp)
	}
	tmp = strings.Trim(tmp, ws.STR_SIGN)
	return tmp
}

// look raw will look for raw data.
// please use this function when and only when
// no flags are provided for our commands.
func lookRaw(text *string, e *EventArgs) {
	myStr := strings.SplitN(*text, e.command, ws.BaseTwoIndex)
	if len(myStr) < ws.BaseTwoIndex {
		return
	}

	tmp := strings.Join(myStr[ws.BaseOneIndex:], ws.EMPTY)
	tmp = strings.TrimSpace(tmp)

	e.rawData = tmp
}

func ToBoolType(value string) (v, isBool bool) {
	value = strings.TrimSpace(value)
	value = strings.ToLower(value)
	switch value {
	case TrueHlc, YesHlc, OnHlc:
		return true, true
	case FalseHlc, NoHlc, OffHlc:
		return false, true
	default:
		return false, false
	}
}
