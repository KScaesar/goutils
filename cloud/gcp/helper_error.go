package gcp

import (
	"google.golang.org/api/compute/v1"

	"github.com/Min-Feng/goutils/errors"
)

func extractError(src *compute.Operation) error {
	if src.Status != "DONE" {
		return errors.Wrap(
			errors.ErrSystem,
			"developer should make sure that the status is DONE",
		)
	}

	if src.Error == nil {
		return nil
	}

	// example:
	// Quota 'NETWORKS' exceeded.  Limit: 5.0 globally
	if ok, msg := isTargetError(src.Error, "QUOTA_EXCEEDED"); ok {
		return errors.Wrap(errors.ErrSystem, msg)
	}

	// example:
	// Invalid IPCidrRange: 192.168.4.0/24 conflicts with existing subnetwork
	if ok, msg := isTargetError(src.Error, "INVALID_USAGE"); ok {
		return errors.Wrap(errors.ErrInvalidParams, msg)
	}

	// 未知的錯誤, 也不知道是回傳一個還是多個錯誤, 所以回傳最原始的資訊
	byteErr, _ := src.Error.MarshalJSON()
	msg := string(byteErr)
	return errors.Wrap(errors.ErrSystem, msg)
}

func isTargetError(opError *compute.OperationError, targetErrCode string) (ok bool, errMsg string) {
	for _, err := range opError.Errors {
		if err.Code == targetErrCode {
			return true, err.Message
		}
	}
	return false, ""
}
