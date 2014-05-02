package core

import (
	"fmt"
	"strconv"
)

func sAdd(i In, o Out) { o(i().(int) + i().(int)) }
func sMul(i In, o Out) { o(i().(int) * i().(int)) }

func sSub(i In, o Out) {
	a := i().(int)
	b := i().(int)
	o(b - a)
}
func sDiv(i In, o Out) {
	a := i().(int)
	b := i().(int)
	o(b / a)
}

func sEq(i In, o Out) {
	o(i() == i())
}
func sLt(i In, o Out) {
	o(i().(int) > i().(int))
}
func sGt(i In, o Out) {
	o(i().(int) < i().(int))
}

func sMap(i In, o Out) {
	fnList := i().([]Value)
	list := i().([]Value)

	fullFnList := make([]Value, len(fnList)+1)
	copy(fullFnList[1:], fnList)

	for i := range list {
		fullFnList[0] = list[i]
		result, _ := Run(fullFnList)
		list[i] = result[0]
	}

	o(list)
}

func sEval(i In, o Out) {
	code := i().([]Value)
	result, _ := Run(code)
	for _, value := range result {
		o(value)
	}
}

func sDup(i In, o Out) {
	v := i()
	o(v)
	o(v)
}

func sIf(i In, o Out) {
	else_ := i().([]Value)
	then := i().([]Value)
	value := i()

	var condition bool
	switch value.(type) {
	case int:
		condition = value != 0
	case string:
		condition = value != ""
	case bool:
		condition = value.(bool)
	default:
		fmt.Println("Ops:", value)
	}

	if condition {
		pushAndRun(i, o, then)
	} else {
		pushAndRun(i, o, else_)
	}
}

func sPrint(i In, o Out) {
	v := i()
	fmt.Println(v)
}

func sNumber(i In, o Out) {
	result, _ := strconv.Atoi(i().(string))
	o(result)
}

func sStr(i In, o Out) {
	result := strconv.Itoa(i().(int))
	o(result)
}

func sLen(i In, o Out) {
	value := i()

	switch value.(type) {
	case string:
		o(len(value.(string)))
	case []Value:
		o(len(value.([]Value)))
	}
}

func sIndex(i In, o Out) {
	index := i().(int)
	value := i()

	switch value.(type) {
	case string:
		char := value.(string)[index]
		o(string(char))
	case []Value:
		o(value.([]Value)[index])
	}
}

func sTranspose(i In, o Out) {
	value1 := i()
	value2 := i()
	o(value1)
	o(value2)
}

// Using a global variable here is actually useful because it allows the
// REPL to remember previously declared functions.
var declared = map[string][]Value{}
func sCall(i In, o Out) {
	name := i().(string)
	body := declared[name]
	pushAndRun(i, o, body)
}
func sDecl(i In, o Out) {
	name := i().(string)
	body := i().([]Value)
	declared[name] = body
}

func sAppend(i In, o Out) {
	value := i()
	list := i().([]Value)
	o(append(list, value))
}

var ops = map[string]func (In, Out){
	"print": sPrint,
	"len": sLen,
	"append": sAppend,

	"declare": sDecl,
	"call": sCall,

	"map": sMap,
	"if": sIf,
	"eval": sEval,
	"dup": sDup,

	"int": sNumber,
	"str": sStr,

	"matches": sMatches,
	"contains": sContains,
	"find": sFind,
	"findall": sFindAll,
	"replace": sReplace,

	"input": sInput,
	"read": sRead,
	"write": sWrite,
	"delete": sDelete,

	"get": sGet,
	"download": sDownload,
	"index": sIndex,
	
	"transpose": sTranspose,
}
