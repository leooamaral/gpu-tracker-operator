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
	"fmt"
	"strings"
	"time"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	gpunodev1 "github.com/leooamaral/gpu-tracker-operator/api/v1"
)

// GPUTrackerReconciler reconciles a GPUTracker object
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
	// logger := log.FromContext(ctx)
	// logger.V(consts.LogLevelInfo).Info("Reconciling GPUTracker")

	tracker := &gpunodev1.GPUTracker{}
	if err := r.Get(ctx, req.NamespacedName, tracker); err != nil {
		err = fmt.Errorf("error getting GPUTracker object: %w", err)
		// logger.V(consts.LogLevelError).Error(nil, err.Error())
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Check if GPU Node object exist

	// Check nodes and its labels - node-type: gpu-node
	nodeList := &corev1.NodeList{}
	if err := r.List(ctx, nodeList, client.MatchingLabels{"node-type": "gpu-node"}); err != nil {
		return ctrl.Result{}, err
	}

	var nodeNames []string
	for _, node := range nodeList.Items {
		nodeNames = append(nodeNames, node.Name)
	}

	commaSeparatedNodes := strings.Join(nodeNames, ",")

	// check if node list is the same as the present on CRD
	if tracker.GPUNodes != commaSeparatedNodes {
		tracker.GPUNodes = commaSeparatedNodes

		if err := r.Update(ctx, tracker); err != nil {
			log.Log.Error(err, "Failed to update GPUTracker node list")
			return ctrl.Result{}, err
		}
		if err := r.Status().Update(ctx, tracker); err != nil {
			log.Log.Error(err, "Failed to update GPUTracker status")
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{RequeueAfter: 5 * time.Minute}, nil
}

func checkTrackerList(nodeListA, nodeListB []string) bool {
	if len(nodeListA) != len(nodeListB) {
		return false
	}

	for e := range nodeListA {
		if nodeListA[e] != nodeListB[e] {
			return false
		}
	}
	return true
}

// SetupWithManager sets up the controller with the Manager.
func (r *GPUTrackerReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gpunodev1.GPUTracker{}).
		Complete(r)
}
