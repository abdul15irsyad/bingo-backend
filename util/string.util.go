package util

import (
	"regexp"
	"strings"
)

func Slugify(text string) string {
	slug := strings.ToLower(text)

	re := regexp.MustCompile(`[^\w\s-]`)
	slug = re.ReplaceAllString(slug, "")

	re = regexp.MustCompile(`[\s-]+`)
	slug = re.ReplaceAllString(slug, "-")

	slug = strings.Trim(slug, "-")

	return slug
}
