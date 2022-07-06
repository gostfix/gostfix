package util

import (
	"fmt"
	"strconv"
	"strings"
	"text/scanner"
	"unicode"

	"github.com/gostfix/gostfix/internal/ascii"
)

const (
	MAC_EXP_BVAL_TRUE  = "true"
	MAC_EXP_BVAL_FALSE = ""
)

const (
	MAC_EXP_MODE_TEST = 0
	MAC_EXP_MODE_USE  = 1
)

const (
	MAC_EXP_FLAG_NONE      int = (0)
	MAC_EXP_FLAG_RECURSE   int = (1 << 0)
	MAC_EXP_FLAG_APPEND    int = (1 << 1)
	MAC_EXP_FLAG_SCAN      int = (1 << 2)
	MAC_EXP_FLAG_PRINTABLE int = (1 << 3)
)

const (
	MAC_EXP_OP_RES_TRUE  = 0
	MAC_EXP_OP_RES_FALSE = 1
	MAC_EXP_OP_RES_ERROR = 2
) // MAC_EXP_OP_RES

var mac_exp_op_res_bool = map[bool]int{
	false: MAC_EXP_OP_RES_FALSE,
	true:  MAC_EXP_OP_RES_TRUE,
}

/*
 * Token codes, public so tha they are available to mac_expand_add_relop()
 */
const (
	MAC_EXP_OP_TOK_NONE = 0 /* Sentinel */
	MAC_EXP_OP_TOK_EQ   = 1 /* == */
	MAC_EXP_OP_TOK_NE   = 2 /* != */
	MAC_EXP_OP_TOK_LT   = 3 /* < */
	MAC_EXP_OP_TOK_LE   = 4 /* <= */
	MAC_EXP_OP_TOK_GE   = 5 /* >= */
	MAC_EXP_OP_TOK_GT   = 6 /* > */
)

/*
 * Relational operators. The MAC_EXP_OP_TOK_* are defined in the header
 * file.
 */
const (
	MAC_EXP_OP_STR_EQ  = "=="
	MAC_EXP_OP_STR_NE  = "!="
	MAC_EXP_OP_STR_LT  = "<"
	MAC_EXP_OP_STR_LE  = "<="
	MAC_EXP_OP_STR_GE  = ">="
	MAC_EXP_OP_STR_GT  = ">"
	MAC_EXP_OP_STR_ANY = "\"" + MAC_EXP_OP_STR_EQ +
		"\" or \"" + MAC_EXP_OP_STR_NE + "\"" +
		"\" or \"" + MAC_EXP_OP_STR_LT + "\"" +
		"\" or \"" + MAC_EXP_OP_STR_LE + "\"" +
		"\" or \"" + MAC_EXP_OP_STR_GE + "\"" +
		"\" or \"" + MAC_EXP_OP_STR_GT + "\""
)

var mac_exp_op_table = map[string]int{
	MAC_EXP_OP_STR_EQ: MAC_EXP_OP_TOK_EQ,
	MAC_EXP_OP_STR_NE: MAC_EXP_OP_TOK_NE,
	MAC_EXP_OP_STR_LT: MAC_EXP_OP_TOK_LT,
	MAC_EXP_OP_STR_LE: MAC_EXP_OP_TOK_LE,
	MAC_EXP_OP_STR_GE: MAC_EXP_OP_TOK_GE,
	MAC_EXP_OP_STR_GT: MAC_EXP_OP_TOK_GT,
}

var mac_exp_ext_table = map[string]func(string, int, string) int{}

var mac_exp_op_table_str = map[int]string{
	MAC_EXP_OP_TOK_EQ: MAC_EXP_OP_STR_EQ,
	MAC_EXP_OP_TOK_NE: MAC_EXP_OP_STR_NE,
	MAC_EXP_OP_TOK_LT: MAC_EXP_OP_STR_LT,
	MAC_EXP_OP_TOK_LE: MAC_EXP_OP_STR_LE,
	MAC_EXP_OP_TOK_GE: MAC_EXP_OP_STR_GE,
	MAC_EXP_OP_TOK_GT: MAC_EXP_OP_STR_GT,
}

const MAC_EXP_WHITESPACE string = CHARS_SPACE

type MacExpandContext struct {
	Result  *strings.Builder
	Flags   int
	Filter  *string
	Lookup  func(string, int, interface{}) (string, error)
	Context interface{}
	Status  int
	Level   int
}

func (m *MacExpandContext) String() string {
	return fmt.Sprintf("{result:'%s', flags:%d, filter:'%v', lookup:%p, context:%v, status:%d, level:%d}",
		m.Result, m.Flags, m.Filter, m.Lookup, m.Context, m.Status, m.Level)
}

