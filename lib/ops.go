package lib

import "fmt"

func sAdd(i In, o Out) { o(i().(int) + i().(int)) }
func sSub(i In, o Out) { o(i().(int) - i().(int)) }
func sDiv(i In, o Out) { o(i().(int) / i().(int)) }
func sMul(i In, o Out) { o(i().(int) * i().(int)) }

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
	o(v)
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

var ops = map[string]func (In, Out){
	"eval": sEval,
	"print": sPrint,
}
