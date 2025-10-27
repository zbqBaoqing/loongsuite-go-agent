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

package instrument

import (
	"path/filepath"

	"github.com/alibaba/loongsuite-go-agent/tool/ast"
	"github.com/alibaba/loongsuite-go-agent/tool/ex"
	"github.com/alibaba/loongsuite-go-agent/tool/rules"
	"github.com/alibaba/loongsuite-go-agent/tool/util"
	"github.com/dave/dst"
)

func (rp *RuleProcessor) addStructField(rule *rules.InstStructRule, decl dst.Decl) {
	util.Assert(rule.FieldName != "" && rule.FieldType != "",
		"rule must have field and type")
	util.Log("Apply struct rule %v (%v)", rule, rp.compileArgs)
	ast.AddStructField(decl, rule.FieldName, rule.FieldType)
}

func (rp *RuleProcessor) applyStructRules(bundle *rules.InstRuleSet) error {
	for file, stRules := range bundle.StructRules {
		util.Assert(filepath.IsAbs(file), "file path must be absolute")
		// Apply struct rules to the file
		astRoot, err := rp.parseAst(file)
		if err != nil {
			return err
		}
		for _, stRule := range stRules {
			structDecl := ast.FindStructDecl(astRoot, stRule.StructType)
			if structDecl != nil {
				rp.addStructField(stRule, structDecl)
			} else {
				return ex.Newf("struct %s not found", stRule.StructType)
			}
		}
		// Once all struct rules are applied, we restore AST to file and use it
		// in future compilation
		newFile, err := rp.writeInstrumented(file, astRoot)
		if err != nil {
			return err
		}
		// Line directive must be placed at the beginning of the line, otherwise
		// it will be ignored by the compiler
		err = rp.enableLineDirective(newFile)
		if err != nil {
			return err
		}
		rp.keepForDebug(newFile)
	}
	return nil
}
