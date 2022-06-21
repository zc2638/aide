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
	"fmt"
	"html/template"
	"regexp"
	"sort"
)

var tpl = template.New("aide")

const maxCharLength = 63

var NameRegexp = regexp.MustCompile(`^[a-z0-9]([_a-z0-9]?[a-z0-9])*([a-z0-9]([-a-z0-9]?[a-z0-9]+)*)*$`)

func ValidateName(name string) error {
	if len(name) > maxCharLength {
		return fmt.Errorf("name is too long and cannot exceed %d characters", maxCharLength)
	}
	if !NameRegexp.MatchString(name) {
		return fmt.Errorf("name is invalid, must match regexp %q", NameRegexp.String())
	}
	return nil
}

func envToSlice(set map[string]string) []string {
	env := make([]string, 0, len(set))
	for k, v := range set {
		env = append(env, fmt.Sprintf("%s=%s", k, v))
	}
	sort.Strings(env)
	return env
}
