package column

import "medicinal_share/main/middleware"

type NotFound struct {
	col string
}

func (c NotFound) Error() string {
	return c.col + "未找到"
}

func Factory(typ string) Column {
	switch typ {
	case "text":
		return &TextColumn{}
	case "range":
		return &RangeColumn{}
	case "list":
		return &ListColumn{}
	default:
		panic(middleware.NewCustomErr(middleware.ERROR, NotFound{
			col: typ,
		}.Error()))
	}
}
