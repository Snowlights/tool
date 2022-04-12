package vsql

type Instance struct{}

func NewInstance() *Instance {

	return &Instance{}
}

func (i *Instance) dsnList() []DSN {
	return nil
}
