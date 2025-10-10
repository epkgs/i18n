package i18n

import "golang.org/x/text/language"

type Matcher struct {
	strict    bool
	Languages []language.Tag
	matcher   language.Matcher
}

var _ language.Matcher = (*Matcher)(nil)

func newMatcher(defaultLanguage language.Tag, limits ...language.Tag) *Matcher {
	langs := append([]language.Tag{defaultLanguage}, limits...)
	return &Matcher{
		strict:    len(limits) > 0,
		Languages: langs,
		matcher:   language.NewMatcher(langs),
	}
}

// Match returns the best match for any of the given tags, along with
// a unique index associated with the returned tag and a confidence
// score.
func (m *Matcher) Match(t ...language.Tag) (language.Tag, int, language.Confidence) {
	return m.matcher.Match(t...)
}

// MatchOrAdd acts like Match but it checks and adds a language tag, if not found,
// when the `Matcher.strict` field is true (when no tags are provided by the caller)
// and they should be dynamically added to the list.
func (m *Matcher) MatchOrAdd(t language.Tag) (tag language.Tag, index int, conf language.Confidence) {
	tag, index, conf = m.Match(t)
	if conf <= language.Low && !m.strict {
		// not found, add it now.
		m.Languages = append(m.Languages, t)
		tag = t
		index = len(m.Languages) - 1
		conf = language.Exact
		m.matcher = language.NewMatcher(m.Languages) // reset matcher to include the new language.
	}

	return
}
