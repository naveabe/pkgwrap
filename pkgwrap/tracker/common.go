package tracker

type IJobstore interface {
	AddJob(BuildJob) error
}
