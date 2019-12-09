package appservice

import (
	"os"

	"context"
	//"fmt"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	//metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	//"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	certv1alpha1 "github.com/jetstack/cert-manager/pkg/apis/certmanager/v1alpha1"
	routev1 "github.com/openshift/api/route/v1"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	logf "sigs.k8s.io/controller-runtime/pkg/runtime/log"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_appservice")

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new AppService Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	// Adding the routev1
	if err := routev1.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}
	// Adding certv1alpha1
	if err := certv1alpha1.AddToScheme(mgr.GetScheme()); err != nil {
		log.Error(err, "")
		os.Exit(1)
	}

	return &ReconcileAppService{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("appservice-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource Route
	err = c.Watch(&source.Kind{Type: &routev1.Route{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	err = c.Watch(&source.Kind{Type: &corev1.Secret{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &corev1.Secret{},
	})
	if err != nil {
		return err
	}



	return nil
}

// blank assignment to verify that ReconcileAppService implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileAppService{}

// ReconcileAppService reconciles a AppService object
type ReconcileAppService struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a AppService object and makes changes based on the state read
// and what is in the AppService.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileAppService) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling AppService")

	// Fetch the Route
	route := &routev1.Route{}
	err := r.client.Get(context.TODO(), request.NamespacedName, route)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not secret, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}
	if route.ObjectMeta.Labels["certificate"] == "" {
		return reconcile.Result{}, nil
	}

	// Check the certificate for the secret name
	certificate := &certv1alpha1.Certificate{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: route.ObjectMeta.Labels["certificate"], Namespace: route.Namespace}, certificate)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Create variable that is a corev1.Secret type
	secret := &corev1.Secret{}
	// Check if this Secret already exists, if it does pass all its info into the secret variable
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: certificate.Spec.SecretName, Namespace: route.Namespace}, secret)
	if err != nil {
		return reconcile.Result{}, err
	}

	cert_manager_crt := secret.Data["tls.crt"]
	cert_manager_key := secret.Data["tls.key"]
	cert_manager_ca := secret.Data["ca.crt"]

	// Does route have TLS on? if not, assume edge
	if route.Spec.TLS == nil {
		route.Spec.TLS = &routev1.TLSConfig{}
		route.Spec.TLS.Termination = "edge"
		// If passthrough, dont try and manage it
	} else if route.Spec.TLS.Termination == "passthrough" {
		reqLogger.Info("Skip reconcile: termination: passthrough not supported for  ", "Route.Namespace", route.Namespace, "Route.Name", route.Name)
		return reconcile.Result{Requeue: true}, nil
		// If TLS defined, check certs and keys to see if update is needed
	} else if route.Spec.TLS.Certificate == string(cert_manager_crt) {
		if route.Spec.TLS.Key == string(cert_manager_key) {
			if route.Spec.TLS.CACertificate == string(cert_manager_ca) {
				reqLogger.Info("Skip reconcile: Certs match for ", "Secret.Namespace", secret.Namespace, "Secret.Name", secret.Name)
				return reconcile.Result{Requeue: true}, nil
			}
		}
	}

	// If we got this far, update is needed
	route.Spec.TLS.Key = string(cert_manager_key)
	route.Spec.TLS.CACertificate = string(cert_manager_ca)
	route.Spec.TLS.Certificate = string(cert_manager_crt)

	err = r.client.Update(context.TODO(), route)
	if err != nil {
		return reconcile.Result{}, err
	}

	// Secret already exists - don't requeue
	reqLogger.Info("Reconciled: Certs updated for ", "Secret.Namespace", secret.Namespace, "Route.Name", route.Name)
	return reconcile.Result{Requeue: true}, nil
}
