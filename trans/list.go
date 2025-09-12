// Package trans ...
package trans

import (
	"context"
)

// PageList ...
func PageList[T string | uint64](ctx context.Context, list []T, pageNum int, pageSize int) (int, []T) {
	total := len(list)
	start := pageSize * pageNum
	end := pageSize * (pageNum + 1)

	if start >= total {
		return total, list[:0]
	}

	if end > total {
		return total, list[start:]
	}

	return total, list[start:end]
}

// PageListIntf ...
func PageListIntf[T string | uint64](ctx context.Context, list []T, pageNum int, pageSize int) (int, []any) {
	total := len(list)
	start := pageSize * pageNum
	end := pageSize * (pageNum + 1)

	if start >= total {
		return total, ListIntf(list[:0])
	}

	if end > total {
		return total, ListIntf(list[start:])
	}

	return total, ListIntf(list[start:end])
}

func ListIntf[T any](list []T) []any {
	ll := make([]any, len(list))
	for i, v := range list {
		ll[i] = v
	}
	return ll
}
