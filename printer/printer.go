package printer

import (
	"fmt"

	"github.com/JonathonGore/api-check/builder"
	"github.com/JonathonGore/api-check/runner"
)

func buildDescription(test builder.APITest) string {
	if len(test.Description) != 0 {
		return test.Description
	}

	return test.Hostname + test.Endpoint
}

func succeededText(succeeded bool) string {
	if succeeded {
		return "succeeded"
	}

	return "failed"
}

func printStats(successes, failures int) {
	total := successes + failures

	fmt.Printf("\n%v tests ran. %v successful. %v failures.\n", total, successes, failures)
}

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
