package jwtx

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"regexp"
)

const regFmt = `^(Bearer ){1}((%s)+)(\.){1}([-_a-zA-Z0-9]+)(\.){1}([-_a-zA-Z0-9]+)$`

var (
	parse = func(tkStr string, header Header) (string, bool) {
		prefix, pattern, err := parserFmt(header)
		if err != nil {
			return "", false
		}

		return string(prefix.ReplaceAll([]byte(tkStr), []byte(""))), pattern.Match([]byte(tkStr))
	}

	parserFmt = func(header Header) (prefix, pattern *regexp.Regexp, err error) {
		parserH, err := encodeHeader(header)
		if err != nil {
			return nil, nil, err
		}

		prefix = regexp.MustCompile(`^Bearer `)
		fmt.Println(fmt.Sprintf(regFmt, parserH))
		pattern = regexp.MustCompile(fmt.Sprintf(regFmt, parserH))

		return
	}

	encodeHeader = func(header Header) (string, error) {
		bytes, err := json.Marshal(header)
		if err != nil {
			return "", err
		}

		fmt.Println(base64.StdEncoding.EncodeToString(bytes))
		return base64.StdEncoding.EncodeToString(bytes), nil
	}
)

func (p *JWTes256) Parse(token string) (string, bool) {
	return parse(token, p.Header)
}

func (p *JWThs256) Parse(token string) (string, bool) {
	return parse(token, p.Header)
}

func (p *JWTrs256) Parse(token string) (string, bool) {
	return parse(token, p.Header)
}
