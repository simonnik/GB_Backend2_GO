package check

// Repository holds datastore
type Repository interface {
	Check() error
}

type Check struct {
	Repository Repository
}

func NewService(repo Repository) *Check {
	return &Check{Repository: repo}
}

func (ch Check) Check() error {
	err := ch.Repository.Check()
	if err != nil {
		return err
	}

	return nil
}
