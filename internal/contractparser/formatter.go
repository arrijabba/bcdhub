package contractparser

import (
	"fmt"
	"strings"

	"github.com/tidwall/gjson"

	"github.com/aopoltorzhicky/bcdhub/internal/helpers"
)

const lineSize = 100

func isFramed(n gjson.Result) bool {
	prim := n.Get("prim").String()
	if helpers.StringInArray(prim, []string{
		"Pair", "Left", "Right", "Some",
		"pair", "or", "option", "map", "big_map", "list", "set", "contract", "lambda",
	}) {
		return true
	} else if helpers.StringInArray(prim, []string{
		"key", "unit", "signature", "operation",
		"int", "nat", "string", "bytes", "mutez", "bool", "key_hash", "timestamp", "address",
	}) {
		return n.Get("annots").Exists()
	}
	return false
}

func isComplex(n gjson.Result) bool {
	prim := n.Get("prim").String()
	return prim == "LAMBDA" || prim[:2] == "IF"
}

func isInline(n gjson.Result) bool {
	prim := n.Get("prim").String()
	return prim == "PUSH"
}

func isScript(n gjson.Result) bool {
	if !n.IsArray() {
		return false
	}
	for _, item := range n.Array() {
		prim := item.Get("prim").String()
		if !helpers.StringInArray(prim, []string{
			"parameter", "storage", "code",
		}) {
			return false
		}
	}
	return true
}

// MichelineToMichelson -
func MichelineToMichelson(n gjson.Result, inline bool) string {
	return formatNode(n, "", inline, true, false)
}

func formatNode(node gjson.Result, indent string, inline, isRoot, wrapped bool) string {
	if node.IsArray() {
		return formatArray(node, indent, inline, isRoot)
	}

	if node.IsObject() {
		return formatObject(node, indent, inline, isRoot, wrapped)
	}

	fmt.Println("NODE:", node)
	panic("shit happens")
}

func formatArray(node gjson.Result, indent string, inline, isRoot bool) string {
	seqIndent := indent
	isScriptRoot := isRoot && isScript(node)
	if !isScriptRoot {
		seqIndent = indent + "  "
	}

	items := make([]string, len(node.Array()))

	for i, n := range node.Array() {
		items[i] = formatNode(n, seqIndent, inline, false, true)
	}

	if len(items) == 0 {
		return "{}"
	}

	length := len(indent) + 4
	for _, i := range items {
		length += len(i)
	}

	space := ""
	if !isScriptRoot {
		space = " "
	}

	var seq string

	if inline || length < lineSize {
		seq = strings.Join(items, fmt.Sprintf("%v; ", space))
	} else {
		seq = strings.Join(items, fmt.Sprintf("%v;\n%v", space, seqIndent))
	}

	if !isScriptRoot {
		return fmt.Sprintf("{ %v }", seq)
	}

	return seq
}

func formatObject(node gjson.Result, indent string, inline, isRoot, wrapped bool) string {
	if node.Get("prim").Exists() {
		return formatPrimObject(node, indent, inline, isRoot, wrapped)
	}

	return formatNonPrimObject(node)
}

func formatPrimObject(node gjson.Result, indent string, inline, isRoot, wrapped bool) string {
	res := []string{node.Get("prim").String()}

	if annots := node.Get("annots"); annots.Exists() {
		for _, a := range annots.Array() {
			res = append(res, a.String())
		}
	}

	expr := strings.Join(res, " ")

	var args []gjson.Result
	if rawArgs := node.Get("args"); rawArgs.Exists() {
		args = rawArgs.Array()
	}

	if isComplex(node) {
		argIndent := indent + "  "
		items := make([]string, len(args))
		for i, a := range args {
			items[i] = formatNode(a, argIndent, inline, false, false)
		}

		length := len(indent) + len(expr) + len(items) + 1

		for _, item := range items {
			length += len(item)
		}

		if inline || length < lineSize {
			expr = fmt.Sprintf("%v %v", expr, strings.Join(items, " "))
		} else {
			res := []string{expr}
			res = append(res, items...)
			expr = strings.Join(res, fmt.Sprintf("\n%v", argIndent))
		}
	} else if len(args) == 1 {
		argIndent := indent + strings.Repeat(" ", len(expr)+1)
		expr = fmt.Sprintf("%v %v", expr, formatNode(args[0], argIndent, inline, false, false))
	} else if len(args) > 1 {
		argIndent := indent + "  "
		altIndent := indent + strings.Repeat(" ", len(expr)+2)

		for _, arg := range args {
			item := formatNode(arg, argIndent, inline, false, false)
			length := len(indent) + len(expr) + len(item) + 1
			if inline || isInline(node) || length < lineSize {
				argIndent = altIndent
				expr = fmt.Sprintf("%v %v", expr, item)
			} else {
				expr = fmt.Sprintf("%v\n%v%v", expr, argIndent, item)
			}
		}
	}

	if isFramed(node) && !isRoot && !wrapped {
		return fmt.Sprintf("(%v)", expr)
	}
	return expr
}

func formatNonPrimObject(node gjson.Result) string {
	if len(node.Map()) != 1 {
		fmt.Println("NODE:", node)
		panic("node keys count != 1")
	}

	for coreType, value := range node.Map() {
		if coreType == "int" {
			return value.String()
		} else if coreType == "bytes" {
			return fmt.Sprintf("0x%v", value.String())
		} else if coreType == "string" {
			return value.Raw
		}
	}

	fmt.Println("NODE:", node)
	panic("invalid coreType")
}
