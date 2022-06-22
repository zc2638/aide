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

package aide

import (
	"bytes"
	"context"
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/99nil/gopkg/sets"
	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
)

const APIVersion = "v1"

const PipelineKind = "Pipeline"

type Metadata struct {
	Name   string            `json:"name" yaml:"name"`
	Labels map[string]string `json:"labels" yaml:"labels"`
}

type Spec struct {
	Prompts []SpecPrompt `json:"prompts" yaml:"prompts"`
	Steps   []SpecStep   `json:"steps" yaml:"steps"`
}

type PromptType string

const (
	PromptInput       PromptType = "Input"
	PromptText        PromptType = "Text"
	PromptPassword    PromptType = "Password"
	PromptConfirm     PromptType = "Confirm"
	PromptSelect      PromptType = "Select"
	PromptMultiSelect PromptType = "MultiSelect"
)

type SpecPrompt struct {
	Name    string     `json:"name" yaml:"name"`
	Type    PromptType `json:"type" yaml:"type"`
	Message string     `json:"message" yaml:"message"`
	Enum    []string   `json:"enum" yaml:"enum"`
	Default string     `json:"default" yaml:"default"`
	Help    string     `json:"help" yaml:"help"`
}

type SpecStep struct {
	Name    string          `json:"name" yaml:"name"`
	Render  *SpecStepRender `json:"render" yaml:"render"`
	Command *string         `json:"command" yaml:"command"`
}

type SpecStepRender struct {
	fsys fs.FS

	Src  string `json:"src" yaml:"src"`
	Dest string `json:"dest" yaml:"dest"`
}

func NewStepRender(src, dest string) *SpecStepRender {
	return NewEmbedStepRender(nil, src, dest)
}

func NewEmbedStepRender(fsys fs.FS, src, dest string) *SpecStepRender {
	return &SpecStepRender{fsys: fsys, Src: src, Dest: dest}
}

func NewPipeline(name string) *Pipeline {
	return &Pipeline{
		APIVersion: APIVersion,
		Kind:       PipelineKind,
		Metadata: Metadata{
			Name: name,
		},
	}
}

type Pipeline struct {
	skipPrompt bool

	APIVersion string   `json:"apiVersion" yaml:"apiVersion"`
	Kind       string   `json:"kind" yaml:"kind"`
	Metadata   Metadata `json:"metadata" yaml:"metadata"`
	Spec       Spec     `json:"spec" yaml:"spec"`
}

func (p *Pipeline) addPrompt(typ PromptType, name, message string, enum []string, defVal, help string) {
	prompt := SpecPrompt{
		Name:    name,
		Type:    typ,
		Message: message,
		Default: defVal,
		Help:    help,
	}
	if len(enum) > 0 {
		prompt.Enum = enum
	}
	p.Spec.Prompts = append(p.Spec.Prompts, prompt)
}

func (p *Pipeline) AddInputPrompt(name, message string, defVal, helpVal string) {
	p.addPrompt(PromptInput, name, message, nil, defVal, helpVal)
}

func (p *Pipeline) AddPasswordPrompt(name, message, helpVal string) {
	p.addPrompt(PromptPassword, name, message, nil, "", "")
}

func (p *Pipeline) AddTextPrompt(name, message string, defVal, helpVal string) {
	p.addPrompt(PromptText, name, message, nil, defVal, helpVal)
}

func (p *Pipeline) AddConfirmPrompt(name, message string, defVal bool, help string) {
	boolean := strconv.FormatBool(defVal)
	p.addPrompt(PromptConfirm, name, message, nil, boolean, help)
}

func (p *Pipeline) AddSelectPrompt(name, message string, enum []string, defVal, helpVal string) {
	p.addPrompt(PromptSelect, name, message, enum, defVal, helpVal)
}

func (p *Pipeline) AddMultiSelectPrompt(name, message string, enum []string, defVal []string, helpVal string) {
	p.addPrompt(PromptMultiSelect, name, message, enum, strings.Join(defVal, ","), helpVal)
}

