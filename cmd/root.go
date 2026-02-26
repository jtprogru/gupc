package cmd

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// SLAInfo represents the JSON output structure.
type SLAInfo struct {
	SLA               float64 `json:"SLA"` //nolint,tagliatelle // This format is normal
	Nines             string  `json:"nines"`
	DailyDownSecs     float64 `json:"dailyDownSecs"`
	DailyDown         string  `json:"dailyDown"`
	WeeklyDownSecs    float64 `json:"weeklyDownSecs"`
	WeeklyDown        string  `json:"weeklyDown"`
	MonthlyDownSecs   float64 `json:"monthlyDownSecs"`
	MonthlyDown       string  `json:"monthlyDown"`
	QuarterlyDownSecs float64 `json:"quarterlyDownSecs"`
	QuarterlyDown     string  `json:"quarterlyDown"`
	YearlyDownSecs    float64 `json:"yearlyDownSecs"`
	YearlyDown        string  `json:"yearlyDown"`
}

var ninesWords = map[int]string{
	1: "one nine",
	2: "two nines",
	3: "three nines",
	4: "four nines",
	5: "five nines",
	6: "six nines",
	7: "seven nines",
	8: "eight nines",
	9: "nine nines",
}

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "gupc",
	Short: "A simple uptime calculator",
	Long: `A simple uptime and downtime calculator.
The first argument should be a SLA percentage (e.g., 99.9 or 99).`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			fmt.Println("Please provide a SLA percent (e.g., 99.9 or 99).")
			return
		}

		input := args[0]
		percent, err := strconv.ParseFloat(input, 64)
		if err != nil {
			fmt.Printf("Invalid SLA percent '%s': %v\n", input, err)
			return
		}

		if percent <= 0 || percent >= 100 {
			fmt.Println("SLA percent must be greater than 0 and less than 100.")
			return
		}

		sla := percent / 100.0

		// Compute the number of nines, if applicable.
		var ninesStr string
		if percent < 100 {
			nines := int(math.Round(-math.Log10(1 - sla)))
			if nines >= 1 && nines <= 9 {
				ninesStr = ninesWords[nines]
			}
		}

		if jsonFlag, _ := cmd.Flags().GetBool("json"); jsonFlag {
			info := SLAInfo{
				SLA:               percent,
				Nines:             ninesStr,
				DailyDownSecs:     86400 * (1 - sla),
				WeeklyDownSecs:    604800 * (1 - sla),
				MonthlyDownSecs:   30.44 * 86400 * (1 - sla), // average month
				QuarterlyDownSecs: 90 * 86400 * (1 - sla),
				YearlyDownSecs:    365.2425 * 86400 * (1 - sla),
			}
			info.DailyDown = formatDuration(info.DailyDownSecs)
			info.WeeklyDown = formatDuration(info.WeeklyDownSecs)
			info.MonthlyDown = formatDuration(info.MonthlyDownSecs)
			info.QuarterlyDown = formatDuration(info.QuarterlyDownSecs)
			info.YearlyDown = formatDuration(info.YearlyDownSecs)

			out, err := json.MarshalIndent(info, "", "    ")
			if err != nil {
				fmt.Println("Error generating JSON:", err)
				return
			}
			fmt.Println(string(out))
		} else {
			// Human‑readable output
			if ninesStr != "" {
				fmt.Printf("%.2f%% SLA (%s)\n", percent, ninesStr)
			} else {
				fmt.Printf("%.2f%% SLA\n", percent)
			}
			fmt.Println("Daily downtime:", formatDuration(86400*(1-sla)))
			fmt.Println("Weekly downtime:", formatDuration(604800*(1-sla)))
			fmt.Println("Monthly downtime:", formatDuration(30.44*86400*(1-sla)))
			fmt.Println("Quarterly downtime:", formatDuration(90*86400*(1-sla)))
			fmt.Println("Yearly downtime:", formatDuration(365.2425*86400*(1-sla)))
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().Bool("json", false, "Output in JSON format")
}

// formatDuration converts a duration in seconds to a human‑readable string.
func formatDuration(seconds float64) string {
	totalSec := int64(seconds)
	h := totalSec / 3600
	m := (totalSec % 3600) / 60
	s := totalSec % 60
	parts := []string{}
	if h > 0 {
		parts = append(parts, fmt.Sprintf("%dh", h))
	}
	if m > 0 {
		parts = append(parts, fmt.Sprintf("%dm", m))
	}
	if s > 0 {
		parts = append(parts, fmt.Sprintf("%ds", s))
	}
	if len(parts) == 0 {
		return "0s"
	}
	return strings.Join(parts, " ")
}
