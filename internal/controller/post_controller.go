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
	"io"
	"net/http"
	"strings"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

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

	//logger.Info("customresource", req.NamespacedName)

	// Fetch the CustomResource instance
	//instance := &httpv1alpha1.Post{}
	// if err := r.Get(ctx, req.NamespacedName, instance); err != nil {
	// 	return ctrl.Result{}, client.IgnoreNotFound(err)
	// }

	// Call REST API
	//myurl := "https://api.restful-api.dev/objects"
	var requeue bool
	post := &httpv1alpha1.Post{}
	logger.Info("evaluating custom resource")

	err := r.Get(ctx, req.NamespacedName, post)
	if err != nil {
		logger.Info("Error getting resource: %v", err)
		requeue = true

		if requeue {
			logger.Info(fmt.Sprintf("Pod created is %v", req.NamespacedName))

			const myurl = "https://api.restful-api.dev/objects"
			fmt.Println(myurl)
			requestBody := strings.NewReader(`
			 
				{
					"name": "3801-XGS-PON",
					"data": {
						"Rx_Operating_Wavelength": 1280,
						"type": "10G Passive  Optical Network (PON) transceivers"
					}
				 }

			 `)
			fmt.Println(requestBody)
			response, err := http.Post(myurl, "application/json", requestBody)
			if err != nil {
				panic(err)
			}
			fmt.Println("post call is suceesful")
			defer response.Body.Close()
			content, _ := io.ReadAll(response.Body)

			fmt.Println(string(content))
		}
	}

	// //req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(payload))
	// req, err := http.Post(myurl, "application/json", payload)
	// if err != nil {
	// 	logger.Info("Failed to create HTTP request")
	// 	return ctrl.Result{}, err
	// }
	// req.Header.Set("Content-Type", "application/json")

	// client := &http.Client{}
	// resp, err := client.Do(req)
	// if err != nil {
	// 	logger.Info("Failed to call REST API")
	// 	return ctrl.Result{}, err
	// }
	// defer resp.Body.Close()
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PostReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&httpv1alpha1.Post{}).
		Complete(r)
}
