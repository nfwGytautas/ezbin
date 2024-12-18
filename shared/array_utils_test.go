package shared_test

import (
	"testing"

	"github.com/nfwGytautas/ezbin/shared"
)

func TestArrayContains(t *testing.T) {
	{
		arrayInt := []int{1, 2, 3, 4, 5}

		if !shared.ArrayContains(arrayInt, 3) {
			t.Errorf("ArrayContains failed")
		}
		if shared.ArrayContains(arrayInt, 6) {
			t.Errorf("ArrayContains failed")
		}
	}

	{
		arrayString := []string{"a", "b", "c", "d", "e"}

		if !shared.ArrayContains(arrayString, "c") {
			t.Errorf("ArrayContains failed")
		}

		if shared.ArrayContains(arrayString, "f") {
			t.Errorf("ArrayContains failed")
		}
	}
}
