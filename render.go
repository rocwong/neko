package neko

type (
	// HtmlEngine is an interface for parsing html templates and redering HTML.
	HtmlEngine interface {
		Render(view string, context interface{}, status ...int) error
	}
)
