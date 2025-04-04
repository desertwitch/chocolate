package chocolate

import (
	"encoding/json"
	"fmt"
	"slices"
	"sort"

	"github.com/charmbracelet/lipgloss"
	"github.com/lithdew/casso"
)

type constraintElement struct {
	width  casso.Symbol
	height casso.Symbol
	xpos   casso.Symbol
	ypos   casso.Symbol
}

type constraintLayout struct {
	width       int
	height      int
	children    map[string]barChild
	constraints []Constraint
	failsMax    int
	dirty       bool
}

func (c *constraintLayout) addBar(n string, v barChild) {
	if c.children == nil {
		c.children = make(map[string]barChild)
	}
	v.setParent(c)
	c.children[n] = v

	c.dirty = true
}

func (c *constraintLayout) Resize(width, height int) {
	// s := lipgloss.NewStyle().Border(lipgloss.NormalBorder())
	if c.width != width || c.height != height {
		c.setDirty()
		// c.width = width - s.GetHorizontalFrameSize()
		// c.height = height - s.GetVerticalFrameSize()
	}
	c.width = width
	c.height = height
}

func (c *constraintLayout) View() string {
	// s := lipgloss.NewStyle().Border(lipgloss.NormalBorder())

	bars, err := c.resolve(0)
	if err != nil {
		return ""
	}

	// ret := s.Width(c.width).
	// 	Height(c.height).
	// 	BorderBottom(false).
	// 	BorderTop(false).
	// 	BorderLeft(false).
	// 	BorderRight(false).
	// 	Render("")

	ret := lipgloss.Place(c.width, c.height, 0, 0, "")

	for _, b := range bars {
		xPos := b.xpos()
		yPos := b.ypos()
		ret = placeOverlay(xPos, yPos, b.View(), ret)
		if c.dirty {
			return c.View()
		}
	}

	// ret = s.Width(c.width).
	// 	Height(c.height).
	// 	Render(ret)
	return ret
}

func (c *constraintLayout) setDirty() { c.dirty = true }

func (c *constraintLayout) AddConstraint(constraint Constraint) {
	c.constraints = append(c.constraints, constraint)
	c.dirty = true
}

func (c *constraintLayout) resolve(f int) (map[string]barChild, error) {
	if !c.dirty {
		return c.children, nil
	}

	if f > c.failsMax {
		return nil, fmt.Errorf("unresolvable")
	}
	solver := casso.NewSolver()

	for _, child := range c.children {
		for _, con := range child.getInitConstraints() {
			solver.AddConstraintWithPriority(casso.Required, con)
		}
		ce := child.getCelem()
		solver.AddConstraintWithPriority(casso.Required, casso.NewConstraint(casso.GTE, 0, ce.xpos.T(1)))
		solver.AddConstraintWithPriority(casso.Required, casso.NewConstraint(casso.GTE, 0, ce.ypos.T(1)))
		solver.AddConstraintWithPriority(casso.Required, casso.NewConstraint(casso.LTE, -float64(c.width), ce.xpos.T(1), ce.width.T(1)))
		solver.AddConstraintWithPriority(casso.Required, casso.NewConstraint(casso.LTE, -float64(c.height), ce.ypos.T(1), ce.height.T(1)))

		// if child.CWidth > 0 {
		// 	solver.AddConstraintWithPriority(casso.Medium, casso.NewConstraint(casso.EQ, -float64(child.CWidth), elements[k].width.T(1)))
		// }
		// if child.Height > 0 {
		// 	solver.AddConstraintWithPriority(casso.Required, casso.NewConstraint(casso.GTE, -float64(child.Height), elements[k].height.T(1)))
		// }
	}

	for _, constraint := range c.constraints {
		if err := c.parseConstraint(solver, constraint); err != nil {
			// fmt.Printf("FUFUF: %v\n", constraint)
			// TODO: error handling
			continue
		}
	}

	for f <= c.failsMax {
		for _, v := range c.children {
			v.update(solver)
		}
		if ok := c.bias(solver); ok {
			f++
			continue
		} else {
			for _, v := range c.children {
				if !v.canBias() {
					continue
				}
				if v.anyZero() {
					return c.resolve(f + 1)
				}
			}
			break
		}
	}

	c.dirty = false
	return c.children, nil
}

