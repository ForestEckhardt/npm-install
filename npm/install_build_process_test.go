package npm_test

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/cloudfoundry/npm-cnb/npm"
	"github.com/cloudfoundry/npm-cnb/npm/fakes"
	"github.com/cloudfoundry/packit/pexec"
	"github.com/cloudfoundry/packit/scribe"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testInstallBuildProcess(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		modulesDir string
		cacheDir   string
		workingDir string
		executable *fakes.Executable
		buffer     *bytes.Buffer

		process npm.InstallBuildProcess
	)

	it.Before(func() {
		var err error
		modulesDir, err = ioutil.TempDir("", "modules")
		Expect(err).NotTo(HaveOccurred())

		cacheDir, err = ioutil.TempDir("", "cache")
		Expect(err).NotTo(HaveOccurred())

		workingDir, err = ioutil.TempDir("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		executable = &fakes.Executable{}

		buffer = bytes.NewBuffer(nil)

		process = npm.NewInstallBuildProcess(executable, scribe.NewLogger(buffer))
	})

	it.After(func() {
		Expect(os.RemoveAll(modulesDir)).To(Succeed())
		Expect(os.RemoveAll(cacheDir)).To(Succeed())
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	context("ShouldRun", func() {
		it("returns true", func() {
			run, sha, err := process.ShouldRun(workingDir, nil)
			Expect(err).NotTo(HaveOccurred())
			Expect(run).To(BeTrue())
			Expect(sha).To(BeEmpty())
		})
	})

	context("Run", func() {
		it("succeeds", func() {
			Expect(process.Run(modulesDir, cacheDir, workingDir)).To(Succeed())
			Expect(executable.ExecuteCall.Receives.Execution).To(Equal(pexec.Execution{
				Args:   []string{"install", "--unsafe-perm", "--cache", cacheDir},
				Dir:    workingDir,
				Stdout: buffer,
				Stderr: buffer,
				Env:    append(os.Environ(), "NPM_CONFIG_PRODUCTION=true", "NPM_CONFIG_LOGLEVEL=error"),
			}))

			path, err := os.Readlink(filepath.Join(workingDir, "node_modules"))
			Expect(err).NotTo(HaveOccurred())
			Expect(path).To(Equal(filepath.Join(modulesDir, "node_modules")))
		})

		context("failure cases", func() {
			context("when unable to write node_modules directory in layer", func() {
				it.Before(func() {
					Expect(os.Chmod(modulesDir, 0000)).To(Succeed())
				})

				it("fails", func() {
					err := process.Run(modulesDir, cacheDir, workingDir)
					Expect(err).To(MatchError(ContainSubstring("permission denied")))
				})
			})

			context("when the node_modules directory cannot be symlinked into the working directory", func() {
				it.Before(func() {
					Expect(os.Chmod(workingDir, 0000)).To(Succeed())
				})

				it("returns an error", func() {
					err := process.Run(modulesDir, cacheDir, workingDir)
					Expect(err).To(MatchError(ContainSubstring("permission denied")))
				})
			})

			context("when the executable fails", func() {
				it.Before(func() {
					executable.ExecuteCall.Stub = func(execution pexec.Execution) error {
						fmt.Fprintln(execution.Stdout, "install error on stdout")
						fmt.Fprintln(execution.Stderr, "install error on stderr")
						return errors.New("failed to execute")
					}
				})

				it("returns an error", func() {
					err := process.Run(modulesDir, cacheDir, workingDir)
					Expect(buffer.String()).To(ContainSubstring("    install error on stdout\n    install error on stderr\n"))
					Expect(err).To(MatchError("npm install failed: failed to execute"))
				})
			})
		})
	})
}
