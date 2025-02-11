package cmd

import (
	"errors"
	"os"
	"time"

	"github.com/avast/retry-go"
	"github.com/majd/ipatool/v2/pkg/appstore"
	"github.com/schollz/progressbar/v3"
	"github.com/spf13/cobra"
)

// nolint:wrapcheck
func downloadCmd() *cobra.Command {
	var (
		acquireLicense bool
		outputPath     string
		appID          int64
		bundleID       string
	)

	cmd := &cobra.Command{
		Use:   "download",
		Short: "Download (encrypted) iOS app packages from the App Store",
		RunE: downloadRun,
	}

	cmd.Flags().Int64VarP(&appID, "app-id", "i", 0, "ID of the target iOS app (required)")
	cmd.Flags().StringVarP(&bundleID, "bundle-identifier", "b", "", "The bundle identifier of the target iOS app (overrides the app ID)")
	cmd.Flags().StringVarP(&outputPath, "output", "o", "", "The destination path of the downloaded app package")
	cmd.Flags().BoolVar(&acquireLicense, "purchase", false, "Obtain a license for the app if needed")

	return cmd
}

func downloadRun(cmd *cobra.Command, args []string) error {
	if bundleID == "" {
		return fmt.Errorf("bundle identifier is required")
	}

	// Check if the app is already purchased
	isPurchased, err := store.IsPurchased(bundleID)
	if err != nil {
		return err
	}

	// If not purchased, attempt to purchase first
	if !isPurchased {
		fmt.Printf("App %s is not purchased. Attempting to purchase...\n", bundleID)
		if err := store.Purchase(bundleID); err != nil {
			return fmt.Errorf("failed to purchase app: %v", err)
		}
		fmt.Printf("Successfully purchased %s\n", bundleID)
	}

	// Proceed with download
	return store.Download(bundleID, output)
}
