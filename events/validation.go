package events

import "regexp"

var titleRe = regexp.MustCompile(`^[а-яА-Яa-zA-Z0-9 ,.]{3,50}$`)

func isValidTitle(title string) bool {
	return titleRe.MatchString(title)
}
