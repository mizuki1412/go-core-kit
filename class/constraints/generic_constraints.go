package constraints

type Number interface {
	Integer | Float
}

type Integer interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

type Float interface {
	float32 | float64
}

type IsValidInterface interface {
	IsValid() bool
}

type SetDBDriverInterface interface {
	SetDBDriver(driver string)
}
