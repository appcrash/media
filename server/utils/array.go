package utils

import "fmt"

func RemoveElementFromArray[T any](slice []T, index int) ([]T, error) {
	if index > len(slice)-1 || index < 0 {
		return nil, fmt.Errorf("index out of range")
	}
	return append(slice[:index], slice[index+1:]...), nil
}
