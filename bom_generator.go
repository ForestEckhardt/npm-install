package npminstall

import (
	"bytes"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/paketo-buildpacks/packit/pexec"
	"github.com/paketo-buildpacks/packit/scribe"
)

type BOMGenerator struct {
	executable  Executable
	executable2 Executable
	logger      scribe.Emitter
}

func NewBOMGenerator(executable Executable, executable2 Executable, logger scribe.Emitter) BOMGenerator {
	return BOMGenerator{
		executable:  executable,
		executable2: executable2,
		logger:      logger,
	}
}

// InstallAndRun installs the cyclonedx/bom tool for BOM creation and runs it.
// It returns the path to the resulting JSON bom file, and an error.
func (bg BOMGenerator) InstallAndRun(workingDir string) (string, error) {
	buffer := bytes.NewBuffer(nil)
	args := []string{"install", "-g", "@cyclonedx/bom"}
	bg.logger.Subprocess("Running 'npm %s'", strings.Join(args, " "))
	err := bg.executable.Execute(pexec.Execution{
		Args:   args,
		Stdout: buffer,
		Stderr: buffer,
	})
	if err != nil {
		bg.logger.Subprocess("%s", buffer.String())
		return "", fmt.Errorf("failed to install: %w", err)
	}

	bg.logger.Subprocess("Successful install of cyclonedx/bom")

	args = []string{"-o", "bom.json"}
	bg.logger.Subprocess("Running  'cyclonedx-bom %s'", strings.Join(args, " "))
	err = bg.executable2.Execute(pexec.Execution{
		Args:   args,
		Dir:    workingDir,
		Stdout: buffer,
		Stderr: buffer,
	})

	if err != nil {
		bg.logger.Subprocess("%s", buffer.String())
		return "", fmt.Errorf("failed to run: %w", err)
	}

	return filepath.Join(workingDir, "bom.json"), nil
}
