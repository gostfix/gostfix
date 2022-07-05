package util

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zaptest/observer"
)

func TestMacExpand(t *testing.T) {

	observedZapCore, observedLogs := observer.New(zap.InfoLevel)
	observedLogger := zap.New(observedZapCore)
	SetLogger(observedLogger)
	_ = observedLogs

	table := map[string]string{
		"name1": "name1-value",
	}

	testCases := []struct {
		pattern string
		status  int
		result  string
		logs    []string
	}{
		{pattern: "$name1", status: 0, result: "name1-value"},
		{pattern: "$(name1", status: 1, result: "", logs: []string{"truncated macro reference"}},
		{pattern: "$(name1)", status: 0, result: "name1-value"},
		{pattern: "$( name1)", status: 0, result: "name1-value"},
		{pattern: "$(name1 )", status: 0, result: "name1-value"},
		{pattern: "$(na me1)", status: 1, result: "", logs: []string{"attribute name syntax error at: \"...na>>> me1\""}},
		{pattern: "${na me1}", status: 1, result: "", logs: []string{"attribute name syntax error at: \"...na>>> me1\""}},
		{pattern: "${${name1} != {}?name 1 defined, |$name1|$name2|}", status: 1, result: "", logs: []string{"attribute name syntax error at: \"...>>>${name1} != {}?name \""}},
		{pattern: "${ ${name1} != {}?name 1 defined, |$name1|$name2|}", status: 1, result: "", logs: []string{"attribute name syntax error at: \"... >>>${name1} != {}?name \""}},
		{pattern: "${ ${name1} ?name 1 defined, |$name1|$name2|}", status: 1, result: "", logs: []string{"attribute name syntax error at: \"... >>>${name1} ?name 1 def\""}},
		{pattern: "${{$name1} ? {name 1 defined, |$name1|$name2|} : {name 1 undefined, |$name1|$name2|} }", status: 1, result: "", logs: []string{"\"==\" or \"!=\"\" or \"<\"\" or \"<=\"\" or \">=\"\" or \">\" expected at: \"...$name1}>>>? {name 1 defined, |$\""}},
		{pattern: "${x{$name1} != {}?{name 1 defined, |$name1|$name2|}}", status: 1, result: "", logs: []string{"attribute name syntax error at: \"...x>>>{$name1} != {}?{name\""}},
		{pattern: "${{$name1}x?{name 1 defined, |$name1|$name2|}}", status: 1, result: "", logs: []string{"\"==\" or \"!=\"\" or \"<\"\" or \"<=\"\" or \">=\"\" or \">\" expected at: \"...$name1}>>>x?{name 1 defined, |\""}},
		{pattern: "${{$name1} != {}x{name 1 defined, |$name1|$name2|}}", status: 1, result: "", logs: []string{"\"?\" or \":\" expected at: \"...}>>>x{name 1 defined, |$\""}},
		{pattern: "${{$name1} != {}?x{name 1 defined, |$name1|$name2|}}", status: 2, result: "x{name 1 defined, |name1-value||}"},
		{pattern: "${{$name2} != {}?x{name 2 defined, |$name1|$name2|}:{name 2 undefined, |$name1|$name2|}}", status: 2, result: ""},
		{pattern: "${{$name1} != {}?{name 1 defined, |$name1|$name2|}x}", status: 3, result: "name 1 defined, |name1-value||", logs: []string{"\":\" expected at: \"...name 1 defined, |$name1|$name2|}>>>x\""}},
		{pattern: "${{$name1} != {}?{name 1 defined, |$name1|$name2|}x:{name 1 undefined, |$name1|$name2|}}", status: 3, result: "name 1 defined, |name1-value||", logs: []string{"\":\" expected at: \"...name 1 defined, |$name1|$name2|}>>>x:{name 1 undefined,\""}},
		{pattern: "${{$name1} != {}?{name 1 defined, |$name1|$name2|}:x{name 1 undefined, |$name1|$name2|}}", status: 2, result: "name 1 defined, |name1-value||"},
		{pattern: "${{$name2} != {}?{name 2 defined, |$name1|$name2|}:x{name 2 undefined, |$name1|$name2|}}", status: 2, result: "x{name 2 undefined, |name1-value||}"},
		{pattern: "${{text}}", status: 1, result: "", logs: []string{"\"==\" or \"!=\"\" or \"<\"\" or \"<=\"\" or \">=\"\" or \">\" expected at: \"...text}>>>\""}},
		{pattern: "${{text}?{non-empty}:{empty}}", status: 1, result: "", logs: []string{"\"==\" or \"!=\"\" or \"<\"\" or \"<=\"\" or \">=\"\" or \">\" expected at: \"...text}>>>?{non-empty}:{empty}\""}},
		{pattern: "${{text} = {}}", status: 1, result: "", logs: []string{"\"==\" or \"!=\"\" or \"<\"\" or \"<=\"\" or \">=\"\" or \">\" expected at: \"...text}>>>= {}\""}},
		{pattern: "${{${ name1}} == {}}", status: 0, result: ""},
		{pattern: "${name1?{${ name1}}:{${name2}}}", status: 0, result: "name1-value"},
		{pattern: "${name2?{${ name1}}:{${name2}}}", status: 2, result: ""},
		{pattern: "${name2?{${name1}}:{${ name2}}}", status: 2, result: ""},
		{pattern: "${name2:{${name1}}:{${name2}}}", status: 1, result: "name1-value", logs: []string{"unexpected input at: \"...${name1}}>>>:{${name2}}\""}},
		{pattern: "${name2?{${name1}}?{${name2}}}", status: 1, result: "", logs: []string{"\":\" expected at: \"...${name1}}>>>?{${name2}}\""}},
		{pattern: "${{${name1?bug:test}} != {bug:test}?{Error: NOT}:{Good:}} Postfix 2.11 compatible", status: 0, result: "Good: Postfix 2.11 compatible"},
		{pattern: "${{${name1??bug}} != {?bug}?{Error: NOT}:{Good:}} Postfix 2.11 compatible", status: 0, result: "Good: Postfix 2.11 compatible"},
		{pattern: "${{${name2::bug}} != {:bug}?{Error: NOT}:{Good:}} Postfix 2.11 compatible", status: 0, result: "Good: Postfix 2.11 compatible"},
		{pattern: "${{xx}==(yy)?{oops}:{phew}}", status: 1, result: "", logs: []string{"\"{expression}\" expected at: \"...{xx} ==>>>(yy)?{oops}:{phew}\""}},
	}

	// MsgVerbose = 2
	for _, tC := range testCases {
		t.Run(tC.pattern, func(t *testing.T) {
			result := strings.Builder{}
			status := MacExpand(&result, tC.pattern, MAC_EXP_FLAG_NONE, nil, func(key string, _ int, context interface{}) (string, error) {
				htable := context.(map[string]string)
				val, exists := htable[key]
				if exists {
					return val, nil
				}
				return "", fmt.Errorf("missing key `%s` in table", key)
			}, table)
			assert.Equal(t, tC.status, status)
			assert.Equal(t, tC.result, result.String())
			logs := observedLogs.TakeAll()
			assert.Equal(t, len(tC.logs), len(logs))
			for i := range tC.logs {
				assert.Equal(t, tC.logs[i], logs[i].Message)
			}
		})
	}
}

