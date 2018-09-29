package domain

type Imports map[string]struct{}

func NewImports() Imports {
	return make(Imports)
}

func (t Imports) Add(imports ...string) {
	for _, i := range imports {
		if i == "" {
			continue
		}
		t[i] = struct{}{}
	}
}

func (t Imports) AsSlice() []string {
	list := make([]string, 0, len(t))
	for l, _ := range t {
		list = append(list, l)
	}
	return list
}
