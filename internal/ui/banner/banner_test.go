package banner

import (
	"testing"
)

func TestPrint(t *testing.T) {
	// Print() writes directly to os.Stdout
	// We'll verify the function exists and doesn't panic when called
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("Print() panicked: %v", r)
		}
	}()

	// Call the function - it will output to stdout but that's acceptable for coverage
	Print()
}
