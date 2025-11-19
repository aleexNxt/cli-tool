package cli

import (
	"context"

	"github.com/spf13/cobra"

	"github.com/aleexNxt/cli-tool/internal/usecase"
)

func newTestCommand() *cobra.Command {
	var (
		coverage bool
		race     bool
		short    bool
		pkg      string
		run      string
		count    int
		parallel int
		timeout  string
	)

	testCmd := &cobra.Command{
		Use:   "test",
		Short: "Führt Tests aus",
		Long:  "Führt Go-Tests mit verschiedenen Optionen aus",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := InitDependencies()
			if err != nil {
				return err
			}

			opts := &usecase.TestOptions{
				Verbose:  verbose,
				Coverage: coverage,
				Race:     race,
				Short:    short,
				Package:  pkg,
				Run:      run,
				Count:    count,
				Parallel: parallel,
				Timeout:  timeout,
			}

			return deps.TestUseCase.Execute(context.Background(), opts)
		},
	}

	testCmd.Flags().BoolVarP(&coverage, "coverage", "c", false, "Coverage-Report erstellen")
	testCmd.Flags().BoolVarP(&race, "race", "r", false, "Race-Detector aktivieren")
	testCmd.Flags().BoolVarP(&short, "short", "s", false, "Kurze Tests ausführen")
	testCmd.Flags().StringVarP(&pkg, "package", "p", "", "Spezifisches Package testen")
	testCmd.Flags().StringVar(&run, "run", "", "Nur Tests ausführen, die dem Muster entsprechen")
	testCmd.Flags().IntVar(&count, "count", 1, "Anzahl der Test-Durchläufe")
	testCmd.Flags().IntVar(&parallel, "parallel", 0, "Parallele Ausführung (0 = CPU count)")
	testCmd.Flags().StringVarP(&timeout, "timeout", "t", "10m", "Timeout für Tests")

	// Subcommands
	testCmd.AddCommand(newTestBenchCommand())
	testCmd.AddCommand(newTestCoverageCommand())
	testCmd.AddCommand(newTestVetCommand())

	return testCmd
}

func newTestBenchCommand() *cobra.Command {
	var (
		pkg   string
		run   string
		count int
	)

	benchCmd := &cobra.Command{
		Use:   "bench",
		Short: "Führt Benchmarks aus",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := InitDependencies()
			if err != nil {
				return err
			}

			opts := &usecase.TestOptions{
				Package:  pkg,
				Run:      run,
				Count:    count,
				BenchMem: true,
			}

			return deps.TestUseCase.RunBenchmarks(context.Background(), opts)
		},
	}

	benchCmd.Flags().StringVarP(&pkg, "package", "p", "", "Spezifisches Package")
	benchCmd.Flags().StringVar(&run, "run", "", "Benchmark-Muster")
	benchCmd.Flags().IntVar(&count, "count", 1, "Anzahl der Durchläufe")

	return benchCmd
}

func newTestCoverageCommand() *cobra.Command {
	var outputFile string

	coverageCmd := &cobra.Command{
		Use:   "coverage",
		Short: "Erstellt Coverage-Report",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := InitDependencies()
			if err != nil {
				return err
			}

			if outputFile == "" {
				outputFile = "coverage.out"
			}

			return deps.TestUseCase.GenerateCoverage(context.Background(), outputFile)
		},
	}

	coverageCmd.Flags().StringVarP(&outputFile, "output", "o", "coverage.out", "Output-Datei")

	return coverageCmd
}

func newTestVetCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "vet",
		Short: "Führt go vet aus",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := InitDependencies()
			if err != nil {
				return err
			}

			return deps.TestUseCase.Vet(context.Background())
		},
	}
}
