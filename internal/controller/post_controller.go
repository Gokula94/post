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
	"flag"
	"log"
	"path/filepath"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

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
	//log := log.FromContext(ctx)
	//rules := clientcmd.NewDefaultClientConfigLoadingRules()
	//kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(rules, &clientcmd.ConfigOverrides{})

	// c, err := kubeconfig.ClientConfig()
	// if err != nil {
	// 	panic(err)
	// }

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	cfg, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Println("could not get config")
	}

	clusterClient, err := dynamic.NewForConfig(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	resource := httpv1alpha1.GroupVersion.WithResource("posts")
	factory := dynamicinformer.NewFilteredDynamicSharedInformerFactory(clusterClient,
		time.Minute, "", nil)
	informer := factory.ForResource(resource).Informer()
	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			log.Println("AddFunc")
		},
		UpdateFunc: func(oldObj, newObj interface{}) {
			log.Println("UpdateFunc")
		},
		DeleteFunc: func(obj interface{}) {
			log.Println("DeleteFunc")
		},
	})

	informer.Run(make(chan struct{}))

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *PostReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&httpv1alpha1.Post{}).
		Complete(r)
}