func IsHSpace(r rune) bool {
	return r == ' ' || r == '\t' || r == '\r' || r == '\n'
}

func atol_or_die(strval string) int {
	if n, err := strconv.Atoi(strval); err == nil {
		return n
	} else {
		MsgFatal("mac_exp_eval: bad conversion", "strval", strval)
	}
	return 0
}

func mac_exp_parse_error(mc *MacExpandContext, format string, a ...any) int {
	MsgWarnf(format, a...)
	mc.Status |= MAC_PARSE_ERROR
	return mc.Status
}

/* mac_exp_eval - evaluate binary expression */
func mac_exp_eval(left string, tok_val int, rite string) int {
	var delta int
	// numerical or string comparison
	if AllDig(left) && AllDig(rite) {
		delta = atol_or_die(left) - atol_or_die(rite)
	} else {
		delta = strings.Compare(left, rite)
	}
	switch tok_val {
	case MAC_EXP_OP_TOK_EQ:
		return mac_exp_op_res_bool[delta == 0]
	case MAC_EXP_OP_TOK_NE:
		return mac_exp_op_res_bool[delta != 0]
	case MAC_EXP_OP_TOK_LT:
		return mac_exp_op_res_bool[delta < 0]
	case MAC_EXP_OP_TOK_LE:
		return mac_exp_op_res_bool[delta <= 0]
	case MAC_EXP_OP_TOK_GE:
		return mac_exp_op_res_bool[delta >= 0]
	case MAC_EXP_OP_TOK_GT:
		return mac_exp_op_res_bool[delta > 0]
	default:
		MsgPanic("unknown operator", "function", "mac_exp_eval", "tok_value", tok_val)
	}
	return MAC_EXP_OP_RES_ERROR
}

func mac_exp_extract_curly_payload(mc *MacExpandContext, bp *RuneScanner) *string {
	/*
	 * Extract the payload and balance the {}. The caller is expected to skip
	 * leading whitespace before the {. See MAC_EXP_FIND_LEFT_CURLY().
	 */
	level := 1

	bp.Next() // skip the first opening brace
	for ch := bp.Scan(); ; ch = bp.Scan() {
		if ch == scanner.EOF {
			mac_exp_parse_error(mc, "unbalanced {} in attribute expression: %s\n", bp.String())
			return nil
		} else if ch == '{' {
			level++
		} else if ch == '}' {
			if level--; level <= 0 {
				break
			}
		}
	}

	bp.Push() // temporarily push the closing '}'
	payload := bp.Emit()
	bp.Next() // consume the closing '}'

	// Skip trailing whitespace after }.
	bp.Skip(IsHSpace)

	return &payload
}

