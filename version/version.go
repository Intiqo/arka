package version

import "fmt"

var (
	BuildVersion = "development"
	BuildNumber  = "development"
	CommitHash   = "development"
)

func PrintInfo() {
	fmt.Printf("Build Version: %s\n", BuildVersion)
	fmt.Printf("Build Number: %s\n", BuildNumber)
	fmt.Printf("Commit Hash: %s\n", CommitHash)
}
