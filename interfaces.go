package chocolate

import (
	"github.com/lithdew/casso"
)

type barSizeConstraint struct {
	Relation ConstraintRelation
	Value    float64
}

type barConstrainer interface {
	sizeConstraints() (width, height []barSizeConstraint)
	constraintTarget(a ConstraintAttribute) bool
	canBias() bool
}

type barRenderer interface {
	setSize(width, height int)
	render() string
}

type barViewer interface {
	View() string
}

type barModel interface {
	View() string
	Resize(width, height int)
}

type barContainer interface {
	setDirty()
}

type barChild interface {
	barConstrainer
	barModel
	getInitConstraints() []casso.Constraint
	getCelem() constraintElement
	update(*casso.Solver)
	width() int
	height() int
	xpos() int
	ypos() int
	xend() int
	yend() int
	anyZero() bool
	setParent(parent barContainer)
}
