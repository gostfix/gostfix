package global

import "github.com/gostfix/gostfix/internal/util"

func MailParamSplit(name string, value string) []string {
	tokens := util.MyStrTokQ(value, util.CHARS_COMMA_SP, util.CHARS_BRACE)
	for n, tok := range tokens {
		if tok[0] == util.CHARS_BRACE[0] {
			var etok string
			var err error
			if etok, err = util.ExtPar(tok, util.CHARS_BRACE, util.EXTPAR_FLAG_STRIP); err != nil {
				util.MsgWarn(err.Error())
			}
			tokens[n] = etok
		}
	}
	return tokens
}
