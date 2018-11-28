package migrates

type CreateUsersIndexes001 struct {
}

func (mig *CreateUsersIndexes001) Up() error {

	return nil
}

func NewCreateUsersIndexes001() *CreateUsersIndexes001 {
	return &CreateUsersIndexes001{}
}
