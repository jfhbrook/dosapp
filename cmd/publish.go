/*
Copyright © 2024 Josh Holbrook <josh.holbrook@gmail.com>
*/
package cmd

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"github.com/jfhbrook/dosapp/config"
	"github.com/jfhbrook/dosapp/manifest"
)

var publishCmd = &cobra.Command{
	Use:   "publish",
	Short: "Publish the current app",
	Long:  `Publish the current app by creating and uploading a GitHub release.`,
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.NewConfig()
		man, err := manifest.FromFile("./Package.yml")

		if err != nil {
			log.Fatal().Err(err).Msg("Failed to read manifest")
		}

		var missingFields []string

		if man.Name == "" {
			missingFields = append(missingFields, "name")
		}

		if man.Version == nil {
			missingFields = append(missingFields, "version")
		}

		if man.ReleaseVersion == nil {
			missingFields = append(missingFields, "release_version")
		}

		if len(missingFields) > 0 {
			log.Fatal().Strs("fields", missingFields).Msg("Missing required fields")
		}

		tag := fmt.Sprintf("%s-%s-%s", man.Name, man.Version, man.ReleaseVersion)
		tarFileName := fmt.Sprintf("%s.tar.gz", tag)

		var tarBuf bytes.Buffer
		var gzBuf bytes.Buffer
		tw := tar.NewWriter(&tarBuf)

		err = filepath.Walk("./", func(path string, info fs.FileInfo, err error) error {
			log.Warn().Msg(path)
			if err != nil {
				return err
			}

			if info.IsDir() || path == tarFileName {
				return nil
			}

			log.Info().Str("path", path).Msg("Adding file to artifact")

			data, err := os.ReadFile(path)

			if err != nil {
				return err
			}

			hdr := &tar.Header{
				Name:    filepath.Join(tag, path),
				Mode:    int64(info.Mode()),
				ModTime: info.ModTime(),
				Size:    info.Size(),
			}

			if err := tw.WriteHeader(hdr); err != nil {
				return err
			}

			if _, err := tw.Write(data); err != nil {
				return err
			}

			return nil
		})

		if err != nil {
			log.Fatal().Err(err).Msg("Failed to create artifact")
		}

		if err := tw.Close(); err != nil {
			log.Fatal().Err(err).Msg("Failed to create artifact")
		}

		zw := gzip.NewWriter(&gzBuf)

		if _, err := zw.Write(tarBuf.Bytes()); err != nil {
			log.Fatal().Err(err).Msg("Failed to compress artifact")
		}

		if err := zw.Close(); err != nil {
			log.Fatal().Err(err).Msg("Failed to compress artifact")
		}

		err = os.WriteFile(tarFileName, gzBuf.Bytes(), 0644)

		if err != nil {
			log.Fatal().Err(err).Msg("Failed to write artifact to disk")
		}

		log.Info().Str("file", tarFileName).Msg("Artifact created")
		log.Info().Msg("Publishing release with 'gh'...")
		gh := exec.Command(
			"gh",
			"release",
			"create",
			"-t", fmt.Sprintf(
				"%s v%s - release %s",
				man.Name, man.Version, man.ReleaseVersion,
			),
			"--notes", fmt.Sprintf(
				"Package %s v%s, release %s",
				man.Name, man.Version, man.ReleaseVersion,
			),
			tarFileName,
		)

		gh.Env = cfg.Environ()
		gh.Stdout = os.Stdout
		gh.Stderr = os.Stderr

		err = gh.Run()

		if err != nil {
			log.Fatal().Err(err).Msg("Failed to publish release")
		}
		log.Info().Msg("Release published")
	},
}

func init() {
	rootCmd.AddCommand(publishCmd)
}
