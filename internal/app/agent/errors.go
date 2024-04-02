package agent

import "errors"

var (
	errClientDoErr              = errors.New("CLIENT Client Do ERROR")
	errEncryptSHA256Err         = errors.New("CLIENT EncryptSHA256 ERROR")
	errNewRequestWithContextErr = errors.New("CLIENT New Request Body Create ERROR")
)
