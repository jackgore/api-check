package printer

import (
	"fmt"

	"github.com/JonathonGore/api-check/runner"
)

func succeededText(succeeded bool) string {
	if succeeded {
		return "succeeded"
	}

	return "failed"
}

func printReport(report runner.RunReport) {
	fmt.Printf("API Check Test for: %v%v %v\n", report.Test.Hostname,
		report.Test.Endpoint, succeededText(report.Successful))
	if !report.Successful {
		fmt.Printf("Failure reason: %v\n", report.Error)
	}
}

func PrintReports(reports []runner.RunReport) {
	for _, report := range reports {
		printReport(report)
	}
}
