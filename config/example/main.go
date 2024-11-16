package main

import (
	"fmt"
	"github.com/syntaxfa/syntax-backend/config"
	"strings"
)

func main() {
	//if err := os.Setenv("SYNTAX_LOGGER__USE_LOCAL_TIME", "true"); err != nil {
	//	panic(err)
	//}
	//
	//fmt.Printf("%+v\n", exampletwo.LoadConfig())

	cfg := config.New(config.Options{
		Prefix:       "SYNTAX_",
		Delimiter:    ".",
		Separator:    "__",
		YamlFilePath: "config.yml",
		CallBackEnV: func(s string) string {
			base := strings.ToLower(strings.TrimPrefix(s, "SYNTAX_"))
			return strings.ReplaceAll(base, "__", ".")
		},
	})

	fmt.Printf("cfg: %+v\n", cfg)
}
