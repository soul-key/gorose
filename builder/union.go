package builder

type UnionItem struct {
	IBuilder
	IsUnionAll bool
}
type UnionClause struct {
	Unions []UnionItem
}

func (u *UnionClause) Union(b ...IBuilder) *UnionClause {
	for _, v := range b {
		u.Unions = append(u.Unions, UnionItem{IBuilder: v})
	}
	return u
}

func (u *UnionClause) UnionAll(b ...IBuilder) *UnionClause {
	for _, v := range b {
		u.Unions = append(u.Unions, UnionItem{IBuilder: v, IsUnionAll: true})
	}
	return u
}
