/*
Copyright 2018 The Knative Authors

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

package v1alpha1

import (
	"testing"

	"github.com/google/go-cmp/cmp"
	duckv1alpha1 "github.com/knative/pkg/apis/duck/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

func TestClusterProvisionerStatusIsReady(t *testing.T) {
	tests := []struct {
		name string
		ps   *ClusterProvisionerStatus
		want bool
	}{{
		name: "uninitialized",
		ps:   &ClusterProvisionerStatus{},
		want: false,
	}, {
		name: "initialized",
		ps: func() *ClusterProvisionerStatus {
			ps := &ClusterProvisionerStatus{}
			ps.InitializeConditions()
			return ps
		}(),
		want: false,
	}, {
		name: "ready true condition",
		ps: &ClusterProvisionerStatus{
			Conditions: []duckv1alpha1.Condition{{
				Type:   ChannelConditionReady,
				Status: corev1.ConditionTrue,
			}},
		},
		want: true,
	}, {
		name: "ready false condition",
		ps: &ClusterProvisionerStatus{
			Conditions: []duckv1alpha1.Condition{{
				Type:   ChannelConditionReady,
				Status: corev1.ConditionFalse,
			}},
		},
		want: false,
	}, {
		name: "unknown condition",
		ps: &ClusterProvisionerStatus{
			Conditions: []duckv1alpha1.Condition{{
				Type:   "foo",
				Status: corev1.ConditionTrue,
			}},
		},
		want: false,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.ps.IsReady()
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("%s: unexpected condition (-want, +got) = %v", test.name, diff)
			}
		})
	}
}

func TestClusterProvisionerStatusGetCondition(t *testing.T) {
	tests := []struct {
		name      string
		ps        *ClusterProvisionerStatus
		condQuery duckv1alpha1.ConditionType
		want      *duckv1alpha1.Condition
	}{{
		name: "single condition",
		ps: &ClusterProvisionerStatus{
			Conditions: []duckv1alpha1.Condition{
				condReady,
			},
		},
		condQuery: duckv1alpha1.ConditionReady,
		want:      &condReady,
	}, {
		name: "unknown condition",
		ps: &ClusterProvisionerStatus{
			Conditions: []duckv1alpha1.Condition{
				condReady,
				condUnprovisioned,
			},
		},
		condQuery: duckv1alpha1.ConditionType("foo"),
		want:      nil,
	}}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := test.ps.GetCondition(test.condQuery)
			if diff := cmp.Diff(test.want, got); diff != "" {
				t.Errorf("unexpected condition (-want, +got) = %v", diff)
			}
		})
	}
}
