/*


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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	hellov1beta1 "example.com/api/v1beta1"
)

// SaiyamReconciler reconciles a Saiyam object
type SaiyamReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=hello.example.com,resources=saiyams,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=hello.example.com,resources=saiyams/status,verbs=get;update;patch

func (r *SaiyamReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("saiyam", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *SaiyamReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hellov1beta1.Saiyam{}).
		Complete(r)
}
