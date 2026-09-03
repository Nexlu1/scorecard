// Copyright 2026 OpenSSF Scorecard Authors
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

package raw

import (
	"reflect"
	"testing"

	"github.com/ossf/scorecard/v5/checker"
)

func TestGitHubWorkflowRunLineOffsets(t *testing.T) {
	t.Parallel()
	content := []byte(`on: push
jobs:
  jobOne:
    runs-on: ubuntu-latest
    steps:
      - run: |
          curl https://example.com | bash
      - run: curl https://example.com | bash
`)

	var result checker.PinningDependenciesData
	_, err := validateGitHubWorkflowIsFreeOfInsecureDownloads(
		".github/workflows/test.yaml", content, &result,
	)
	if err != nil {
		t.Fatalf("validateGitHubWorkflowIsFreeOfInsecureDownloads: %v", err)
	}

	var got []uint
	for _, dependency := range result.Dependencies {
		if dependency.Type != checker.DependencyUseTypeDownloadThenRun {
			continue
		}
		got = append(got, dependency.Location.Offset)
	}
	want := []uint{7, 8}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("downloadThenRun offsets = %v, want %v", got, want)
	}
}
