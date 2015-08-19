// The MIT License (MIT)
//
// Copyright (c) 2015 Nick Powell
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package main

// Context is what's given to every command handler.  It should contain
// everything a command will need
type Context struct {
	Bot *slack.Client
	RTM *slack.RTM
	Ev  *slack.MessageEvent
}

func (c *Context) Respond(s string) {
	Respond(s, c.RTM, c.Ev)
}

func parseTemplate(filename string) *template.Template {
	// Default directory if we're in a Docker environment
	lookupName := filepath.Join("/assets", filename)

	// Sometimes it might just be in the current working directory
	if _, err := os.Stat(lookupName); err != nil {
		lookupName = filename
	}

	return template.Must(template.ParseFiles(lookupName))
}

func (ctx *Context) RenderTemplate(tmpl *template.Template, g map[string]interface{}) {
	buf := bytes.NewBufferString("")
	tmpl.Execute(buf, g)
	ctx.Respond(buf.String())
}

// helper context functions
// these can be thought of as middleware or replacements.
// or small functions that take a context as an argument

func notReadyYet(ctx *Context) {
	ctx.Respond("That command isn't ready yet")
}
