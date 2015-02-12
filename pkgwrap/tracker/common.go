package tracker

type IJobstore interface {
	AddJob(BuildJob) (string, error)
}
