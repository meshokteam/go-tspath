// Package cmd process command line arguments
package cmd

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/joseluisq/go-tspath/pkg/replacer"
	"github.com/joseluisq/go-tspath/pkg/tsconfig"

	zglob "github.com/mattn/go-zglob"
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})

	configPath := flag.String("c", "./tsconfig.json", "Specifies the Typescript configuration file.")
	targetDir := flag.String("t", ".", "Specifies the directory containing output of the tsc compiler.")
	rootDir := flag.String("r", ".", "Specifies the directory which will be used to resolve relative paths.")
	debug := flag.Bool("debug", false, "Sets log level to debug")
	flag.Parse()

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	log.Info().Str("targetDir", *targetDir).Str("rootDir", *rootDir).Msg("start processing")
	startTime := time.Now()

	config := tsconfig.New(*configPath).Read()

	outFilesPath := filepath.Join(*targetDir, "**/*.js")
	files, err := zglob.Glob(outFilesPath)

	if err != nil {
		log.Fatal().Err(err).Msg("failed to get list of files to process")
	}

	log.Info().Str("glob", outFilesPath).Msgf("found %d files to process", len(files))

	// Create the replacement string array (pattern-replacement)
	replacements := make([]tsconfig.PathReplacement, 0, len(config.CompilerOptions.Paths))

	for keyPathStr, valuePathStr := range config.CompilerOptions.Paths {
		patternStr := strings.TrimSpace(keyPathStr)

		if len(patternStr) == 0 {
			continue
		}

		if len(valuePathStr) > 1 {
			log.Fatal().Str("pattern", patternStr).Msg("multiple replacements are not supported")
		}

		replacements = append(replacements, tsconfig.PathReplacement{
			Pattern:     []byte(patternStr),
			Replacement: []byte(filepath.Join(*rootDir, valuePathStr[0])),
		})
	}

	for _, file := range files {
		log.Debug().Str("file", file).Msg("processing file")
		replacer.Replace(file, replacements)
	}

	log.Info().Msgf("processed %d files in %v", len(files), time.Since(startTime))
}
