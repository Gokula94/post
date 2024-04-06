/*
Copyright 2024 Gokul.

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
	"os"
	"path/filepath"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	//"sigs.k8s.io/controller-runtime/pkg/log"

	httpv1alpha1 "post.com/api/v1alpha1"
)

// PostReconciler reconciles a Post object
type PostReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=http.gokula.zinkworks,resources=posts,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=http.gokula.zinkworks,resources=posts/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=http.gokula.zinkworks,resources=posts/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Post object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.16.3/pkg/reconcile
func (r *PostReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	//post := &httpv1alpha1.Post{}

	//err := r.Get(ctx, req.NamespacedName, post)
	kubeconfig := filepath.Join(homedir.HomeDir(), ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		fmt.Println("Error building kubeconfig:", err)
		os.Exit(1)
	}

	// Create a Kubernetes clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		fmt.Println("Error creating clientset:", err)
		os.Exit(1)
	}
	logger.Info("gonna start restclient")
	group := "http"
	version := "v1alpha1"
	namespace := "default" // Specify the namespace where your custom resource is located
	plural := "posts"      // Specify the plural name of your custom resource

	// Get the absolute path for accessing your custom resource
	absPath := clientset.RESTClient().Get().AbsPath(fmt.Sprintf("/apis/%s/%s/namespaces/%s/%s", group, version, namespace, plural))

	// Print the absolute path
	fmt.Println("Absolute path for your custom resource:", absPath)

	return ctrl.Result{}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *PostReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&httpv1alpha1.Post{}).
		Complete(r)
}
