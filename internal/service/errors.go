package service

import "errors"

var (
	ErrCannotSignToken  = errors.New("cannot sign token")
	ErrCannotParseToken = errors.New("cannot parse token")

	ErrPingerAlreadyExists = errors.New("pinger already exists")
	ErrCannotCreatePinger  = errors.New("cannot create pinger")
	ErrPingerNotFound      = errors.New("pinger not found")
	ErrCannotGetPinger     = errors.New("cannot get pinger")

	ErrCannotConvertJson = errors.New("cannot convert json")
	ErrCannotStoreReport = errors.New("cannot store report")
	ErrCannotGetReports  = errors.New("cannot get reports")
)
