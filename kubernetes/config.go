package kubernetes

type Configuration struct {
	InCluster  bool   `hcl:"in_cluster"`
	ApiServer  string `hcl:"api_server,optional"`
	Kubeconfig string `hcl:"kubeconfig,optional"`
	Namespace  string `hcl:"namespace"`
}
