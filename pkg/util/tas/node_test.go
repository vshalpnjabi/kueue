/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package tas

import (
	"testing"

	corev1 "k8s.io/api/core/v1"

	testingnode "sigs.k8s.io/kueue/pkg/util/testingjobs/node"
)

func TestNodeHostname(t *testing.T) {
	cases := map[string]struct {
		node *corev1.Node
		want string
	}{
		"hostname label differs from the Node name": {
			node: testingnode.MakeNode("node-x1").Label(corev1.LabelHostname, "x1").Obj(),
			want: "x1",
		},
		"hostname label equals the Node name": {
			node: testingnode.MakeNode("x1").Label(corev1.LabelHostname, "x1").Obj(),
			want: "x1",
		},
		"hostname label missing falls back to the Node name": {
			node: testingnode.MakeNode("x1").Obj(),
			want: "x1",
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := NodeHostname(tc.node); got != tc.want {
				t.Errorf("NodeHostname() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestNodeMatchesFlavor(t *testing.T) {
	const labelKey = "example.com/pool"
	cases := map[string]struct {
		nodeLabels     map[string]string
		requiredLabels map[string]string
		requiredLevels []string
		want           bool
	}{
		"matching non-empty value": {
			nodeLabels:     map[string]string{labelKey: "gpu"},
			requiredLabels: map[string]string{labelKey: "gpu"},
			want:           true,
		},
		"wrong value": {
			nodeLabels:     map[string]string{labelKey: "cpu"},
			requiredLabels: map[string]string{labelKey: "gpu"},
			want:           false,
		},
		"missing required non-empty label": {
			nodeLabels:     map[string]string{"kubernetes.io/hostname": "x1"},
			requiredLabels: map[string]string{labelKey: "gpu"},
			want:           false,
		},
		"empty required value matches present empty label": {
			nodeLabels:     map[string]string{labelKey: ""},
			requiredLabels: map[string]string{labelKey: ""},
			want:           true,
		},
		"empty required value does not match missing label": {
			nodeLabels:     map[string]string{"kubernetes.io/hostname": "x1"},
			requiredLabels: map[string]string{labelKey: ""},
			want:           false,
		},
		"empty required value does not match non-empty label": {
			nodeLabels:     map[string]string{labelKey: "gpu"},
			requiredLabels: map[string]string{labelKey: ""},
			want:           false,
		},
	}
	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			if got := NodeMatchesFlavor(tc.nodeLabels, tc.requiredLabels, tc.requiredLevels); got != tc.want {
				t.Errorf("NodeMatchesFlavor() = %t, want %t", got, tc.want)
			}
		})
	}
}
