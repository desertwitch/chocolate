package chocolate

import (
	"slices"

	"github.com/charmbracelet/lipgloss"
)

type chocolateBarConstrainer struct {
	widthConstraints  []barSizeConstraint
	heightConstraints []barSizeConstraint

	nonTargetConstraints []ConstraintAttribute

	bias bool
}

func (cb *chocolateBarConstrainer) canBias() bool { return cb.bias }

func (cb *chocolateBarConstrainer) sizeConstraints() (width, height []barSizeConstraint) {
	return cb.widthConstraints, cb.heightConstraints
}

func (cb *chocolateBarConstrainer) constraintTarget(a ConstraintAttribute) bool {
	return !slices.Contains(cb.nonTargetConstraints, a)
}

type hiddenBarConstrainer struct {
	chocolateBarConstrainer
}

func newHiddenConstrainer() *hiddenBarConstrainer {
	return &hiddenBarConstrainer{
		chocolateBarConstrainer: chocolateBarConstrainer{
			widthConstraints: []barSizeConstraint{
				{
					Relation: EQ,
					Value:    0,
				},
			},
			heightConstraints: []barSizeConstraint{
				{
					Relation: EQ,
					Value:    0,
				},
			},
			nonTargetConstraints: []ConstraintAttribute{
				WIDTH,
				HEIGHT,
			},
			bias: false,
		},
	}
}

type guideBarConstrainer struct {
	chocolateBarConstrainer
}

func newGuideConstrainer() *guideBarConstrainer {
	return &guideBarConstrainer{
		chocolateBarConstrainer: chocolateBarConstrainer{
			widthConstraints: []barSizeConstraint{
				{
					Relation: GE,
					Value:    -0,
				},
			},
			heightConstraints: []barSizeConstraint{
				{
					Relation: GE,
					Value:    -0,
				},
			},
			bias: true,
		},
	}
}

type styledBarConstrainer struct {
	chocolateBarConstrainer
	style   *lipgloss.Style
	content *string
}

func (sbc *styledBarConstrainer) sizeConstraints() (width, height []barSizeConstraint) {
	wconstant := 1.0
	hconstant := 1.0
	if sbc.content != nil {
		w, h := lipgloss.Size(*sbc.content)
		wconstant = float64(w)
		hconstant = float64(h)
	}
	if sbc.style != nil {
		wconstant += float64(sbc.style.GetHorizontalFrameSize())
		hconstant += float64(sbc.style.GetVerticalFrameSize())
	}

	sbc.widthConstraints = []barSizeConstraint{
		{
			Relation: GE,
			Value:    -wconstant,
		},
	}
	sbc.heightConstraints = []barSizeConstraint{
		{
			Relation: GE,
			Value:    -hconstant,
		},
	}

	return sbc.chocolateBarConstrainer.sizeConstraints()
}

func newStyledConstrainer(style *lipgloss.Style, content ...*string) *styledBarConstrainer {
	var c *string
	if len(content) >= 1 {
		c = content[0]
	}

	return &styledBarConstrainer{
		chocolateBarConstrainer: chocolateBarConstrainer{
			bias: true,
		},
		style:   style,
		content: c,
	}
}
