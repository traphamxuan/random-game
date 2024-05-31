package utils

import "fmt"

// fmt.Errorf("failed to get %s. Requied by %s: %w", src, dst, err)
func Failed2SetupService(src string, dst string, err error) error {
	return fmt.Errorf("failed to get %s. Requied by %s: %w", src, dst, err)
}
