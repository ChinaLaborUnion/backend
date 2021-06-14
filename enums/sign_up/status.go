package signUpEnum

import enumsbase "grpc-demo/enums"

const (
	NoDone = 1
	Done = 2

)

func NewStatusEnums() enumsbase.EnumBaseInterface {
	return enumsbase.EnumBase{
		Enums: []int{1, 2},
	}
}
