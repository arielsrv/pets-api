package property

type AppType int

const (
	Backend AppType = iota + 1
	Frontend
	_Limit
)

func (a AppType) IsValid() bool {
	return a > 0 && a < _Limit
}
