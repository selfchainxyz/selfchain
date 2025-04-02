package cmd

import (
	"encoding/json"
	"fmt"
	"runtime"
	"strings"

	"github.com/cosmos/cosmos-sdk/version"
	"github.com/spf13/cobra"

	"selfchain/app"
)

const (
	AppVersion = "v1.0.1" // App version
	ChainID    = "self-1" // Chain ID
)

type VersionInfo struct {
	Name              string `json:"name"`
	AppName           string `json:"app_name"`
	Version           string `json:"version"`
	ChainID           string `json:"chain_id"`
	Bech32Prefix      string `json:"bech32_prefix"`
	GitCommit         string `json:"git_commit"`
	GoVersion         string `json:"go_version"`
	BuildTags         string `json:"build_tags"`
	CosmosSdkVersion  string `json:"cosmos_sdk_version"`
	TendermintVersion string `json:"tendermint_version"`
}

// VersionCmd returns a CLI command to display the application binary version information.
func VersionCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the application binary version information",
		Long: `Print the application binary version information.
Use --json flag for JSON output.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			verInfo := VersionInfo{
				Name:              app.Name,
				AppName:           "selfchaind",
				Version:           AppVersion,
				ChainID:           ChainID,
				Bech32Prefix:      app.AccountAddressPrefix,
				GitCommit:         version.Commit,
				GoVersion:         fmt.Sprintf("go version %s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH),
				BuildTags:         version.BuildTags,
				CosmosSdkVersion:  "v0.47.10", // From go.mod
				TendermintVersion: "v0.37.4",  // From go.mod
			}

			jsonFormat, err := cmd.Flags().GetBool("json")
			if err != nil {
				return err
			}

			if jsonFormat {
				output, err := json.MarshalIndent(verInfo, "", "  ")
				if err != nil {
					return err
				}
				fmt.Println(string(output))
				return nil
			}

			// Default human-readable output
			fmt.Printf("Name: %s\n", verInfo.Name)
			fmt.Printf("App Name: %s\n", verInfo.AppName)
			fmt.Printf("Version: %s\n", verInfo.Version)
			fmt.Printf("Chain ID: %s\n", verInfo.ChainID)
			fmt.Printf("Bech32 Prefix: %s\n", verInfo.Bech32Prefix)
			fmt.Printf("Git Commit: %s\n", strings.TrimSpace(verInfo.GitCommit))
			fmt.Printf("Go Version: %s\n", verInfo.GoVersion)
			if verInfo.BuildTags != "" {
				fmt.Printf("Build Tags: %s\n", verInfo.BuildTags)
			}
			fmt.Printf("Cosmos SDK Version: %s\n", verInfo.CosmosSdkVersion)
			fmt.Printf("Tendermint Version: %s\n", verInfo.TendermintVersion)
			return nil
		},
	}

	cmd.Flags().Bool("json", false, "Output version info in JSON format")
	return cmd
}
