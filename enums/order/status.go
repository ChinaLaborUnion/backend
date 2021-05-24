package orderEnum

import enumsbase "grpc-demo/enums"

const (
	WaitToPay = 1
	Done = 2
	Cancel = 4

)

func NewStatusEnums() enumsbase.EnumBaseInterface {
	return enumsbase.EnumBase{
		Enums: []int{1, 2,4},
	}
}