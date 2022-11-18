package property

type AppType int

const (
	Backend  AppType = iota + 1
	Frontend         = 2
)

var AppTypeValues = []AppType{
	Backend,
	Frontend,
}
