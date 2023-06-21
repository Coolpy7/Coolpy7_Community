package topic

import (
	"errors"
	"strings"
)

var ErrZeroLength = errors.New("zero length topic")
var ErrOutLength = errors.New("out length topic")
var ErrWildcards = errors.New("invalid use of wildcards")

func trimSlash(r rune) bool {
	return r == '/'
}

func Parse(topic string, allowWildcards bool) (string, error) {
	if topic == "" {
		return "", ErrZeroLength
	}

	if len(topic) > 65535 {
		return "", ErrOutLength
	}

	if hasAdjacentSlashes(topic) {
		topic = collapseSlashes(topic)
	}

	topic = strings.TrimRightFunc(topic, trimSlash)
	if topic == "" {
		return "", ErrZeroLength
	}

	remainder := topic
	segment := topicSegment(topic, "/")

	for segment != topicEnd {
		if (strings.Contains(segment, "+") || strings.Contains(segment, "#")) && len(segment) > 1 {
			return "", ErrWildcards
		}

		if !allowWildcards && (segment == "#" || segment == "+") {
			return "", ErrWildcards
		}

		if segment == "#" && topicShorten(remainder, "/") != topicEnd {
			return "", ErrWildcards
		}

		remainder = topicShorten(remainder, "/")
		segment = topicSegment(remainder, "/")
	}

	return topic, nil
}

func ContainsWildcards(topic string) bool {
	return strings.Contains(topic, "+") || strings.Contains(topic, "#")
}

func hasAdjacentSlashes(str string) bool {
	var last rune
	for _, r := range str {
		if r == '/' && last == '/' {
			return true
		}
		last = r
	}

	return false
}

func collapseSlashes(str string) string {
	var b strings.Builder
	var last rune
	for _, r := range str {
		if r == '/' && last == '/' {
			continue
		}
		b.WriteRune(r)
		last = r
	}

	return b.String()
}