func TestMacExpand2(t *testing.T) {

	table := map[string]string{
		"name1": "name1-value",
	}

	testCases := []struct {
		pattern string
		status  int
		result  string
	}{
		{pattern: "${name1?name 1 defined, |$name1|$name2|}", status: 2, result: "name 1 defined, |name1-value||"},
		{pattern: "${name1:name 1 undefined, |$name1|$name2|}", status: 0, result: ""},
		{pattern: "${name2?name 2 defined, |$name1|$name2|}", status: 0, result: ""},
		{pattern: "${name2:name 2 undefined, |$name1|$name2|}", status: 2, result: "name 2 undefined, |name1-value||"},
		{pattern: "|$name1|$name2|", status: 2, result: "|name1-value||"},
		{pattern: "${{$name1} != {}?{name 1 defined, |$name1|$name2|}}", status: 2, result: "name 1 defined, |name1-value||"},
		{pattern: "${{$name1} != {}:{name 1 undefined, |$name1|$name2|}}", status: 0, result: ""},
		{pattern: "${{$name1} == {}?{name 1 undefined, |$name1|$name2|}}", status: 0, result: ""},
		{pattern: "${{$name1} == {}:{name 1 defined, |$name1|$name2|}}", status: 2, result: "name 1 defined, |name1-value||"},
		{pattern: "${name1?{name 1 defined, |$name1|$name2|}:{name 1 undefined, |$name1|$name2|}}", status: 2, result: "name 1 defined, |name1-value||"},
		{pattern: "${{$name1} != {}?{name 1 defined, |$name1|$name2|}:{name 1 undefined, |$name1|$name2|}}", status: 2, result: "name 1 defined, |name1-value||"},
		{pattern: "${{$name1} != {} ? {name 1 defined, |$name1|$name2|} : {name 1 undefined, |$name1|$name2|}}", status: 2, result: "name 1 defined, |name1-value||"},
		{pattern: "${{$name1} != {}?{name 1 defined, |$name1|$name2|}:name 1 undefined, |$name1|$name2|}", status: 2, result: "name 1 defined, |name1-value||"},
		{pattern: "${{$name1} != {} ? {name 1 defined, |$name1|$name2|} : name 1 undefined, |$name1|$name2|}", status: 2, result: "name 1 defined, |name1-value||"},
		{pattern: "${{$name1} != {}}", status: 0, result: "true"},
		{pattern: "${{$name1} == {}}", status: 0, result: ""},
		{pattern: "${{$name2} != {}?{name 2 defined, |$name1|$name2|}}", status: 2, result: ""},
		{pattern: "${{$name2} != {}:{name 2 undefined, |$name1|$name2|}}", status: 2, result: "name 2 undefined, |name1-value||"},
		{pattern: "${{$name2} == {}?{name 2 undefined, |$name1|$name2|}}", status: 2, result: "name 2 undefined, |name1-value||"},
		{pattern: "${{$name2} == {}:{name 2 defined, |$name1|$name2|}}", status: 2, result: ""},
		{pattern: "${name2?{name 2 defined, |$name1|$name2|}:{name 2 undefined, |$name1|$name2|}}", status: 2, result: "name 2 undefined, |name1-value||"},
		{pattern: "${{$name2} != {}?{name 2 defined, |$name1|$name2|}:{name 2 undefined, |$name1|$name2|}}", status: 2, result: "name 2 undefined, |name1-value||"},
		{pattern: "${{$name2} != {} ? {name 2 defined, |$name1|$name2|} : {name 2 undefined, |$name1|$name2|}}", status: 2, result: "name 2 undefined, |name1-value||"},
		{pattern: "${{$name2} != {}?{name 2 defined, |$name1|$name2|}:name 2 undefined, |$name1|$name2|}", status: 2, result: "name 2 undefined, |name1-value||"},
		{pattern: "${{$name2} != {} ? {name 2 defined, |$name1|$name2|} : name 2 undefined, |$name1|$name2|}", status: 2, result: "name 2 undefined, |name1-value||"},
		{pattern: "${{$name2} != {}}", status: 2, result: ""},
		{pattern: "${{$name2} == {}}", status: 2, result: "true"},
	}

	for _, tC := range testCases {
		t.Run(tC.pattern, func(t *testing.T) {
			result := strings.Builder{}
			status := MacExpand(&result, tC.pattern, MAC_EXP_FLAG_NONE, nil, func(key string, _ int, context interface{}) (string, error) {
				htable := context.(map[string]string)
				val, exists := htable[key]
				if exists {
					return val, nil
				}
				return "", fmt.Errorf("missing key `%s` in table", key)
			}, table)
			assert.Equal(t, tC.status, status)
			assert.Equal(t, tC.result, result.String())
		})
	}
}

