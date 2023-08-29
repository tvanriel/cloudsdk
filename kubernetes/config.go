package kubernetes

type Configuration struct {
        InCluster bool
        ApiServer string
        Kubeconfig string
        Namespace string
}
