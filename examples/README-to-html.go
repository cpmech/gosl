// Copyright 2016 The Gosl Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// build +ignore

package main

import (
	"bytes"
	"os"
	"strings"

	"github.com/cpmech/gosl/io"
	"github.com/russross/blackfriday"
)

func main() {

	// read README.md file
	md, err := io.ReadFile("README.md")
	if err != nil {
		io.PfRed("cannot read README.md\n")
		return
	}

	flags := 0 |
		blackfriday.HTML_USE_XHTML |
		blackfriday.HTML_USE_SMARTYPANTS |
		blackfriday.HTML_SMARTYPANTS_LATEX_DASHES

	extensions := 0 |
		blackfriday.EXTENSION_NO_INTRA_EMPHASIS |
		blackfriday.EXTENSION_TABLES |
		blackfriday.EXTENSION_FENCED_CODE |
		blackfriday.EXTENSION_AUTOLINK |
		blackfriday.EXTENSION_STRIKETHROUGH |
		blackfriday.EXTENSION_SPACE_HEADERS |
		blackfriday.EXTENSION_HEADER_IDS |
		blackfriday.EXTENSION_BACKSLASH_LINE_BREAK |
		blackfriday.EXTENSION_DEFINITION_LISTS

	renderer := blackfriday.HtmlRenderer(flags, "", "")
	html := string(blackfriday.MarkdownOptions(md, renderer, blackfriday.Options{Extensions: extensions}))

	// set path to source files
	path := "https://github.com/cpmech/gosl/blob/master/examples/"
	html = strings.Replace(html, "a href=\"", io.Sf("a href=\"%s", path), -1)

	// set path of figures
	path = os.ExpandEnv("${HOME}/10.go/src/github.com/cpmech/gosl/examples/")
	html = strings.Replace(html, "img src=\"", io.Sf("img src=\"%s", path), -1)

	// set header and footer
	html = `<!DOCTYPE HTML>
<html>
<head>
<title>Gosl Examples</title>
<meta charset="utf-8" />

<style>
h1 {color:#0064cb; font-family:verdana; font-size:200%;}
h2 {color:#0064cb}
h3 {color:#0064cb}
a:hover {background-color:#5397dc;}
#container {
	width:500px;
	text-align:center;
}
#container img {
	max-width:100%;
	height:auto;
}
</style>

</head>
<body>
` + html + `
</body>
</html>`

	// write file
	io.WriteFileVD("/tmp", "gosl-README.html", bytes.NewBuffer([]byte(html)))
}