func TestMacExpand3(t *testing.T) {

	table := map[string]string{}

	testCases := []struct {
		pattern string
		status  int
		result  string
	}{
		{pattern: "${{1} == {1}}", status: 0, result: "true"},
		{pattern: "${{1} <  {1}}", status: 0, result: ""},
		{pattern: "${{1} <= {1}}", status: 0, result: "true"},
		{pattern: "${{1} >= {1}}", status: 0, result: "true"},
		{pattern: "${{1} >  {1}}", status: 0, result: ""},
		{pattern: "${{1} == {2}}", status: 0, result: ""},
		{pattern: "${{1} <  {2}}", status: 0, result: "true"},
		{pattern: "${{1} <= {2}}", status: 0, result: "true"},
		{pattern: "${{1} >= {2}}", status: 0, result: ""},
		{pattern: "${{1} >  {2}}", status: 0, result: ""},
		{pattern: "${{a} == {a}}", status: 0, result: "true"},
		{pattern: "${{a} <  {a}}", status: 0, result: ""},
		{pattern: "${{a} <= {a}}", status: 0, result: "true"},
		{pattern: "${{a} >= {a}}", status: 0, result: "true"},
		{pattern: "${{a} >  {a}}", status: 0, result: ""},
		{pattern: "${{a} == {b}}", status: 0, result: ""},
		{pattern: "${{a} <  {b}}", status: 0, result: "true"},
		{pattern: "${{a} <= {b}}", status: 0, result: "true"},
		{pattern: "${{a} >= {b}}", status: 0, result: ""},
		{pattern: "${{a} >  {b}}", status: 0, result: ""},
	}

	for _, tC := range testCases {
		t.Run(tC.pattern, func(t *testing.T) {
			result := strings.Builder{}
			status := MacExpand(&result, tC.pattern, MAC_EXP_FLAG_NONE, nil, func(key string, _ int, context interface{}) (string, error) {
				htable := context.(map[string]string)
				val, exists := htable[key]
				if exists {
					return val, nil
				}
				return "", fmt.Errorf("missing key `%s` in table", key)
			}, table)
			assert.Equal(t, tC.status, status)
			assert.Equal(t, tC.result, result.String())
		})
	}
}