func (p *Pipeline) AddStep(name string, render *SpecStepRender, command string) {
	step := SpecStep{Name: name}
	if render != nil {
		step.Render = render
	}
	if len(command) > 0 {
		step.Command = &command
	}
	p.Spec.Steps = append(p.Spec.Steps, step)
}

func (p *Pipeline) BindFlags(set *flag.FlagSet) {
	set.BoolVar(&p.skipPrompt, "skip-prompt", false, "Used to skip prompt interactions")
	for k, prompt := range p.Spec.Prompts {
		set.StringVar(&p.Spec.Prompts[k].Default, prompt.Name, prompt.Default, prompt.Help)
	}
}

func (p *Pipeline) ParseFlags() {
	p.BindFlags(flag.CommandLine)
	flag.Parse()
}

func (p *Pipeline) Validate() error {
	if p.Kind != PipelineKind {
		return fmt.Errorf("kind expect Pipeline, not %s", p.Kind)
	}
	if err := ValidateName(p.Metadata.Name); err != nil {
		return fmt.Errorf("metadata.name validate failed: %v", err)
	}
	if len(p.Spec.Steps) == 0 {
		return errors.New("step is not define")
	}
	for k, prompt := range p.Spec.Prompts {
		if err := ValidateName(prompt.Name); err != nil {
			return fmt.Errorf("prompt[%d].Name validate failed: %v", k, err)
		}
		switch prompt.Type {
		case PromptInput:
		case PromptPassword:
		case PromptText:
		case PromptConfirm:
		case PromptSelect, PromptMultiSelect:
			if len(prompt.Enum) == 0 {
				return fmt.Errorf("prompt[%d].Enum validate failed: select requires at least one enum", k)
			}
		default:
			return fmt.Errorf("prompt[%d].Type validate failed: unknow prompt type(%s)", k, prompt.Type)
		}
	}

	for k, step := range p.Spec.Steps {
		if step.Render == nil && step.Command == nil {
			return errors.New("step render or command must be defined")
		}
		if step.Render != nil || len(step.Name) > 0 {
			if err := ValidateName(step.Name); err != nil {
				return fmt.Errorf("step[%d].Name validate failed: %v", k, err)
			}
		}
	}
	return nil
}

func (p *Pipeline) Execute(ctx context.Context) error {
	envSet := make(map[string]string)

	originEnv := os.Environ()
	for _, env := range originEnv {
		parts := strings.Split(env, "=")
		envSet[parts[0]] = envSet[parts[1]]
	}

	for k, v := range p.Metadata.Labels {
		envSet[k] = v
	}
	if err := p.executePrompts(ctx, envSet); err != nil {
		return err
	}
	return p.executeSteps(ctx, envSet)
}

func (p *Pipeline) executePrompts(_ context.Context, envSet map[string]string) error {
	answers := make(map[string]interface{})
	questions := make([]*survey.Question, 0, len(p.Spec.Prompts))
	for _, v := range p.Spec.Prompts {
		var prompt survey.Prompt
		switch v.Type {
		case PromptInput:
			prompt = &survey.Input{Message: v.Message, Default: v.Default, Help: v.Help}
		case PromptPassword:
			prompt = &survey.Password{Message: v.Message, Help: v.Help}
		case PromptText:
			prompt = &survey.Multiline{Message: v.Message, Default: v.Default, Help: v.Help}
		case PromptConfirm:
			boolean, _ := strconv.ParseBool(v.Default)
			prompt = &survey.Confirm{Message: v.Message, Default: boolean, Help: v.Help}
		case PromptSelect:
			if v.Default == "" && len(v.Enum) > 0 {
				v.Default = v.Enum[0]
			}
			prompt = &survey.Select{
				Message: v.Message,
				Options: v.Enum,
				Default: v.Default,
				Help:    v.Help,
			}
		case PromptMultiSelect:
			parts := strings.Split(v.Default, ",")
			defSet := sets.NewString(parts...)
			defSet = sets.NewString(v.Enum...).Intersection(defSet)
			prompt = &survey.MultiSelect{
				Message: v.Message,
				Options: v.Enum,
				Default: defSet.List(),
				Help:    v.Help,
			}
		default:
			continue
		}
		envSet[v.Name] = v.Default
		questions = append(questions, &survey.Question{
			Name:   v.Name,
			Prompt: prompt,
		})
	}
	if p.skipPrompt {
		return nil
	}
	if err := survey.Ask(questions, &answers); err != nil {
		return err
	}

	for k, answer := range answers {
		switch val := answer.(type) {
		case bool:
			envSet[k] = strconv.FormatBool(val)
		case string:
			envSet[k] = val
		case core.OptionAnswer:
			envSet[k] = val.Value
		case []core.OptionAnswer:
			parts := make([]string, 0, len(val))
			for _, v := range val {
				parts = append(parts, v.Value)
			}
			envSet[k] = strings.Join(parts, ",")
		}
	}
	return nil
}

