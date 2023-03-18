package dao

import "gorm.io/gen"

type Hospital interface {

	//GetNearHospital
	//
	//select temp.*, f.*, fd.* from(
	//	select *, (st_distance(point(longitude,latitude),point(10.600000, 10.500000))*111195/1000) as distance
	//	from hospital where id > @last order by distance
	//) as temp left join file_data fd on fd.id = temp.avatar left join file f on f.id = fd.id where distance < @rge limit 20
	GetNearHospital(g1 float64, g2 float64, last int64, rge int) (*gen.T, error)
}
