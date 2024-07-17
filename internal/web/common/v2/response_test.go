package v2

import (
	"fmt"
	"testing"
)

func Test_response_Clone(t *testing.T) {
	original := &response{Data: "Original Data"}
	clone := original.Clone().(*response)

	fmt.Printf("Original Data: %v\n", original.Data)
	fmt.Printf("Clone Data: %v\n", clone.Data)

	clone.SetData("New Data")

	fmt.Printf("After modification:\n")
	fmt.Printf("Original Data: %v\n", original.Data)
	fmt.Printf("Clone Data: %v\n", clone.Data)
}
