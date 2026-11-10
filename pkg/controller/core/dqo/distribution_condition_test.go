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

package dqo

import (
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
	apimeta "k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	kueuealpha "sigs.k8s.io/kueue/apis/kueue/v1alpha1"
	"sigs.k8s.io/kueue/pkg/util/api"
	utiltestingalpha "sigs.k8s.io/kueue/pkg/util/testing/v1alpha1"
)

func TestSetDistributionConditionTruncatesMessage(t *testing.T) {
	const generation int64 = 7
	status := metav1.ConditionFalse
	reason := kueuealpha.DynamicQuotaOrchestratorReasonMisconfigured

	cases := map[string]struct {
		message string
	}{
		"short message is unchanged": {
			message: "Quotas successfully distributed",
		},
		"message longer than 32KiB is truncated": {
			message: strings.Repeat("a", 32*1024+128),
		},
	}

	for name, tc := range cases {
		t.Run(name, func(t *testing.T) {
			orchestrator := utiltestingalpha.MakeDynamicQuotaOrchestrator("dqo").
				Generation(generation).
				Obj()
			r := &Reconciler{}
			r.setDistributionCondition(orchestrator, status, reason, tc.message)

			got := apimeta.FindStatusCondition(orchestrator.Status.Conditions, kueuealpha.DynamicQuotaOrchestratorDistributed)
			if got == nil {
				t.Fatal("Distributed condition was not set")
			}
			if got.Status != status || got.Reason != reason || got.ObservedGeneration != generation {
				t.Errorf("unexpected condition metadata: status=%s reason=%s observedGeneration=%d", got.Status, got.Reason, got.ObservedGeneration)
			}
			wantMessage := api.TruncateConditionMessage(tc.message)
			if diff := cmp.Diff(wantMessage, got.Message); diff != "" {
				t.Errorf("Unexpected message (-want +got):\n%s", diff)
			}
		})
	}
}
