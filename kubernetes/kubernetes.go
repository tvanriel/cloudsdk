package kubernetes

import (
	batchv1 "k8s.io/api/batch/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
        "context"
)


type KubernetesClient struct {
        Clientset *kubernetes.Clientset
        Namespace string

}

func NewKubernetesClient(config *Configuration) (*KubernetesClient, error) {
        clientset, err := newClientSet(config)

        if err != nil {
                return nil, err
        }

        return &KubernetesClient{
                Clientset: clientset,
                Namespace: config.Namespace,
                
        }, err
}

func newClientSet(config *Configuration) (*kubernetes.Clientset, error) {

        if config.InCluster {
                kubeconfig, err := rest.InClusterConfig()
                if err != nil {
                        return nil, err
                }
                return kubernetes.NewForConfig(kubeconfig)
        }


        restconfig, err := clientcmd.BuildConfigFromFlags("", config.Kubeconfig)

        if err != nil {
                return nil, err 
        }
        return kubernetes.NewForConfig(restconfig)
}


func (k *KubernetesClient) RunJob(job *batchv1.Job) error {
        _, err := k.Clientset.BatchV1().Jobs(k.Namespace).Create(
                context.Background(), 
                job, 
                v1.CreateOptions{},
        )
        return err
}
