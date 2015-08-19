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
