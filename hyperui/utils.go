package hyperui

import (
	twmerge "github.com/Oudwins/tailwind-merge-go"
	h "github.com/assaidy/hyper/v2"
)

func mergeStyles(element *h.Element, styles string) {
	for i := len(element.Attributes) - 1; i >= 0; i-- {
		if a, ok := element.Attributes[i].(h.PairAttribute); ok {
			if a.Key == "class" {
				a.Value = twmerge.Merge(styles, a.Value)
				element.Attributes[i] = a
				return
			}
		}
	}

	element.Attributes = append(element.Attributes, h.AttrClass(styles))
}
