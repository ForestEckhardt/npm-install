package npminstall

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/paketo-buildpacks/packit/pexec"
	"github.com/paketo-buildpacks/packit/scribe"
)

func NewInstallBuildProcess(executable Executable, environment EnvironmentConfig, logger scribe.Emitter) InstallBuildProcess {
	return InstallBuildProcess{
		executable:  executable,
		environment: environment,
		logger:      logger,
	}
}

type InstallBuildProcess struct {
	executable  Executable
	environment EnvironmentConfig
	logger      scribe.Emitter
}

func (r InstallBuildProcess) ShouldRun(workingDir string, metadata map[string]interface{}) (bool, string, error) {
	return true, "", nil
}

func (r InstallBuildProcess) Run(modulesDir, cacheDir, workingDir string) error {
	err := os.Mkdir(filepath.Join(modulesDir, "node_modules"), os.ModePerm)
	if err != nil {
		return err
	}

	err = os.Symlink(filepath.Join(modulesDir, "node_modules"), filepath.Join(workingDir, "node_modules"))
	if err != nil {
		return err
	}

	buffer := bytes.NewBuffer(nil)
	args := []string{"install", "--unsafe-perm", "--cache", cacheDir}

	r.logger.Subprocess("Running 'npm %s'", strings.Join(args, " "))
	err = r.executable.Execute(pexec.Execution{
		Args:   args,
		Dir:    workingDir,
		Stdout: buffer,
		Stderr: buffer,
		Env: append(
			os.Environ(),
			fmt.Sprintf("NPM_CONFIG_LOGLEVEL=%s", r.environment.GetValue("NPM_CONFIG_LOGLEVEL")),
		),
	})
	if err != nil {
		r.logger.Subprocess("%s", buffer.String())
		return fmt.Errorf("npm install failed: %w", err)
	}

	return nil
}