func mac_exp_parse_relational(mc *MacExpandContext, lookup *string, bp *RuneScanner) int {
	var left_op_strval *string
	var rite_op_strval *string
	var relop_eval func(string, int, string) int
	var op_result int

	/*
	 * Left operand. The caller is expected to skip leading whitespace before
	 * the {. See scan.Span(IsHSpace).
	 */
	if left_op_strval = mac_exp_extract_curly_payload(mc, bp); left_op_strval == nil {
		return mc.Status
	}

	/*
	 * Operator. TODO: regexp operator.
	 */
	op_pos := bp.Off
	bp.Span(func(r rune) bool { return strings.ContainsRune("<>!=?+-*/~&|%", r) })
	op_strval := bp.Emit()
	op_tokval := mac_exp_op_table[op_strval]
	if op_tokval == MAC_EXP_OP_TOK_NONE {
		return mac_exp_parse_error(mc, "%s expected at: \"...%s}>>>%s%.20s\"", MAC_EXP_OP_STR_ANY, *left_op_strval, op_strval, bp.PosString())
	}

	/*
	 * Custom operator suffix.
	 */
	if len(mac_exp_ext_table) != 0 && ascii.IsAlnum(bp.Peek()) {
		var exists bool
		bp.Span(ascii.IsAlnum)
		op_ext := bp.Emit()
		mac_exp_ext_key := fmt.Sprintf("%s%s", op_strval, op_ext)

		if relop_eval, exists = mac_exp_ext_table[mac_exp_ext_key]; !exists {
			return mac_exp_parse_error(mc, "bad operator suffix at: \"...%s>>>%s\"", op_strval, op_ext)
		}
	} else {
		relop_eval = mac_exp_eval
	}

	/*
	 * Right operand. Todo: syntax may depend on operator.
	 */
	if bp.Skip(IsHSpace); bp.Peek() != '{' {
		return mac_exp_parse_error(mc, "\"{expression}\" expected at: \"...{%s} %s>>>%.20s\"", *left_op_strval, bp.String()[op_pos:bp.Pos], bp.PosString())
	}
	if rite_op_strval = mac_exp_extract_curly_payload(mc, bp); rite_op_strval == nil {
		return mc.Status
	}

	/*
	 * Evaluate the relational expression. Todo: regexp support.
	 */
	left_op_buf := strings.Builder{}
	mc.Status |= MacExpand(&left_op_buf, *left_op_strval, mc.Flags, mc.Filter, mc.Lookup, mc.Context)
	rite_op_buf := strings.Builder{}
	mc.Status |= MacExpand(&rite_op_buf, *rite_op_strval, mc.Flags, mc.Filter, mc.Lookup, mc.Context)
	if mc.Flags&MAC_EXP_FLAG_SCAN == 0 {
		op_result = relop_eval(left_op_buf.String(), op_tokval, rite_op_buf.String())
		if op_result == MAC_EXP_OP_RES_ERROR {
			mc.Status |= MAC_PARSE_ERROR
		}
	}
	if mc.Status&MAC_PARSE_ERROR == MAC_PARSE_ERROR {
		return mc.Status
	}

	/*
	 * Here, we fake up a non-empty or empty parameter value lookup result,
	 * for compatibility with the historical code that looks named parameter
	 * values.
	 */
	if mc.Flags&MAC_EXP_FLAG_SCAN == MAC_EXP_FLAG_SCAN {
		lookup = nil
	} else {
		switch op_result {
		case MAC_EXP_OP_RES_TRUE:
			*lookup = MAC_EXP_BVAL_TRUE
		case MAC_EXP_OP_RES_FALSE:
			*lookup = MAC_EXP_BVAL_FALSE
		default:
			MsgPanic("mac_expand: unexpected operator result", "op_result", op_result)
		}
	}

	return 0
}

/* mac_expand_add_relop - register operator extensions */

func mac_expand_add_relop(tok_list []int, suffix string, relop_eval func(string, int, string) int) {
	var myname string = "mac_expand_add_relop"
	var tok_name string
	// int    *tp;

	/*
	 * Sanity checks.
	 */
	if !AllAlnum(suffix) {
		MsgPanic("bad operator suffix", "function", myname, "suffix", suffix)
	}

	for _, tok := range tok_list {
		var exists bool
		if tok_name, exists = mac_exp_op_table_str[tok]; !exists {
			MsgPanic("unknown token code", "function", myname, "token", tok)
		}
		mac_exp_ext_key := fmt.Sprintf("%s%s", tok_name, suffix)

		if _, exists = mac_exp_ext_table[mac_exp_ext_key]; exists {
			MsgPanic("duplicate key", "function", myname, "key", mac_exp_ext_key)
		}
		mac_exp_ext_table[mac_exp_ext_key] = relop_eval
	}
}

