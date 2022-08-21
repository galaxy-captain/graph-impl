package recommend

import "grahp-impl/pkg/schema"

type Attribute struct {
	Name  string
	Model string
	Count int
}

type Exchange struct {
	input  []*schema.Unit
	output []*schema.Unit
}
