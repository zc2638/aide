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

package app

import (
	"errors"
	"os"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"

	"github.com/zc2638/aide"
)

func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use: "aide",
	}
	cmd.AddCommand(NewApplyCmd())
	return cmd
}

type ApplyOption struct {
	Path string
}

func NewApplyCmd() *cobra.Command {
	opt := &ApplyOption{}
	cmd := &cobra.Command{
		Use:          "apply",
		SilenceUsage: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(opt.Path) == 0 {
				return errors.New("please specify the configuration file to execute")
			}
			filedata, err := os.ReadFile(opt.Path)
			if err != nil {
				return err
			}
			var pipeline aide.Pipeline
			if err := yaml.Unmarshal(filedata, &pipeline); err != nil {
				return err
			}

			pipeline.ParseFlags()
			if err := pipeline.Validate(); err != nil {
				return err
			}
			return pipeline.Execute(cmd.Context())
		},
	}
	cmd.Flags().StringVarP(&opt.Path, "file", "f", opt.Path, "that contains the configuration to apply")
	return cmd
}
