package operator

import (
	"fmt"
	"os"
	"time"
	"github.com/openshift/cluster-svcat-apiserver-operator/pkg/operator/operatorclient"
	"github.com/openshift/cluster-svcat-apiserver-operator/pkg/operator/resourcesynccontroller"
	"github.com/openshift/cluster-svcat-apiserver-operator/pkg/operator/workloadcontroller"
	configv1 "github.com/openshift/api/config/v1"
	configv1client "github.com/openshift/client-go/config/clientset/versioned"
	configinformers "github.com/openshift/client-go/config/informers/externalversions"
	operatorv1client "github.com/openshift/client-go/operator/clientset/versioned"
	operatorv1informers "github.com/openshift/client-go/operator/informers/externalversions"
	"github.com/openshift/library-go/pkg/controller/controllercmd"
	"github.com/openshift/library-go/pkg/operator/status"
	"github.com/openshift/library-go/pkg/operator/unsupportedconfigoverridescontroller"
	"github.com/openshift/library-go/pkg/operator/v1helpers"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	apiregistrationclient "k8s.io/kube-aggregator/pkg/client/clientset_generated/clientset"
	apiregistrationinformers "k8s.io/kube-aggregator/pkg/client/informers/externalversions"
)

func RunOperator(ctx *controllercmd.ControllerContext) error {
	_logClusterCodePath()
	defer _logClusterCodePath()
	_logClusterCodePath()
	defer _logClusterCodePath()
	kubeClient, err := kubernetes.NewForConfig(ctx.ProtoKubeConfig)
	if err != nil {
		return err
	}
	apiregistrationv1Client, err := apiregistrationclient.NewForConfig(ctx.ProtoKubeConfig)
	if err != nil {
		return err
	}
	operatorConfigClient, err := operatorv1client.NewForConfig(ctx.KubeConfig)
	if err != nil {
		return err
	}
	configClient, err := configv1client.NewForConfig(ctx.KubeConfig)
	if err != nil {
		return err
	}
	operatorConfigInformers := operatorv1informers.NewSharedInformerFactory(operatorConfigClient, 10*time.Minute)
	kubeInformersForNamespaces := v1helpers.NewKubeInformersForNamespaces(kubeClient, "", operatorclient.GlobalUserSpecifiedConfigNamespace, operatorclient.GlobalMachineSpecifiedConfigNamespace, operatorclient.KubeAPIServerNamespaceName, operatorclient.OperatorNamespace, operatorclient.TargetNamespaceName)
	apiregistrationInformers := apiregistrationinformers.NewSharedInformerFactory(apiregistrationv1Client, 10*time.Minute)
	configInformers := configinformers.NewSharedInformerFactory(configClient, 10*time.Minute)
	operatorClient := &operatorclient.OperatorClient{Informers: operatorConfigInformers, Client: operatorConfigClient.OperatorV1()}
	resourceSyncController, err := resourcesynccontroller.NewResourceSyncController(operatorClient, kubeInformersForNamespaces, v1helpers.CachedConfigMapGetter(kubeClient.CoreV1(), kubeInformersForNamespaces), v1helpers.CachedSecretGetter(kubeClient.CoreV1(), kubeInformersForNamespaces), ctx.EventRecorder)
	if err != nil {
		return err
	}
	versionRecorder := status.NewVersionGetter()
	clusterOperator, err := configClient.ConfigV1().ClusterOperators().Get("service-catalog-apiserver", metav1.GetOptions{})
	if err != nil && !errors.IsNotFound(err) {
		return err
	}
	for _, version := range clusterOperator.Status.Versions {
		versionRecorder.SetVersion(version.Name, version.Version)
	}
	versionRecorder.SetVersion("operator", os.Getenv("OPERATOR_IMAGE_VERSION"))
	workloadController := workloadcontroller.NewWorkloadController(os.Getenv("IMAGE"), versionRecorder, operatorConfigInformers.Operator().V1().ServiceCatalogAPIServers(), kubeInformersForNamespaces.InformersFor(operatorclient.TargetNamespaceName), kubeInformersForNamespaces.InformersFor(operatorclient.GlobalUserSpecifiedConfigNamespace), kubeInformersForNamespaces.InformersFor(operatorclient.KubeAPIServerNamespaceName), kubeInformersForNamespaces.InformersFor(operatorclient.GlobalUserSpecifiedConfigNamespace), apiregistrationInformers, configInformers, operatorConfigClient.OperatorV1(), configClient.ConfigV1(), kubeClient, apiregistrationv1Client.ApiregistrationV1(), ctx.EventRecorder)
	finalizerController := NewFinalizerController(kubeInformersForNamespaces.InformersFor(operatorclient.TargetNamespaceName), kubeClient.CoreV1(), ctx.EventRecorder)
	clusterOperatorStatus := status.NewClusterOperatorStatusController("service-catalog-apiserver", append([]configv1.ObjectReference{{Resource: "namespaces", Name: operatorclient.GlobalUserSpecifiedConfigNamespace}, {Resource: "namespaces", Name: operatorclient.GlobalMachineSpecifiedConfigNamespace}, {Resource: "namespaces", Name: operatorclient.OperatorNamespace}, {Resource: "namespaces", Name: operatorclient.TargetNamespaceName}}, workloadcontroller.APIServiceReferences()...), configClient.ConfigV1(), configInformers.Config().V1().ClusterOperators(), operatorClient, versionRecorder, ctx.EventRecorder)
	configUpgradeableController := unsupportedconfigoverridescontroller.NewUnsupportedConfigOverridesController(operatorClient, ctx.EventRecorder)
	operatorConfigInformers.Start(ctx.Done())
	kubeInformersForNamespaces.Start(ctx.Done())
	apiregistrationInformers.Start(ctx.Done())
	configInformers.Start(ctx.Done())
	go workloadController.Run(1, ctx.Done())
	go clusterOperatorStatus.Run(1, ctx.Done())
	go finalizerController.Run(1, ctx.Done())
	go resourceSyncController.Run(1, ctx.Done())
	go configUpgradeableController.Run(1, ctx.Done())
	<-ctx.Done()
	return fmt.Errorf("stopped")
}
