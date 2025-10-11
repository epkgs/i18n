package internal

import "golang.org/x/text/language"

type Matcher struct {
	strict  bool
	langs   []language.Tag
	matcher language.Matcher
}

func NewMatcher(defaultLanguage language.Tag, limits ...language.Tag) *Matcher {
	langs := append([]language.Tag{defaultLanguage}, limits...)
	return &Matcher{
		strict:  len(limits) > 0,
		langs:   langs,
		matcher: language.NewMatcher(langs),
	}
}

// Match returns the best match for any of the given tags, along with
// a unique index associated with the returned tag and a confidence
// score.
func (m *Matcher) Match(t ...language.Tag) language.Tag {
	_, i, conf := m.matcher.Match(t...)
	if conf <= language.Low {
		return m.langs[0]
	}

	return m.langs[i]
}

// MatchOrAdd acts like Match but it checks and adds a language tag, if not found,
// when the `Matcher.strict` field is true (when no tags are provided by the caller)
// and they should be dynamically added to the list.
func (m *Matcher) MatchOrAdd(t language.Tag) language.Tag {
	_, i, conf := m.matcher.Match(t)
	if conf <= language.Low {
		if !m.strict {
			// not found, add it now.
			m.langs = append(m.langs, t)
			m.matcher = language.NewMatcher(m.langs) // reset matcher to include the new language.
		}
		return t
	}

	return m.langs[i]
}

func (m *Matcher) DefaultLanguage() language.Tag {
	return m.langs[0]
}

func (m *Matcher) Languages() []language.Tag {
	return m.langs
}

func (m *Matcher) SetLanguages(langs []language.Tag) {
	m.langs = langs
	m.matcher = language.NewMatcher(langs)
}
