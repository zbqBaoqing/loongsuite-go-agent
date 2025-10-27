// Copyright (c) 2024 Alibaba Group Holding Ltd.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package preprocess

import (
	"archive/tar"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/alibaba/loongsuite-go-agent/tool/config"
	"github.com/alibaba/loongsuite-go-agent/tool/data"
	"github.com/alibaba/loongsuite-go-agent/tool/ex"
	"github.com/alibaba/loongsuite-go-agent/tool/rules"
	"github.com/alibaba/loongsuite-go-agent/tool/util"
	"golang.org/x/mod/modfile"
)

const (
	pkgPrefix = "github.com/alibaba/loongsuite-go-agent/pkg"
)

var otelDeps = map[string]string{
	"go.opentelemetry.io/otel":                                          "v1.35.0",
	"go.opentelemetry.io/otel/sdk":                                      "v1.35.0",
	"go.opentelemetry.io/otel/trace":                                    "v1.35.0",
	"go.opentelemetry.io/otel/metric":                                   "v1.35.0",
	"go.opentelemetry.io/otel/sdk/metric":                               "v1.35.0",
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace":                 "v1.35.0",
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc":   "v1.35.0",
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp":   "v1.35.0",
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc": "v1.35.0",
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetrichttp": "v1.35.0",
	"go.opentelemetry.io/otel/exporters/prometheus":                     "v0.57.0",
	"go.opentelemetry.io/contrib/instrumentation/runtime":               "v0.60.0",
	"google.golang.org/protobuf":                                        "v1.35.2",
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric":            "v1.35.0",
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace":             "v1.35.0",
	"go.opentelemetry.io/otel/exporters/zipkin":                         "v1.35.0",
}

func parseGoMod(gomod string) (*modfile.File, error) {
	data, err := util.ReadFile(gomod)
	if err != nil {
		return nil, err
	}
	modFile, err := modfile.Parse(util.GoModFile, []byte(data), nil)
	if err != nil {
		return nil, ex.Wrap(err)
	}
	return modFile, nil
}

func writeGoMod(gomod string, modfile *modfile.File) error {
	bs, err := modfile.Format()
	if err != nil {
		return ex.Wrap(err)
	}
	_, err = util.WriteFile(gomod, string(bs))
	if err != nil {
		return err
	}
	return nil
}

func extractGZip(data []byte, targetDir string) error {
	err := os.MkdirAll(targetDir, 0755)
	if err != nil {
		return ex.Wrap(err)
	}

	gzReader, err := gzip.NewReader(strings.NewReader(string(data)))
	if err != nil {
		return ex.Wrap(err)
	}
	defer func() {
		err := gzReader.Close()
		if err != nil {
			ex.Fatal(err)
		}
	}()

	tarReader := tar.NewReader(gzReader)
	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return ex.Wrap(err)
		}

		// Skip AppleDouble files (._filename) and other hidden files
		if strings.HasPrefix(filepath.Base(header.Name), "._") ||
			strings.HasPrefix(filepath.Base(header.Name), ".") {
			continue
		}

		// Rename pkg_tmp to pkg in the path
		// Normalize path to Unix style for consistent string operations
		cleanName := filepath.ToSlash(filepath.Clean(header.Name))
		if strings.HasPrefix(cleanName, "pkg_tmp/") {
			cleanName = strings.Replace(cleanName, "pkg_tmp/", "pkg/", 1)
		} else if cleanName == "pkg_tmp" {
			cleanName = "pkg"
		}

		// Sanitize the file path to prevent Zip Slip vulnerability
		if cleanName == "." || cleanName == ".." ||
			strings.HasPrefix(cleanName, "..") {
			continue
		}

		// Ensure the resolved path is within the target directory
		targetPath := filepath.Join(targetDir, cleanName)
		resolvedPath, err := filepath.EvalSymlinks(targetPath)
		if err != nil {
			// If symlink evaluation fails, use the original path
			resolvedPath = targetPath
		}

		// Check if the resolved path is within the target directory
		relPath, err := filepath.Rel(targetDir, resolvedPath)
		if err != nil || strings.HasPrefix(relPath, "..") ||
			filepath.IsAbs(relPath) {
			continue // Skip files that would be extracted outside target dir
		}
		switch header.Typeflag {
		case tar.TypeDir:
			err = os.MkdirAll(targetPath, os.FileMode(header.Mode))
			if err != nil {
				return ex.Wrap(err)
			}

		case tar.TypeReg:
			file, err := os.OpenFile(targetPath, os.O_CREATE|os.O_RDWR,
				os.FileMode(header.Mode))
			if err != nil {
				return ex.Wrap(err)
			}

			_, err = io.Copy(file, tarReader)
			if err != nil {
				return ex.Wrap(err)
			}
			err = file.Close()
			if err != nil {
				return ex.Wrap(err)
			}

		default:
			return ex.Newf("unsupported file type: %c in %s",
				header.Typeflag, header.Name)
		}
	}

	return nil
}

