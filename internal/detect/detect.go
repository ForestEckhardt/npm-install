package detect

import (
	"fmt"
	"path/filepath"

	libbuildpackV3 "github.com/buildpack/libbuildpack"
	"github.com/cloudfoundry/libbuildpack"
	"github.com/cloudfoundry/npm-cnb/internal/build"
	"github.com/cloudfoundry/npm-cnb/package_json"
)

func UpdateBuildPlan(detector *libbuildpackV3.Detect) error {
	packageJSONPath := filepath.Join(detector.Application.Root, "package.json")
	if exists, err := libbuildpack.FileExists(packageJSONPath); err != nil {
		return fmt.Errorf("error checking filepath %s", packageJSONPath)
	} else if !exists {
		return fmt.Errorf("no package.json found in %s", packageJSONPath)
	}

	pkgJSON, err := package_json.LoadPackageJSON(packageJSONPath, detector.Logger)
	if err != nil {
		return err
	}

	detector.BuildPlan[build.NodeDependency] = libbuildpackV3.BuildPlanDependency{
		Version: pkgJSON.Engines.Node,
		Metadata: libbuildpackV3.BuildPlanDependencyMetadata{
			"build":true,
			"launch":true,
		},
	}

	return nil
}
