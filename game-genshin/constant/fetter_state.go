package constant

type FetterState struct {
	NONE     uint16
	NOT_OPEN uint16
	OPEN     uint16
	FINISH   uint16
}

func GetFetterStateConst() (r *FetterState) {
	r = new(FetterState)
	r.NONE = 0
	r.NOT_OPEN = 1
	r.OPEN = 1
	r.FINISH = 3
	return r
}