// Fetch the zipped pkg module from the embedded data section and extract it to
// a temporary directory, then return the path to the pkg directory.
// In -pkg flag is specified, we will use the path specified by the flag instead
// of the temporary directory. This flag is specifically used for a scenario where
// the user compiles a shared library and the main program separately. If the
// injected code uses different temporary pkg directories, it will cause the
// plugin to report a version conflict error. This flag ensures that the shared
// library and the main program share the same pkg directory to prevent this issue.
func findPkgModDir() (string, error) {
	pkgPath := config.GetConf().PkgPath
	// If -pkg flag is specified, check if the path exists. If it does,
	// return the absolute path. If it does not, extract the embedded pkg
	// module into the path.
	if pkgPath != "" {
		if util.PathExists(pkgPath) {
			return filepath.Join(pkgPath, "pkg"), nil
		}
		bs, err := data.UseEmbeddedPkg()
		if err != nil {
			return "", err
		}
		tempPkg := pkgPath
		err = extractGZip(bs, tempPkg)
		if err != nil {
			return "", err
		}
		return filepath.Join(tempPkg, "pkg"), nil
	}
	bs, err := data.UseEmbeddedPkg()
	if err != nil {
		return "", err
	}
	tempPkg := util.GetTempBuildDirWith("alibaba-pkg")
	if util.PathExists(tempPkg) {
		_ = os.RemoveAll(tempPkg)
	}
	err = extractGZip(bs, tempPkg)
	if err != nil {
		return "", err
	}
	return filepath.Join(tempPkg, "pkg"), nil
}

// updateRule rectifies the file rules path to the local module cache path.
func (dp *DepProcessor) updateRule(bundles []*rules.InstRuleSet) error {
	util.GuaranteeInPreprocess()
	defer util.PhaseTimer("Fetch")()
	modfile, err := parseGoMod(dp.getGoModPath())
	if err != nil {
		return err
	}
	// Collect all the replace directives from go.mod file, we will use it to
	// rectify the custom rules.
	replaceMap := map[string]string{}
	for _, replace := range modfile.Replace {
		replaceMap[replace.Old.Path] = replace.New.Path
	}
	rectified := map[string]bool{}
	for _, bundle := range bundles {
		for _, funcRules := range bundle.FuncRules {
			for _, rule := range funcRules {
				if rule.UseRaw {
					continue
				}
				if rectified[rule.GetPath()] {
					continue
				}
				_, path, err := dp.findRuleDir(rule.GetPath())
				if err != nil {
					return err
				}
				rule.SetPath(path)
				rectified[path] = true
			}
		}
		for _, fileRule := range bundle.FileRules {
			if rectified[fileRule.GetPath()] {
				continue
			}
			_, path, err := dp.findRuleDir(fileRule.GetPath())
			if err != nil {
				return err
			}
			fileRule.SetPath(path)
			fileRule.FileName = filepath.Join(path, fileRule.FileName)
			rectified[path] = true
		}
	}
	return nil
}

func (dp *DepProcessor) updateGoMod() error {
	// Backup go.mod and go.sum files
	gomodDir := dp.getGoModDir()
	files := []string{}
	files = append(files, filepath.Join(gomodDir, util.GoModFile))
	files = append(files, filepath.Join(gomodDir, util.GoSumFile))
	files = append(files, filepath.Join(gomodDir, util.GoWorkSumFile))
	for _, file := range files {
		if util.PathExists(file) {
			err := dp.backupFile(file)
			if err != nil {
				return err
			}
		}
	}
	// Add the alibaba-otel pkg module to the go.mod file
	addDeps := make([]Dependency, 0)
	dep := Dependency{
		ImportPath:     pkgPrefix,
		Version:        "v0.0.0-00010101000000-000000000000",
		Replace:        true,
		ReplacePath:    dp.pkgModDir,
		ReplaceVersion: "",
	}
	addDeps = append(addDeps, dep)
	// OTel dependencies may publish new versions that are not compatible
	// with the otel tool. In such cases, we need to add a replace directive
	// to use certain versions of the OTel dependencies for given otel tool,
	// otherwise, the otel tool may fail to run.
	for path, version := range otelDeps {
		addDeps = append(addDeps, Dependency{
			ImportPath:     path,
			Version:        version,
			Replace:        true,
			ReplacePath:    path,
			ReplaceVersion: version,
		})
	}
	err := dp.addDependency(dp.getGoModPath(), addDeps)
	if err != nil {
		return err
	}
	// Update the existing replace directives to use the local module cache
	// Very bad, we must guarantee the replace path is consistent either in
	// go.mod or vendor/modules.txt, otherwise, the go build toolchain will fail
	// so we must parse go.mod to check if there is any existing replace directive
	// and update the vendor/modules.txt accordingly.
	modfile, err := parseGoMod(dp.getGoModPath())
	if err != nil {
		return err
	}
	changed := false
	for _, replace := range modfile.Replace {
		if replace.Old.Path == pkgPrefix {
			err = modfile.DropReplace(pkgPrefix, "")
			if err != nil {
				return ex.Wrap(err)
			}
			err = modfile.AddReplace(pkgPrefix, "", dp.pkgModDir, "")
			if err != nil {
				return ex.Wrap(err)
			}
			changed = true
			break
		}
	}
	if changed {
		err = writeGoMod(dp.getGoModPath(), modfile)
		if err != nil {
			return err
		}
	}
	return nil
}
