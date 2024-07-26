package packet

type route func() // put better function definition here

var routers = make(map[[2]uint16]route)

func RegisterRoute(state uint16, id uint16, route route) {
	routers[[2]uint16{state, id}] = route
}

func GetRoute(state uint16, id uint16) (route, bool) {
	route, ok := routers[[2]uint16{state, id}]
	return route, ok
}
