package cassowary

type Priority int64

const (
	PriorityRequired Priority = 1000000000
	PriorityStrong   Priority = 1000000
	PriorityMedium   Priority = 1000
	PriorityWeak     Priority = 1
)
