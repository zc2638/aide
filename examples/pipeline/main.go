// Copyright © 2022 zc2638 <zc2638@qq.com>.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"embed"
	"log"

	"github.com/zc2638/aide"
)

//go:embed resource
var embedFS embed.FS

func main() {
	pipeline := aide.NewPipeline("test")

	pipeline.AddConfirmPrompt("confirm1", "确认", false, "")
	pipeline.AddInputPrompt("input1", "输入", "123", "")
	pipeline.AddPasswordPrompt("pwd1", "密码", "")
	pipeline.AddTextPrompt("text1", "文本", "text123", "")
	pipeline.AddSelectPrompt("select1", "选择", []string{"a", "b", "c"}, "c", "")
	pipeline.AddMultiSelectPrompt("multi_select1", "多选", []string{"a", "b", "c"}, []string{"a", "c"}, "")

	pipeline.AddStep("step1", nil, "env")
	pipeline.AddStep("step2", aide.NewEmbedStepRender(embedFS, "resource/test.in", "testdata2/test.out"), "")
	pipeline.AddStep("step3", nil, "echo $text1")
	if err := pipeline.Validate(); err != nil {
		log.Fatal(err)
	}
	pipeline.ParseFlags()
	if err := pipeline.Execute(context.Background()); err != nil {
		log.Fatal(err)
	}
}
