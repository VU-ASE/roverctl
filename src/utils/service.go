package utils

import "fmt"

func ServiceFqn(author string, name string, version string) string {
	return fmt.Sprintf("%s/%s@%s", author, name, version)
}