func length_relop_eval(left string, relop int, rite string) int {
	var myname string = "length_relop_eval"
	var delta int = len(left) - len(rite)

	switch relop {
	case MAC_EXP_OP_TOK_EQ:
		return (mac_exp_op_res_bool[delta == 0])
	case MAC_EXP_OP_TOK_NE:
		return (mac_exp_op_res_bool[delta != 0])
	case MAC_EXP_OP_TOK_LT:
		return (mac_exp_op_res_bool[delta < 0])
	case MAC_EXP_OP_TOK_LE:
		return (mac_exp_op_res_bool[delta <= 0])
	case MAC_EXP_OP_TOK_GE:
		return (mac_exp_op_res_bool[delta >= 0])
	case MAC_EXP_OP_TOK_GT:
		return (mac_exp_op_res_bool[delta > 0])
	default:
		MsgPanic("unknown operator", "function", myname, "relop", relop)
	}
	return 0
}

func TestMacExpand4(t *testing.T) {
	var length_relops = []int{
		MAC_EXP_OP_TOK_EQ, MAC_EXP_OP_TOK_NE,
		MAC_EXP_OP_TOK_GT, MAC_EXP_OP_TOK_GE,
		MAC_EXP_OP_TOK_LT, MAC_EXP_OP_TOK_LE,
	}
	mac_expand_add_relop(length_relops, "length", length_relop_eval)

	table := map[string]string{
		"name1": "foo",
	}

	testCases := []struct {
		pattern string
		status  int
		result  string
	}{
		{pattern: "${{$name1} >=blah {bar}}", status: 1, result: ""},
		{pattern: "${{aaa} == {bbb}}", status: 0, result: ""},
		{pattern: "${{aaa} ==length {bbb}}", status: 0, result: "true"},
		{pattern: "${{aaa} <=length {bbb}}", status: 0, result: "true"},
		{pattern: "${{aaa} >=length {bbb}}", status: 0, result: "true"},
		{pattern: "${{aaa} != {bbb}}", status: 0, result: "true"},
		{pattern: "${{aaa} !=length {bbb}}", status: 0, result: ""},
		{pattern: "${{aaa} > {bb}}", status: 0, result: ""},
		{pattern: "${{aaa} >length {bb}}", status: 0, result: "true"},
		{pattern: "${{aaa} >= {bb}}", status: 0, result: ""},
		{pattern: "${{aaa} >=length {bb}}", status: 0, result: "true"},
		{pattern: "${{aaa} < {bb}}", status: 0, result: "true"},
		{pattern: "${{aaa} <length {bb}}", status: 0, result: ""},
		{pattern: "${{aaa} <= {bb}}", status: 0, result: "true"},
		{pattern: "${{aaa} <=length {bb}}", status: 0, result: ""},
	}

	for _, tC := range testCases {
		t.Run(tC.pattern, func(t *testing.T) {
			result := strings.Builder{}
			status := MacExpand(&result, tC.pattern, MAC_EXP_FLAG_NONE, nil, func(key string, _ int, context interface{}) (string, error) {
				htable := context.(map[string]string)
				val, exists := htable[key]
				if exists {
					return val, nil
				}
				return "", fmt.Errorf("missing key `%s` in table", key)
			}, table)
			assert.Equal(t, tC.status, status)
			assert.Equal(t, tC.result, result.String())
		})
	}
}
