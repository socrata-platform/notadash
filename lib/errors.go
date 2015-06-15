package lib

import (
	"errors"
)

var (
	ErrParameterMissing = errors.New("Required command parameters missing from call, please see log for more details")
	ErrMarathonError    = errors.New("An error accoured while querying the Marathon API, please see logs for more details")
	ErrMesosError       = errors.New("An error accoured while querying the Mesos API, please see logs for more details")
	ErrLeaderRequired   = errors.New("This task is required to be run on the cluster leader.")
)
