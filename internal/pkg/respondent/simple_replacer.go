package respondent

import "errors"

type SimpleReplacer struct {
	rules []replacementRule
}

type replacementRule struct {
	Original    error
	Replacement error
}

func (sr *SimpleReplacer) Replace(err error) error {
	if err == nil {
		return nil
	}
	for _, rule := range sr.rules {
		if errors.Is(err, rule.Original) {
			return rule.Replacement
		}
	}
	return err
}

func (sr *SimpleReplacer) ReplaceBy(original, replacement error) *SimpleReplacer {
	if original == nil || replacement == nil {
		return nil
	}
	sr.rules = append(sr.rules, replacementRule{
		Original:    original,
		Replacement: replacement,
	})
	return sr
}

func NewSimpleReplacer() *SimpleReplacer {
	return &SimpleReplacer{
		rules: make([]replacementRule, 0),
	}
}
