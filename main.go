package main

import (
	"fmt"
	"os"

	"path/filepath"

	"k8s.io/client-go/rest"

	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/workqueue"

	apimgmtv3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	"github.com/rancher/rancher/pkg/generated/controllers/management.cattle.io"
	//"log"
	//"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	//"k8s.io/apimachinery/pkg/runtime"
)

// NodeEventHandler implements cache.ResourceEventHandler
type NodeEventHandler struct {
	queue workqueue.RateLimitingInterface
}

func (n *NodeEventHandler) OnDelete(obj interface{}) {
	// 노드 삭제 이벤트 처리
	//node, ok := obj.(*corev1.Node)
	node, ok := obj.(*apimgmtv3.Node)
	if !ok {
		return
	}
	fmt.Printf("Node %s deleted\n", node.Name)

	n.queue.Add(node.ObjectMeta.Name)
}

func main() {
	// 홈 디렉터리에서 kubeconfig 경로 설정
	var config *rest.Config
	var err error

	if kubeconfigPath := os.Getenv("KUBECONFIG"); kubeconfigPath != "" {
		//   2. kubeconfig file located by "KUBECONFIG" env
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	} else {
		//   3. in-cluster client configuration (useful when using detek in a kubernetes cluster)
		config, err = rest.InClusterConfig()
	}
	if err != nil {
		//   4. kubeconfig file located in default directory ($HOME/.kube/config)`,
		kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	}

	if err != nil {
		fmt.Printf("Error loading in-cluster configuration: %s\n", err.Error())
		os.Exit(1)
	}

	// 클라이언트셋 생성
	rancherManagement, err := management.NewFactoryFromConfig(config)
	if err != nil {
		panic(err.Error())
	}

	informer := rancherManagement.Management().V3().Node().Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		DeleteFunc: func(obj interface{}) {
			v3Node, _ := obj.(*apimgmtv3.Node)
			fmt.Printf("Deleted Node's m-name: %s\n", v3Node.Name)
			fmt.Printf("Deleted Node's real name: %s\n", v3Node.Status.NodeLabels["kubernetes.io/hostname"])
			fmt.Printf("Deleted Node's osType: %s\n", v3Node.Status.NodeLabels["osType"])
		},
	})

	// 인포머 시작
	stopCh := make(chan struct{})
	defer close(stopCh)
	//dynamicInformerFactory.Start(stopCh)
	go informer.Run(stopCh)

	// 시그널 핸들링
	// signalCh := make(chan os.Signal, 1)
	// signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	select {}
}
