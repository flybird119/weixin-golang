package utils

import "testing"

func TestDemo(t *testing.T) {

	has := []int{345, 345, 3451, 4, 5, 3, 45, 34, 5, 323423445}

	has = append(has[:9], has[10:]...)

	t.Log(has)
}
