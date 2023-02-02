package entity

type Report struct {
	Id     int64
	UserId int64
	Height float64
	Weight float64
	Info   *RealInfo
}

type Blood struct {
	Pressure float64
}
