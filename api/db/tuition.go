package db

type Tuition struct {
}

func (q *Querier) QueryAlertStudentInTuition() (res AllUser, err error) {
	res, err = q.GetAllUser()

	return res, err
}
