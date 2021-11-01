# sample-credential-provider

Sample exec based external kubelet image credential provider. 

To fetch credential kubelet executes `sample-credential-provider get-credentials` command with passing the api request as json-serialized api

```
echo '{"kind": "CredentialProviderRequest", "apiVersion": "credentialprovider.kubelet.k8s.io/v1alpha1", "image": "gcr.io/authenticated-image-pulling/alpine:3.7"}' | sample-credential-provider get-credentials
```
on recieving kubelet requests `sample-credential-provider` reads authentication information from docker config.json and emits the response in the form json-serialized api which is then read and parsed by kubelet.

e.g. response
```
{"kind":"CredentialProviderResponse","apiVersion":"credentialprovider.kubelet.k8s.io/v1alpha1","cacheKeyType":"Registry","auth":{<auth-info>}}
```
> **_NOTE:_**  This provider is intended to be used only in test environments

If you want to read and understand more about kubelet credential providers follow below links

- **Doc**: https://kubernetes.io/docs/tasks/kubelet-credential-provider/kubelet-credential-provider/
- **KEP**: https://github.com/kubernetes/enhancements/tree/master/keps/sig-node/2133-kubelet-credential-providers
- **API**: https://github.com/kubernetes/kubernetes/tree/master/staging/src/k8s.io/kubelet/pkg/apis/credentialprovider/v1alpha1