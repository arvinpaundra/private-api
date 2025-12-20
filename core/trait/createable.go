package trait

type Createable struct {
	val bool
}

func (c *Createable) IsCreated() bool {
	return c.val
}

func (c *Createable) MarkCreate() {
	c.val = true
}

func (c *Createable) UnmarkCreate() {
	c.val = false
}
