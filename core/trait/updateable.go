package trait

type Updateable struct {
	val bool
}

func (u *Updateable) IsUpdated() bool {
	return u.val
}

func (u *Updateable) MarkUpdate() {
	u.val = true
}

func (u *Updateable) UnmarkUpdate() {
	u.val = false
}
