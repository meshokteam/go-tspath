// Package replacer replaces every TS path alias occurrence per file using a tsconfig.json
package replacer

import (
	"bufio"
	"bytes"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"

	"github.com/joseluisq/go-tspath/pkg/tsconfig"
	"github.com/joseluisq/redel/v3"
)

// Replace replaces every TS path occurrence per file
func Replace(filePathAbs string, replacements []tsconfig.PathReplacement) {
	r, err := os.Open(filePathAbs)

	if err != nil {
		log.Fatal().Err(err).Str("file", filePathAbs).Msg("failed to open file")
	}

	defer r.Close()

	filePathTemp := filePathAbs + ".js"

	w, err := os.Create(filePathTemp)

	if err != nil {
		log.Fatal().Err(err).Str("file", filePathTemp).Msg("failed to create temporary file")
	}

	defer w.Close()

	var writer = bufio.NewWriter(w)

	pathRel := filepath.Dir(filePathAbs)

	replaceFunc := func(data []byte, atEOF bool) {
		_, err := writer.Write(data)

		if err != nil {
			log.Fatal().Err(err).Str("file", filePathTemp).Msg("failed to write data to file")
		}

		if atEOF {
			err := os.Remove(filePathAbs)

			if err != nil {
				log.Fatal().Err(err).Str("file", filePathAbs).Msg("failed to remove old file")
			}

			err = os.Rename(filePathTemp, filePathAbs)

			if err != nil {
				log.Fatal().Err(err).Str("temp", filePathTemp).Str("file", filePathAbs).Msg("failed to rename new file")
			}

			log.Debug().Str("file", filePathAbs).Msg("file has been processed")
		}
	}

	filterFunc := func(matchValue []byte) []byte {
		for _, vtspath := range replacements {
			if len(vtspath.Replacement) == 0 {
				continue
			}

			if bytes.Equal(matchValue, vtspath.Pattern) {
				repl := string(vtspath.Replacement)
				replacement, err := filepath.Rel(pathRel, repl)

				if err != nil {
					log.Fatal().Err(err).Str("base", pathRel).Str("target", repl).Msg("failed to obtain relative path")
				}

				log.Debug().Str("file", filePathAbs).Str("from", string(matchValue)).Str("to", replacement).Msg("replacing path")

				return []byte("./" + replacement)
			}
		}

		return matchValue
	}

	rep := redel.New(r, []redel.Delimiter{
		{Start: []byte("require(\""), End: []byte("\");")},
		{Start: []byte("from \""), End: []byte("\";")},
	})

	rep.ReplaceFilterWith(replaceFunc, filterFunc, true)

	writer.Flush()
}