func (p *Pipeline) executeSteps(ctx context.Context, envSet map[string]string) error {
	var err error
	for k, step := range p.Spec.Steps {
		if step.Render != nil {
			envSet[step.Name+"_src"] = step.Render.Src
			envSet[step.Name+"_dest"] = step.Render.Src

			var stat fs.FileInfo
			if step.Render.fsys != nil {
				stat, err = fs.Stat(step.Render.fsys, step.Render.Src)
				if err != nil {
					return fmt.Errorf("stat embed render[%d] src failed: %v", k, err)
				}
			} else {
				stat, err = os.Stat(step.Render.Src)
				if err != nil {
					return fmt.Errorf("stat render[%d] src failed: %v", k, err)
				}
			}
			if !stat.IsDir() {
				err = p.renderFile(envSet, step.Render.fsys, step.Render.Src, step.Render.Dest)
			} else {
				err = p.renderDir(envSet, step.Render.fsys, step.Render.Src, step.Render.Dest)
			}
			if err != nil {
				return fmt.Errorf("render[%d] failed: %v", k, err)
			}
		}
		if step.Command != nil {
			cmd := exec.CommandContext(ctx, "/bin/sh", "-c", *step.Command)
			cmd.Env = envToSlice(envSet)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				return fmt.Errorf("run command (%s) failed: %v", *step.Command, err)
			}
		}
	}
	return nil
}

func (p *Pipeline) renderDir(envSet map[string]string, fsys fs.FS, src, dest string) error {
	err := os.MkdirAll(dest, fs.ModePerm)
	if err != nil {
		return err
	}

	var dir []fs.DirEntry
	if fsys != nil {
		dir, err = fs.ReadDir(fsys, src)
	} else {
		dir, err = os.ReadDir(src)
	}
	if err != nil {
		return err
	}

	for _, e := range dir {
		if !e.IsDir() {
			if err := p.renderFile(envSet, fsys, src, dest); err != nil {
				return err
			}
			continue
		}

		currentSrc := filepath.Join(src, e.Name())
		currentDest := filepath.Join(dest, e.Name())
		if err := p.renderDir(envSet, fsys, currentSrc, currentDest); err != nil {
			return err
		}
	}
	return nil
}

func (p *Pipeline) renderFile(envSet map[string]string, fsys fs.FS, src, dest string) error {
	dir := filepath.Dir(dest)
	_, err := os.Stat(dir)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
		if err := os.MkdirAll(dir, fs.ModePerm); err != nil {
			return err
		}
	}

	var b []byte
	if fsys != nil {
		b, err = fs.ReadFile(fsys, src)
	} else {
		b, err = os.ReadFile(src)
	}
	if err != nil {
		return err
	}

	t, err := tpl.Parse(string(b))
	if err != nil {
		return fmt.Errorf("parse template(%s) failed: %v", src, err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, map[string]interface{}{"env": envSet}); err != nil {
		return err
	}
	return os.WriteFile(dest, buf.Bytes(), fs.ModePerm)
}
