package commonmodule

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
