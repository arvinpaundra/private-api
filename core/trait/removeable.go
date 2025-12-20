package trait

type Removeable struct {
	val bool
}

func (r *Removeable) IsRemoved() bool {
	return r.val
}

func (r *Removeable) MarkRemove() {
	r.val = true
}

func (r *Removeable) UnmarkRemove() {
	r.val = false
}
