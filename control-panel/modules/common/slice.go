package commonmodule

import "slices"

func RemoveFromUnorderedSlice[T any](
	slice []T, //nullable
	index int,
) []T {
	if slice == nil || len(slice) < 2 {
		return []T{}
	}

	slice[index] = slice[len(slice)-1]
	return slice[:len(slice)-1]
}

func RemoveFromOrderedSlice[T any](
	slice []T, //nullable
	index int,
) []T {
	if slice == nil || len(slice) < 2 {
		return []T{}
	}

	sliceA := slice[:index]
	sliceB := slice[index+1:]

	return slices.Concat(sliceA, sliceB)
}
