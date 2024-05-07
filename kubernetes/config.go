package kubernetes

type Configuration struct {
	InCluster  bool   `hcl:"in_cluster"`
	ApiServer  string `hcl:"api_server"`
	Kubeconfig string `hcl:"kubeconfig"`
	Namespace  string `hcl:"namespace"`
}
