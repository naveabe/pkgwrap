package tracker

type IJobstore interface {
	Add(BuildJob) error
}
