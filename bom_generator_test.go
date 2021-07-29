package npminstall_test

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	npminstall "github.com/paketo-buildpacks/npm-install"
	"github.com/paketo-buildpacks/npm-install/fakes"
	"github.com/paketo-buildpacks/packit/pexec"
	"github.com/paketo-buildpacks/packit/scribe"
	"github.com/sclevine/spec"

	. "github.com/onsi/gomega"
)

func testBOMGenerator(t *testing.T, context spec.G, it spec.S) {
	var (
		Expect = NewWithT(t).Expect

		workingDir    string
		executable    *fakes.Executable
		executable2   *fakes.Executable
		buffer        *bytes.Buffer
		commandOutput *bytes.Buffer

		bomGenerator npminstall.BOMGenerator
	)

	it.Before(func() {
		var err error
		workingDir, err = ioutil.TempDir("", "working-dir")
		Expect(err).NotTo(HaveOccurred())

		executable = &fakes.Executable{}
		executable2 = &fakes.Executable{}
		executable2.ExecuteCall.Stub = func(execution pexec.Execution) error {
			Expect(os.WriteFile(filepath.Join(workingDir, "bom.json"), []byte{}, 0644)).To(Succeed())
			return nil
		}

		buffer = bytes.NewBuffer(nil)
		commandOutput = bytes.NewBuffer(nil)

		bomGenerator = npminstall.NewBOMGenerator(executable, executable2, scribe.NewEmitter(buffer))
	})

	it.After(func() {
		Expect(os.RemoveAll(workingDir)).To(Succeed())
	})

	context("InstallAndRun", func() {
		it.Before(func() {

		})
		it("succeeds in installing the BOM generation tool", func() {
			bomPath, err := bomGenerator.InstallAndRun(workingDir)
			Expect(err).ToNot(HaveOccurred())

			Expect(executable.ExecuteCall.Receives.Execution).To(Equal(pexec.Execution{
				Args:   []string{"install", "-g", "@cyclonedx/bom"},
				Stdout: commandOutput,
				Stderr: commandOutput,
			}))

			Expect(buffer.String()).To(ContainSubstring("Successful install of cyclonedx/bom"))

			Expect(executable2.ExecuteCall.Receives.Execution).To(Equal(pexec.Execution{
				Args:   []string{"-o", "bom.json"},
				Dir:    workingDir,
				Stdout: commandOutput,
				Stderr: commandOutput,
			}))

			Expect(bomPath).To(Equal(filepath.Join(workingDir, "bom.json")))
			fileInfo, err := os.Stat(bomPath)
			Expect(err).NotTo(HaveOccurred())
			Expect(fileInfo.IsDir()).To(BeFalse())
		})

		context("failure cases", func() {

			context("installing tool fails", func() {
				it.Before(func() {
					executable.ExecuteCall.Returns.Err = errors.New("error with 'cyclonedx/bom'")
				})
				it("returns an error", func() {
					_, err := bomGenerator.InstallAndRun(workingDir)
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(ContainSubstring("failed to install: error with 'cyclonedx/bom'")))
				})
			})

			context("running tool fails", func() {
				it.Before(func() {
					executable2.ExecuteCall.Stub = func(execution pexec.Execution) error {
						Expect(os.WriteFile(filepath.Join(workingDir, "bom.json"), []byte{}, 0644)).To(Succeed())
						return errors.New("error running 'cyclonedx-bom'")
					}
				})

				it("returns an error", func() {
					_, err := bomGenerator.InstallAndRun(workingDir)
					Expect(err).To(HaveOccurred())
					Expect(err).To(MatchError(ContainSubstring("failed to run: error running 'cyclonedx-bom'")))
				})
			})
		})
	})

}
