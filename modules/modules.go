package modules

var Modules map[string]Module = map[string]Module{}

type Module interface {
	Run() error
	Print()
	//Name() string
}
