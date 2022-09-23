// Copyright 2022 The EasyDL Authors. All rights reserved.
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package controllers

import (
	elasticv1alpha1 "github.com/intelligent-machine-learning/easydl/operator/api/v1alpha1"
	"github.com/stretchr/testify/assert"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

func TestGeneratePod(t *testing.T) {
	job := &elasticv1alpha1.ElasticJob{
		ObjectMeta: metav1.ObjectMeta{
			Name:        "test-job",
			Namespace:   "easydl",
			Annotations: map[string]string{},
			Labels:      map[string]string{},
		},
	}
	container := corev1.Container{
		Name:            "main",
		Image:           "test",
		ImagePullPolicy: corev1.PullAlways,
		Command:         []string{"python", "--version"},
	}
	podTemplate := &corev1.PodTemplateSpec{
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{container},
		},
	}

	manager := newPodManager()
	pod := manager.GeneratePod(job, podTemplate, "test-job-worker-0")
	assert.Equal(t, pod.Name, "test-job-worker-0")
	assert.Equal(t, pod.Spec.Containers[0].Image, "test")
}

func TestGetReplicaStatus(t *testing.T) {
	pods := []corev1.Pod{}
	pod0 := corev1.Pod{}
	pod0.Status.Phase = corev1.PodRunning
	pods = append(pods, pod0)
	pod1 := corev1.Pod{}
	pod1.Status.Phase = corev1.PodFailed
	pods = append(pods, pod1)
	pod2 := corev1.Pod{}
	pod2.Status.Phase = corev1.PodSucceeded
	pods = append(pods, pod2)
	manager := newPodManager()
	replicaStatus := manager.GetReplicaStatus(pods)
	int32One := int32(1)
	assert.Equal(t, replicaStatus.Active, int32One)
	assert.Equal(t, replicaStatus.Failed, int32One)
	assert.Equal(t, replicaStatus.Succeeded, int32One)
}