func mac_expand_callback(typ int, buf string, context interface{}) int {
	var myname = "mac_expand_callback"
	var lookup_mode int
	var lookup string
	var lookup_err error
	var res_iftrue *string
	var res_iffalse *string

	mc := context.(*MacExpandContext)

	if mc.Level++; mc.Level > 100 {
		mac_exp_parse_error(mc, "unreasonable macro call nesting: \"%s\"", buf)
	}
	if mc.Status&MAC_PARSE_ERROR == MAC_PARSE_ERROR {
		return mc.Status
	}

	scan := NewRuneScanner(buf)

	if typ == MAC_PARSE_EXPR { // Relational expression.
		// If recursion is disabled, perform only one
		// level of $name expansion.
		if scan.Skip(IsHSpace); scan.Peek() == '{' {
			if mac_exp_parse_relational(mc, &lookup, scan) != 0 {
				return mc.Status
			}

			// Look for the ? or : operator.
			if ch := scan.Peek(); ch != scanner.EOF {
				if ch != '?' && ch != ':' {
					return mac_exp_parse_error(mc, "\"?\" or \":\" expected at: \"...}>>>%.20s\"", scan.PosString())
				}
			}
		} else { // named parameter
			scan.Sync()
			scan.Skip(IsHSpace)
			for ch := scan.Scan(); ch != scanner.EOF; ch = scan.Scan() {
				if scan.Peek() == scanner.EOF {
					lookup_mode = MAC_EXP_MODE_USE
					break
				}
				if ch == '?' || ch == ':' {
					lookup_mode = MAC_EXP_MODE_TEST
					scan.Push()
					break
				}
				if !ascii.IsAlnum(ch) && ch != '_' {
					// move the position back one rune so our debug log is pointing
					// at the character that caused the parse failure
					scan.Push()
					return mac_exp_parse_error(mc, "attribute name syntax error at: \"...%s>>>%.20s\"",
						scan.String()[:scan.Pos], scan.PosString())
				}
			}

			lookup, lookup_err = mc.Lookup(strings.TrimSpace(scan.Emit()), lookup_mode, mc.Context)
		}

		/*
		 * Return the requested result. After parsing the result operand
		 * following ?, we fall through to parse the result operand following
		 * :. This is necessary with the ternary ?: operator: first, with
		 * MAC_EXP_FLAG_SCAN to parse both result operands with mac_parse(),
		 * and second, to find garbage after any result operand. Without
		 * MAC_EXP_FLAG_SCAN the content of only one of the ?: result
		 * operands will be parsed with mac_parse(); syntax errors in the
		 * other operand will be missed.
		 */
		switch ch := scan.Next(); ch {
		case '?':
			if scan.Skip(IsHSpace); scan.Peek() == '{' {
				if res_iftrue = mac_exp_extract_curly_payload(mc, scan); res_iftrue == nil {
					return mc.Status
				}
			} else {
				str := scan.PosString()
				res_iftrue = &str
				scan.End(true)
			}
			if (lookup_err == nil && lookup != "") || (mc.Flags&MAC_EXP_FLAG_SCAN == MAC_EXP_FLAG_SCAN) {
				mc.Status |= MacParse(*res_iftrue, mac_expand_callback, mc)
			}
			if scan.Peek() == scanner.EOF {
				break
			} else if scan.Peek() != ':' {
				return mac_exp_parse_error(mc, "\":\" expected at: \"...%s}>>>%.20s\"", *res_iftrue, scan.PosString())
			}
			scan.Next()
			fallthrough
		case ':':
			if scan.Skip(IsHSpace); scan.Peek() == '{' {
				if res_iffalse = mac_exp_extract_curly_payload(mc, scan); res_iffalse == nil {
					return mc.Status
				}
			} else {
				str := scan.PosString()
				res_iffalse = &str
				scan.End(true)
			}
			if lookup_err != nil || lookup == "" || (mc.Flags&MAC_EXP_FLAG_SCAN == MAC_EXP_FLAG_SCAN) {
				mc.Status |= MacParse(*res_iffalse, mac_expand_callback, mc)
			}
			if scan.Peek() != scanner.EOF {
				return mac_exp_parse_error(mc, "unexpected input at: \"...%s}>>>%.20s\"", *res_iffalse, scan.PosString())
			}
		case scanner.EOF:
			if lookup_err != nil {
				mc.Status |= MAC_PARSE_UNDEF
			} else if lookup == "" || mc.Flags&MAC_EXP_FLAG_SCAN == MAC_EXP_FLAG_SCAN {
				// empty
			} else if mc.Flags&MAC_EXP_FLAG_RECURSE == MAC_EXP_FLAG_RECURSE {
				mc.Status |= MacParse(lookup, mac_expand_callback, mc)
			} else {
				if mc.Flags&MAC_EXP_FLAG_PRINTABLE == MAC_EXP_FLAG_PRINTABLE {
					lookup = strings.Map(func(r rune) rune {
						if unicode.IsPrint(r) {
							return r
						}
						return '_'
					}, lookup)
				} else if mc.Filter != nil {
					lookup = strings.Map(func(r rune) rune {
						if strings.ContainsRune(*mc.Filter, r) {
							return r
						}
						return '_'
					}, lookup)
				}
				mc.Result.WriteString(lookup)
			}
		default:
			MsgPanic("unknown operator code", "function", myname, "code", ch)
		}

	} else if mc.Flags&MAC_EXP_FLAG_SCAN == 0 { // Literal text.
		mc.Result.WriteString(buf)
	}
	mc.Level--

	return mc.Status
}

func MacExpand(result *strings.Builder, pattern string, flags int, filter *string, lookup func(string, int, interface{}) (string, error), context interface{}) int {
	mc := MacExpandContext{
		Result:  result,
		Flags:   flags,
		Filter:  filter,
		Lookup:  lookup,
		Context: context,
		Status:  0,
		Level:   0,
	}

	if flags&(MAC_EXP_FLAG_APPEND|MAC_EXP_FLAG_SCAN) == 0 {
		result.Reset()
	}
	status := MacParse(pattern, mac_expand_callback, &mc)
	return status
}
