package cli

import (
	"context"
	"strings"

	"github.com/spf13/cobra"

	"github.com/aleexNxt/cli-tool/internal/usecase"
)

func newDockerCommand() *cobra.Command {
	dockerCmd := &cobra.Command{
		Use:   "docker",
		Short: "Docker-Operationen",
		Long:  "Verwaltet Docker Images und Container",
	}

	// Subcommands
	dockerCmd.AddCommand(newDockerBuildCommand())
	dockerCmd.AddCommand(newDockerPushCommand())
	dockerCmd.AddCommand(newDockerRunCommand())
	dockerCmd.AddCommand(newDockerStopCommand())
	dockerCmd.AddCommand(newDockerRemoveCommand())
	dockerCmd.AddCommand(newDockerImagesCommand())
	dockerCmd.AddCommand(newDockerPsCommand())

	return dockerCmd
}

func newDockerBuildCommand() *cobra.Command {
	var (
		imageName  string
		tag        string
		dockerfile string
		context    string
		noCache    bool
		platform   string
		buildArgs  []string
	)

	buildCmd := &cobra.Command{
		Use:   "build",
		Short: "Build ein Docker Image",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := InitDependencies()
			if err != nil {
				return err
			}

			// Parse build args
			buildArgsMap := make(map[string]string)
			for _, arg := range buildArgs {
				parts := strings.SplitN(arg, "=", 2)
				if len(parts) == 2 {
					buildArgsMap[parts[0]] = parts[1]
				}
			}

			opts := &usecase.DockerBuildOptions{
				ImageName:  imageName,
				Tag:        tag,
				Dockerfile: dockerfile,
				Context:    context,
				NoCache:    noCache,
				Platform:   platform,
				BuildArgs:  buildArgsMap,
			}

			return deps.DockerUseCase.Build(cmd.Context(), opts)
		},
	}

	buildCmd.Flags().StringVarP(&imageName, "image", "i", "", "Image-Name (erforderlich)")
	buildCmd.Flags().StringVarP(&tag, "tag", "t", "latest", "Image-Tag")
	buildCmd.Flags().StringVarP(&dockerfile, "file", "f", "Dockerfile", "Pfad zum Dockerfile")
	buildCmd.Flags().StringVarP(&context, "context", "C", ".", "Build-Context")
	buildCmd.Flags().BoolVar(&noCache, "no-cache", false, "Cache nicht verwenden")
	buildCmd.Flags().StringVar(&platform, "platform", "", "Ziel-Platform")
	buildCmd.Flags().StringArrayVar(&buildArgs, "build-arg", []string{}, "Build-Argumente (KEY=VALUE)")
	buildCmd.MarkFlagRequired("image")

	return buildCmd
}

func newDockerPushCommand() *cobra.Command {
	var (
		imageName string
		tag       string
	)

	pushCmd := &cobra.Command{
		Use:   "push",
		Short: "Push ein Image zur Registry",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := InitDependencies()
			if err != nil {
				return err
			}

			return deps.DockerUseCase.Push(context.Background(), imageName, tag)
		},
	}

	pushCmd.Flags().StringVarP(&imageName, "image", "i", "", "Image-Name (erforderlich)")
	pushCmd.Flags().StringVarP(&tag, "tag", "t", "latest", "Image-Tag")
	pushCmd.MarkFlagRequired("image")

	return pushCmd
}

func newDockerRunCommand() *cobra.Command {
	var (
		imageName     string
		containerName string
		ports         []string
		env           []string
	)

	runCmd := &cobra.Command{
		Use:   "run",
		Short: "Startet einen Container",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := InitDependencies()
			if err != nil {
				return err
			}

			// Parse environment variables
			envMap := make(map[string]string)
			for _, e := range env {
				parts := strings.SplitN(e, "=", 2)
				if len(parts) == 2 {
					envMap[parts[0]] = parts[1]
				}
			}

			return deps.DockerUseCase.Run(context.Background(), imageName, containerName, ports, envMap)
		},
	}

	runCmd.Flags().StringVarP(&imageName, "image", "i", "", "Image-Name (erforderlich)")
	runCmd.Flags().StringVarP(&containerName, "name", "n", "", "Container-Name")
	runCmd.Flags().StringArrayVarP(&ports, "port", "p", []string{}, "Port-Mappings (z.B. 8080:80)")
	runCmd.Flags().StringArrayVarP(&env, "env", "e", []string{}, "Umgebungsvariablen (KEY=VALUE)")
	runCmd.MarkFlagRequired("image")

	return runCmd
}

func newDockerStopCommand() *cobra.Command {
	var containerName string

	stopCmd := &cobra.Command{
		Use:   "stop",
		Short: "Stoppt einen Container",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := InitDependencies()
			if err != nil {
				return err
			}

			return deps.DockerUseCase.Stop(context.Background(), containerName)
		},
	}

	stopCmd.Flags().StringVarP(&containerName, "name", "n", "", "Container-Name (erforderlich)")
	stopCmd.MarkFlagRequired("name")

	return stopCmd
}

func newDockerRemoveCommand() *cobra.Command {
	var (
		containerName string
		force         bool
	)

	rmCmd := &cobra.Command{
		Use:   "rm",
		Short: "Entfernt einen Container",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := InitDependencies()
			if err != nil {
				return err
			}

			return deps.DockerUseCase.Remove(context.Background(), containerName, force)
		},
	}

	rmCmd.Flags().StringVarP(&containerName, "name", "n", "", "Container-Name (erforderlich)")
	rmCmd.Flags().BoolVarP(&force, "force", "f", false, "Erzwinge Entfernung")
	rmCmd.MarkFlagRequired("name")

	return rmCmd
}

func newDockerImagesCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "images",
		Short: "Listet Docker Images",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := InitDependencies()
			if err != nil {
				return err
			}

			return deps.DockerUseCase.ListImages(context.Background())
		},
	}
}

func newDockerPsCommand() *cobra.Command {
	var all bool

	psCmd := &cobra.Command{
		Use:   "ps",
		Short: "Listet Docker Container",
		RunE: func(cmd *cobra.Command, args []string) error {
			deps, err := InitDependencies()
			if err != nil {
				return err
			}

			return deps.DockerUseCase.ListContainers(context.Background(), all)
		},
	}

	psCmd.Flags().BoolVarP(&all, "all", "a", false, "Zeige alle Container")

	return psCmd
}
