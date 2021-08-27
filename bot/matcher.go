package bot

import "regexp"

type Matcher string

func (m Matcher) getArguments(text string) (args map[string]string) {
	var re = regexp.MustCompile(m.toPerlSyntax())
	match := re.FindStringSubmatch(text)

	args = make(map[string]string)
	for i, name := range re.SubexpNames() {
		if i > 0 && i <= len(match) {
			args[name] = match[i]
		}
	}
	return args
}

// Matcher syntax to perl syntax
func (m Matcher) toPerlSyntax() (result string) {
	re := regexp.MustCompile("<")
	result = re.ReplaceAllString(string(m), "(?P<")

	re = regexp.MustCompile(">")
	result = re.ReplaceAllString(result, ">.*)")

	return result
}

// Matcher syntax to regex
func (m Matcher) toRegex() (result string) {
	re := regexp.MustCompile("<([^>]*)>")
	return re.ReplaceAllString(string(m), "(.*)")
}