func (c *constraintLayout) bias(solver *casso.Solver) bool {
	for k, child := range c.children {
		hasBias := false
		if !child.canBias() {
			continue
		}
		biases := []string{k}
		xbiases := []string{k}
		ybiases := []string{k}
		for kk, vv := range c.children {
			if !vv.canBias() {
				continue
			}
			if slices.Contains(biases, kk) {
				continue
			}

			for _, b := range biases {
				xbias := 0.0
				ybias := 0.0
				v := c.children[b]

				if v.xpos() >= vv.xpos() && v.xpos() < vv.xend() {
					xbias = float64(vv.xend()) - float64(v.xpos())
				} else if v.xpos() < vv.xpos() && v.xend() > vv.xpos() {
					xbias = float64(v.xend()) - float64(vv.xpos())
				}

				if v.ypos() >= vv.ypos() && v.ypos() < vv.yend() {
					ybias = float64(vv.yend()) - float64(v.ypos())
				} else if v.ypos() < vv.ypos() && v.yend() > vv.ypos() {
					ybias = float64(v.yend()) - float64(vv.ypos())
				}

				if xbias == 0 || ybias == 0 {
					continue
				}
				biases = append(biases, kk)

				if ybias < xbias {
					ybiases = append(ybiases, kk)
				} else {
					xbiases = append(xbiases, kk)
				}
			}
		}
		if l := len(xbiases); l > 1 {
			hasBias = true
			sort.Slice(xbiases, func(i, j int) bool {
				return c.children[xbiases[i]].xpos() < c.children[xbiases[j]].xpos()
			})

			total := float64(0)
			for _, id := range xbiases {
				total += float64(c.children[id].width())
				// terms = append(terms, elements[id].height.T(-1))
			}
			if total > float64(c.width) {
				total = float64(c.width)
			}
			solver.AddConstraintWithPriority(casso.Strong, casso.NewConstraint(casso.EQ, total, c.children[xbiases[0]].getCelem().width.T(-float64(l))))
			for i := 1; i < l; i++ {
				// idp := xbiases[i-1]
				id := xbiases[i]
				solver.AddConstraintWithPriority(casso.Strong, casso.NewConstraint(casso.EQ, total, c.children[id].getCelem().width.T(-float64(l))))
				// solver.AddConstraintWithPriority(casso.Strong, casso.NewConstraint(casso.EQ, 0, elements[idp].width.T(-1), elements[id].width.T(1)))
				// solver.AddConstraintWithPriority(casso.Strong, casso.NewConstraint(casso.EQ, 0, elements[idp].ypos.T(1), elements[idp].height.T(1), elements[id].ypos.T(-1)))
			}
		}
		if l := len(ybiases); l > 1 {
			hasBias = true
			sort.Slice(ybiases, func(i, j int) bool {
				return c.children[ybiases[i]].ypos() < c.children[ybiases[j]].ypos()
			})

			total := float64(0)
			for _, id := range ybiases {
				total += float64(c.children[id].height())
				// terms = append(terms, elements[id].height.T(-1))
			}
			if total > float64(c.height) {
				total = float64(c.height)
			}
			solver.AddConstraintWithPriority(casso.Strong, casso.NewConstraint(casso.EQ, total, c.children[ybiases[0]].getCelem().height.T(-float64(l))))
			for i := 1; i < l; i++ {
				// idp := ybiases[i-1]
				id := ybiases[i]
				solver.AddConstraintWithPriority(casso.Strong, casso.NewConstraint(casso.EQ, total, c.children[id].getCelem().height.T(-float64(l))))
				// solver.AddConstraintWithPriority(casso.Strong, casso.NewConstraint(casso.EQ, 0, elements[idp].height.T(-1), elements[id].height.T(1)))
				// solver.AddConstraint(casso.NewConstraint(casso.EQ, 0, elements[idp].ypos.T(1), elements[idp].height.T(1), elements[id].ypos.T(-1)))
			}
		}
		if hasBias {
			return true
		}
	}

	return false
}

func (c *constraintLayout) parseConstraint(solver *casso.Solver, constraint Constraint) error {
	if !c.children[constraint.Target].constraintTarget(constraint.TargetAttribute) {
		return nil
	}
	target, ok := c.children[constraint.Target]
	if !ok {
		return fmt.Errorf("unknown target: '%s'", constraint.Target)
	}

	terms := getAttributeTerms(constraint.TargetAttribute, target.getCelem(), -constraint.Multiplier)

	if constraint.Source == "" {
		_, err := solver.AddConstraintWithPriority(casso.Priority(constraint.Strength), casso.NewConstraint(casso.Op(constraint.Relation), -constraint.Constant, terms...))
		return err
	}

	if constraint.Source == "super" {
		var sourceVal float64

		switch constraint.SourceAttribute {
		case WIDTH, XEND:
			sourceVal = float64(c.width)
		case HEIGHT, YEND:
			sourceVal = float64(c.height)
		case XSTART, YSTART:
			sourceVal = 0
		}

		_, err := solver.AddConstraintWithPriority(casso.Priority(constraint.Strength), casso.NewConstraint(casso.Op(constraint.Relation), sourceVal*constraint.Multiplier+constraint.Constant, terms...))
		return err
	}

	source, ok := c.children[constraint.Source]
	if !ok {
		return fmt.Errorf("unknown source: '%s'", constraint.Source)
	}
	terms = append(getAttributeTerms(constraint.SourceAttribute, source.getCelem(), constraint.Multiplier), terms...)

	_, err := solver.AddConstraintWithPriority(casso.Priority(constraint.Strength), casso.NewConstraint(casso.Op(constraint.Relation), constraint.Constant, terms...))
	return err
}

func (c *constraintLayout) fromJson(p []byte) error {
	err := json.Unmarshal(p, &c.constraints)
	return err
}

func newConstraintLayout(sourceConstraints ...Constraint) *constraintLayout {
	ret := &constraintLayout{
		failsMax: 50,
		dirty:    true,
	}
	ret.constraints = append(ret.constraints, sourceConstraints...)

	return ret
}
