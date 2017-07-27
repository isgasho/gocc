//Copyright 2012 Vastech SA (PTY) LTD
//
//   Licensed under the Apache License, Version 2.0 (the "License");
//   you may not use this file except in compliance with the License.
//   You may obtain a copy of the License at
//
//       http://www.apache.org/licenses/LICENSE-2.0
//
//   Unless required by applicable law or agreed to in writing, software
//   distributed under the License is distributed on an "AS IS" BASIS,
//   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
//   See the License for the specific language governing permissions and
//   limitations under the License.
package golang

import (
	"bytes"
	"path"
	"text/template"

	"github.com/goccmack/gocc/internal/io"
)

func GenErrors(pkg, outDir string) {
	tmpl, err := template.New("parser errors").Parse(errorsSrc[1:])
	if err != nil {
		panic(err)
	}
	wr := new(bytes.Buffer)
	if err := tmpl.Execute(wr, path.Join(pkg, "token")); err != nil {
		panic(err)
	}
	io.WriteFile(path.Join(outDir, "errors", "errors.go"), wr.Bytes())
}

const errorsSrc = `
// Code generated by gocc; DO NOT EDIT.

package errors

import (
	"bytes"
	"fmt"

	"{{.}}"
)

type ErrorSymbol interface {
}

type Error struct {
	Err            error
	ErrorToken     *token.Token
	ErrorSymbols   []ErrorSymbol
	ExpectedTokens []string
	StackTop       int
}

func (e *Error) String() string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "Error")
	if e.Err != nil {
		fmt.Fprintf(w, " %s\n", e.Err)
	} else {
		fmt.Fprintf(w, "\n")
	}
	fmt.Fprintf(w, "Token: type=%d, lit=%s\n", e.ErrorToken.Type, e.ErrorToken.Lit)
	fmt.Fprintf(w, "Pos: offset=%d, line=%d, column=%d\n", e.ErrorToken.Pos.Offset, e.ErrorToken.Pos.Line, e.ErrorToken.Pos.Column)
	fmt.Fprintf(w, "Expected one of: ")
	for _, sym := range e.ExpectedTokens {
		fmt.Fprintf(w, "%s ", sym)
	}
	fmt.Fprintf(w, "ErrorSymbol:\n")
	for _, sym := range e.ErrorSymbols {
		fmt.Fprintf(w, "%v\n", sym)
	}
	return w.String()
}

func (e *Error) Error() string {
	w := new(bytes.Buffer)
	fmt.Fprintf(w, "Error in S%d: %s, %s", e.StackTop, token.TokMap.TokenString(e.ErrorToken), e.ErrorToken.Pos.String())
	if e.Err != nil {
		fmt.Fprintf(w, e.Err.Error())
	} else {
		fmt.Fprintf(w, ", expected one of: ")
		for _, expected := range e.ExpectedTokens {
			fmt.Fprintf(w, "%s ", expected)
		}
	}
	return w.String()
}
`
