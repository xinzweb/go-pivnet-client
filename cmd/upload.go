package cmd

import (
	"github.com/baotingfang/go-pivnet-client/vlog"
	"sync"

	"github.com/spf13/cobra"
)

var (
	uploadCmdFlagsInit sync.Once

	metaDataFilePath string
	searchPath       string
	verbose          bool
	gpdbVersion      string

	logLevel = vlog.InfoLevel
)

var uploadCmd = &cobra.Command{
	Use:   "upload [-v] [-s search_path] <-m metadata_file> <-g gpdb_version>",
	Short: "Upload artifacts to pivnet",
	Long:  `Given metadata specifying a pivnet release with file groups and/or product files, this program will perform the necessary actions to create those components on pivnet`,
	Run: func(cmd *cobra.Command, args []string) {
		vlog.InitLog("Upload ", logLevel)
	},
}

func init() {
	uploadCmdFlagsInit.Do(func() {
		uploadCmd.Flags().StringVarP(&metaDataFilePath, FlagNameMetaFilePath.String(), "m", "", "Path to a valid pivnet client metadata yaml file")
		uploadCmd.Flags().StringVarP(&searchPath, FlagNameSearchPath.String(), "s", ".", "Path to look for product files defined in metadata")
		uploadCmd.Flags().BoolVarP(&verbose, FlagNameVerbose.String(), "v", false, "Verbose output")
		uploadCmd.Flags().StringVarP(&gpdbVersion, FlagNameGpdbVersion.String(), "g", "", "GPDB version from getversion tool")

		uploadCmdRequiredFlags := []string{
			FlagNameMetaFilePath.String(),
			FlagNameGpdbVersion.String(),
		}

		for _, flag := range uploadCmdRequiredFlags {
			err := uploadCmd.MarkFlagRequired(flag)
			if err != nil {
				vlog.Fatal(err)
			}
		}

		if verbose {
			logLevel = vlog.DebugLevel
		}

		rootCmd.AddCommand(uploadCmd)
	})
}
