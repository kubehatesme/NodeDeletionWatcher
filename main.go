package main

import (
	"fmt"

	//"path/filepath"

	//v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"

	//"k8s.io/client-go/kubernetes"

	//"k8s.io/client-go/util/homedir"

	"github.com/rancher/rancher/pkg/generated/controllers/management.cattle.io"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	//"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	//"k8s.io/apimachinery/pkg/runtime"

	"log"
	//"k8s.io/apimachinery/pkg/runtime/schema"
)

func main() {
	// 홈 디렉터리에서 kubeconfig 경로 설정
	log.Println("node cleanup pod starting...")
	//var config *rest.Config
	config, err := rest.InClusterConfig()

	// if kubeconfigPath := os.Getenv("KUBECONFIG"); kubeconfigPath != "" {
	// 	//   2. kubeconfig file located by "KUBECONFIG" env
	// 	config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	// } else {
	// 	//   3. in-cluster client configuration (useful when using detek in a kubernetes cluster)
	// 	config, err = rest.InClusterConfig()
	// }
	// if err != nil {
	// 	//   4. kubeconfig file located in default directory ($HOME/.kube/config)`,
	// 	kubeconfigPath := filepath.Join(homedir.HomeDir(), ".kube", "config")
	// 	config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	//}

	// kubeconfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
	// 	clientcmd.NewDefaultClientConfigLoadingRules(),
	// 	&clientcmd.ConfigOverrides{},
	// )
	//nconfig, nerr := kubeconfig.ClientConfig()

	log.Println(config)
	if err != nil {
		panic(err)
	}

	// if err != nil {
	// 	fmt.Printf("Error loading in-cluster configuration: %s\n", err.Error())
	// 	os.Exit(1)
	// }
	log.Println("here...")

	// 클라이언트셋 생성
	rancherManagement, err := management.NewFactoryFromConfig(config)
	if err != nil {
		panic(err.Error())
	}
	log.Println(rancherManagement)

	// informer := rancherManagement.Management().V3().Node().Informer()
	// log.Println(informer)

	test, err := rancherManagement.Management().V3().Node().List("", metav1.ListOptions{})
	log.Println(test)
	if err != nil {
		panic(err.Error())
	}
	for _, node := range test.Items {
		fmt.Println(node.Name)
	}

	//groupVersion := schema.GroupVersion{Group: "management.cattle.io", Version: "v3"}
	// ControllerFactory().SharedCacheFactory().ForResource(schema.GroupVersionResource{
	// 	Group:    groupVersion.Group,
	// 	Version:  groupVersion.Version,
	// 	Resource: "nodes",},true).

	// clientset, err := kubernetes.NewForConfig(config)
	// if err != nil {
	// 	panic(err.Error())
	// }

	// // 모든 네임스페이스를 나열
	// namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	// if err != nil {
	// 	panic(err.Error())
	// }

	// fmt.Println("Namespaces in the cluster:")
	// for _, ns := range namespaces.Items {
	// 	fmt.Println(ns.Name)
	// }

	// informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
	// 	DeleteFunc: func(obj interface{}) {
	// 		log.Println("here2...")
	// 		v3Node, _ := obj.(*apimgmtv3.Node)
	// 		//[Delete h53] Set variables
	// 		hostName := v3Node.Status.NodeLabels["hostName"]
	// 		address := v3Node.Status.InternalNodeStatus.Addresses[0].Address
	// 		site := v3Node.Status.NodeLabels["site"][:len(v3Node.Status.NodeLabels["site"])-4]
	// 		zone := v3Node.Status.NodeLabels["site"][len(v3Node.Status.NodeLabels["site"])-4:]
	// 		url := "gitlab.arc.hcloud.io/common/rancher/-/raw/master"
	// 		//zone := os.Getenv("ZONE_TYPE")
	// 		//url := os.Getenv("SCRIPT_URL")

	// 		//[Delete EAI] Set variables
	// 		boxName := v3Node.Status.NodeLabels["boxName"]
	// 		creationTimestamp := v3Node.Status.NodeLabels["vmCreationTime"]
	// 		memory := v3Node.Status.NodeLabels["memory"]
	// 		osType := v3Node.Status.NodeLabels["osType"]
	// 		vCpus := v3Node.Status.NodeLabels["vCpus"]
	// 		projectId := "548"
	// 		//projectId := os.Getenv("PROJECT_ID")

	// 		cmd := fmt.Sprintf("url=%s; site=%s; curl -sfL https://$url/post-scripts/$site/cleanup.sh | VM_HOSTNAME=%s ZONE_TYPE=%s VM_IP_INFO=%s VM_HYPERVISOR=%s VM_CREATED_DATE=%s VM_MEMORY=%s OS_TYPE=%s VM_VCORE=%s PROJECT_ID=%s sh -",
	// 			url, site, hostName, zone, address, boxName, creationTimestamp, memory, osType, vCpus, projectId)
	// 		log.Println("[Command] \"", cmd, "\"")
	// 		output, err := exec.Command("bash", "-c", cmd, "&& sleep infinity").Output()
	// 		if err != nil {
	// 			log.Printf("error %s", err)
	// 		} else {
	// 			log.Printf("[Response] %s\n", string(output))
	// 		}
	// 		//fmt.Printf("Deleted Node's real name: %s\n", hostName)
	// 	},
	// })

	//인포머 시작
	// stopCh := make(chan struct{})
	// defer close(stopCh)
	// //dynamicInformerFactory.Start(stopCh)
	// go informer.Run(stopCh)

	// 시그널 핸들링
	// signalCh := make(chan os.Signal, 1)
	// signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	select {}
}
