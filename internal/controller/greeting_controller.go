package controller

import (
	"context"
	"fmt"

	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"k8s.io/apimachinery/pkg/runtime/schema"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

// GVKGreeting is the GroupVersionKind for the Greeting resource.
var gvkGreeting = schema.GroupVersionKind{
	Group:   "operator.example.com",
	Version: "v1alpha1",
	Kind:    "Greeting",
}

type GreetingReconciler struct {
	client client.Client
}

// Implement the business logic:
// This function will be called when there is a change to a ReplicaSet or a Pod with an OwnerReference
// to a ReplicaSet.
//
// * Read the ReplicaSet
// * Read the Pods
// * Set a Label on the ReplicaSet with the Pod count.
func (r *GreetingReconciler) Reconcile(ctx context.Context, req reconcile.Request) (reconcile.Result, error) {
	ns := req.NamespacedName
	log := ctrl.LoggerFrom(ctx).WithValues("greeting", ns)

	// Fetch Greeting (unstructured)
	gr := &unstructured.Unstructured{}
	gr.SetGroupVersionKind(gvkGreeting)
	if err := r.client.Get(ctx, ns, gr); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	// Read spec.message
	spec, _ := gr.Object["spec"].(map[string]any)
	msg, _ := spec["message"].(string)

	// Fetch/create the Greeting's ConfigMap
	cmName := fmt.Sprintf("%s-cm", gr.GetName())
	cm := &coreV1.ConfigMap{}
	nsStr := gr.GetNamespace()
	err := r.client.Get(ctx, types.NamespacedName{Namespace: nsStr, Name: cmName}, cm)
	if client.IgnoreNotFound(err) == nil && err != nil {
		return ctrl.Result{}, err
	}
	if err != nil { // NotFound
		cm = &coreV1.ConfigMap{
			ObjectMeta: metaV1.ObjectMeta{
				Name:      cmName,
				Namespace: nsStr,
				Labels:    map[string]string{"app.kubernetes.io/managed-by": "greeting-operator"},
			},
			Data: map[string]string{"message": msg},
		}

		// Set owner reference to the Greeting for garbage collection
		// Note: Using &runtime.Scheme{} here for simplicity; in real code, use the actual scheme
		// passed to the controller.
		if err := controllerutil.SetControllerReference(gr, cm, &runtime.Scheme{}); err != nil {
			return ctrl.Result{}, err
		}
		// Create the ConfigMap
		if err := r.client.Create(ctx, cm); err != nil {
			return ctrl.Result{}, err
		}
		log.Info("created ConfigMap", "name", cmName)
		return ctrl.Result{}, nil
	}

	// Update if changed
	if cm.Data == nil {
		cm.Data = map[string]string{}
	}
	if cm.Data["message"] != msg {
		cm.Data["message"] = msg
		if err := r.client.Update(ctx, cm); err != nil {
			return ctrl.Result{}, err
		}
		log.Info("updated ConfigMap", "name", cmName)
	}

	return ctrl.Result{}, nil
}

func (r *GreetingReconciler) SetupWithManager(mgr ctrl.Manager) error {
	u := &unstructured.Unstructured{}
	u.SetGroupVersionKind(gvkGreeting)
	return ctrl.NewControllerManagedBy(mgr).
		For(u, builder.WithPredicates()).
		Owns(&coreV1.ConfigMap{}).
		Complete(r)
}

func SetupGreeting(mgr ctrl.Manager) error {
	return (&GreetingReconciler{client: mgr.GetClient()}).SetupWithManager(mgr)
}
