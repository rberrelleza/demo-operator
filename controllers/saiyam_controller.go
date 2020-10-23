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
	"fmt"
	"time"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
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
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete

func (r *SaiyamReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("saiyam", req.NamespacedName)

	// your logic here
	log.Info("managing resource")

	resource := &hellov1beta1.Saiyam{}
	err := r.Get(ctx, req.NamespacedName, resource)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			log.Info("resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}

		return ctrl.Result{Requeue: false}, err
	}

	found := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: resource.Name, Namespace: resource.Namespace}, found)
	if err == nil {
		log.Info("resource already exists")
		return ctrl.Result{Requeue: true, RequeueAfter: 10 * time.Second}, nil
	}

	log.Info("creating deployment")
	log.Info(fmt.Sprintf("container image %s", resource.Spec.Image))
	d := r.deploymentForImage(resource)
	log.Info("Creating a new Deployment", "Deployment.Namespace", d.Namespace, "Deployment.Name", d.Name)
	err = r.Create(ctx, d)
	if err != nil {
		log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", d.Namespace, "Deployment.Name", d.Name)
		return ctrl.Result{}, err
	}

	resource.Status.DeploymentName = d.Name
	if err := r.Status().Update(ctx, resource); err != nil {
		log.Error(err, "Failed to update status")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *SaiyamReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&hellov1beta1.Saiyam{}).
		Complete(r)
}

func (r *SaiyamReconciler) deploymentForImage(s *hellov1beta1.Saiyam) *appsv1.Deployment {
	ls := map[string]string{"app": "saiyam", "saiyam_cr": s.Name}
	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      s.Name,
			Namespace: s.Namespace,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: int32Ptr(1),
			Selector: &metav1.LabelSelector{
				MatchLabels: ls,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: ls,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: s.Spec.Image,
						Name:  "container",
						Ports: []corev1.ContainerPort{{
							ContainerPort: 8080,
						}},
					}},
				},
			},
		},
	}
	// Set instance as the owner and controller
	ctrl.SetControllerReference(s, dep, r.Scheme)
	return dep
}

func int32Ptr(i int32) *int32 { return &i }
