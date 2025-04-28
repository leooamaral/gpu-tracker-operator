/*
Copyright 2025.

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

package controller

import (
	"context"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	gpunodev1 "github.com/leooamaral/gpu-tracker-operator/api/v1"
)

type GPUTrackerReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=suse.tests.dev,resources=gputrackers,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=suse.tests.dev,resources=gputrackers/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=suse.tests.dev,resources=gputrackers/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the GPUTracker object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *GPUTrackerReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	logger.Info("Starting reconciliation for GPUTracker", "name", req.Name)

	tracker := &gpunodev1.GPUTracker{}
	if err := r.Get(ctx, req.NamespacedName, tracker); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	nodeList := &corev1.NodeList{}
	if err := r.List(ctx, nodeList, client.MatchingLabels{"node-type": "gpu-node"}); err != nil {
		logger.Error(err, "Failed to list nodes with label node-type=gpu-node")
		return ctrl.Result{}, err
	}

	var nodeNames []string
	for _, node := range nodeList.Items {
		nodeNames = append(nodeNames, node.Name)
	}

	commaSeparatedNodes := strings.Join(nodeNames, ",")

	if tracker.GPUNodes != commaSeparatedNodes {
		tracker.GPUNodes = commaSeparatedNodes

		if err := r.Update(ctx, tracker); err != nil {
			log.Log.Error(err, "Failed to update GPUTracker node list")
			return ctrl.Result{}, err
		}
		logger.Info("Successfully updated GPUTracker with matching nodes", "nodes", commaSeparatedNodes)
	} else {
		logger.Info("No changes detected in GPUTracker gpu_nodes, skipping update")
	}

	return ctrl.Result{RequeueAfter: 1 * time.Minute}, nil
}

func (r *GPUTrackerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gpunodev1.GPUTracker{}).
		Complete(r)
}
