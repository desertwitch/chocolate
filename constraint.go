package chocolate

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lithdew/casso"
)

type ConstraintAttribute uint8

func (ca *ConstraintAttribute) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch strings.ToUpper(v) {
	case "WIDTH":
		*ca = WIDTH
	case "HEIGHT":
		*ca = HEIGHT
	case "XSTART":
		*ca = XSTART
	case "YSTART":
		*ca = YSTART
	case "XEND":
		*ca = XEND
	case "YEND":
		*ca = YEND
	default:
		return fmt.Errorf("unknown attribute '%s'", v)
	}

	return nil
}

const (
	WIDTH ConstraintAttribute = iota
	HEIGHT
	XSTART
	YSTART
	XEND
	YEND
)

type ConstraintRelation uint8

func (cr *ConstraintRelation) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch strings.ToUpper(v) {
	case "EQ":
		*cr = EQ
	case "GE":
		*cr = GE
	case "LE":
		*cr = LE
	default:
		return fmt.Errorf("unknown relation '%s'", v)
	}

	return nil
}

const (
	EQ ConstraintRelation = ConstraintRelation(casso.EQ)
	GE ConstraintRelation = ConstraintRelation(casso.GTE)
	LE ConstraintRelation = ConstraintRelation(casso.LTE)
)

type ConstraintStrength float64

func (cs *ConstraintStrength) UnmarshalJSON(data []byte) error {
	var v string
	if err := json.Unmarshal(data, &v); err != nil {
		return err
	}

	switch strings.ToUpper(v) {
	case "WEAK":
		*cs = WEAK
	case "MEDIUM":
		*cs = MEDIUM
	case "STRONG":
		*cs = STRONG
	case "REQUIRED":
		*cs = REQUIRED
	default:
		return fmt.Errorf("unknown relation '%s'", v)
	}

	return nil
}

const (
	WEAK     ConstraintStrength = ConstraintStrength(casso.Weak)
	MEDIUM                      = ConstraintStrength(casso.Medium)
	STRONG                      = ConstraintStrength(casso.Strong)
	REQUIRED                    = ConstraintStrength(casso.Required)
)

func getAttributeTerms(c ConstraintAttribute, v constraintElement, m float64) []casso.Term {
	switch c {
	case WIDTH:
		return []casso.Term{v.width.T(m)}
	case HEIGHT:
		return []casso.Term{v.height.T(m)}
	case XSTART:
		return []casso.Term{v.xpos.T(m)}
	case YSTART:
		return []casso.Term{v.ypos.T(m)}
	case XEND:
		return []casso.Term{v.width.T(m), v.xpos.T(m)}
	case YEND:
		return []casso.Term{v.height.T(m), v.ypos.T(m)}
	}

	return []casso.Term{}
}

type Constraint struct {
	Target          string              `json:"target"`
	Source          string              `json:"source"`
	TargetAttribute ConstraintAttribute `json:"target_attribute"`
	SourceAttribute ConstraintAttribute `json:"source_attribute"`
	Relation        ConstraintRelation  `json:"relation"`
	Constant        float64             `json:"constant"`
	Multiplier      float64             `json:"multiplier"`
	Strength        ConstraintStrength  `json:"strength"`
}

func (c Constraint) WithTarget(v string) Constraint {
	c.Target = v
	return c
}

func (c Constraint) WithSource(v string) Constraint {
	c.Source = v
	return c
}

func (c Constraint) WithTargetAttribute(v ConstraintAttribute) Constraint {
	c.TargetAttribute = v
	return c
}

func (c Constraint) WithSourceAttribute(v ConstraintAttribute) Constraint {
	c.SourceAttribute = v
	return c
}

func (c Constraint) WithRelation(v ConstraintRelation) Constraint {
	c.Relation = v
	return c
}

func (c Constraint) WithConstant(v float64) Constraint {
	c.Constant = v
	return c
}

func (c Constraint) WithMultiplier(v float64) Constraint {
	c.Multiplier = v
	return c
}

func (c Constraint) WithStrength(v ConstraintStrength) Constraint {
	c.Strength = v
	return c
}

func (c *Constraint) UnmarshalJSON(data []byte) error {
	c.Source = ""
	c.Constant = 0
	c.Multiplier = 1.0
	c.Strength = MEDIUM

	type Alias Constraint
	v := (*Alias)(c)

	return json.Unmarshal(data, &v)
}

func NewConstraint() Constraint {
	return Constraint{
		Constant:   0,
		Multiplier: 1.0,
		Strength:   MEDIUM,
	}
}
