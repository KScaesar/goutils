package testingY

import (
	"bytes"
	stdJSON "encoding/json"
	"fmt"
	"regexp"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/json"
)

func RawJSON(prettyJSON string) string {
	m := minify.New()
	m.AddFuncRegexp(regexp.MustCompile("[/+]json$"), json.Minify)
	rawJSON, err := m.String("application/json", prettyJSON)
	if err != nil {
		panic(fmt.Sprintf("transform RawJSON: %v", err)) // 避免循環依賴
	}
	return rawJSON
}

func PrettyJSON(rawJSON []byte) string {
	var prettyJSON bytes.Buffer
	if err := stdJSON.Indent(&prettyJSON, rawJSON, "", "  "); err != nil {
		panic(fmt.Sprintf("transform To PrettyJSON: %v", err))
	}
	return prettyJSON.String()
}

func RawSQL(prettySQL string) (rawSQL string) {
	rawSQL = prettySQL
	for _, s := range []string{"\t", "\n"} {
		rawSQL = strings.ReplaceAll(rawSQL, s, "")
	}
	return rawSQL
}
