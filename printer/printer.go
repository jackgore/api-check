package printer

import (
	"fmt"

	"github.com/JonathonGore/api-check/builder"
	"github.com/JonathonGore/api-check/runner"
)

// buildDescription builds the describing text to use when printing the run
// results. If the individual test has a description it is used otherwise
// that hostname and endpoint are used.
func buildDescription(test builder.APITest) string {
	if len(test.Description) != 0 {
		return test.Description
	}

	return test.Hostname + test.Endpoint
}

// succeededText converts the given boolean into a string representation
// of "succeeded" or "failed".
func succeededText(succeeded bool) string {
	if succeeded {
		return "succeeded"
	}

	return "failed"
}

// printStats prints the statistics from all tests that were run. Describing
// how many tests ran and how many failed/succeeded.
func printStats(successes, failures int) {
	total := successes + failures

	fmt.Printf("\n%v tests ran. %v successful. %v failures.\n", total, successes, failures)
}

// printReport consumes a RunReport for a specific test and prints information
// regarding its success or failure.
func printReport(report runner.RunReport) {
	fmt.Printf("API Check Test for: %v %v\n", buildDescription(report.Test), succeededText(report.Successful))
	if !report.Successful {
		fmt.Printf("Failure reason: %v\n", report.Error)
	}
}

// PrintReports consumes a slice of run reports and iterates through them
// printing results of each and printing aggregate result at the end.
func PrintReports(reports []runner.RunReport) {
	successes := 0
	errors := 0

	for _, report := range reports {
		printReport(report)

		if report.Error != nil {
			errors++
		} else {
			successes++
		}
	}

	printStats(successes, errors)
}
