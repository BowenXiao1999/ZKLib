package zetcd


var (
	errorCodeToErr = map[ErrCode]error{
		errBadArguments: ErrBadArguments,
		errNoNode:       ErrNoNode,
		errNodeExists:   ErrNodeExists,
		errInvalidAcl:   ErrInvalidACL,
		errBadVersion:   ErrBadVersion,
		errAPIError:     ErrAPIError,
		errNotEmpty:     ErrNotEmpty,
	}

)