package main

import (
	"fmt"
	"log"
	"os"
	"syscall"

	//"k8s.io/apimachinery/pkg/fields"
	"k8s.io/client-go/rest"

	"path/filepath"

	//corev1 "k8s.io/api/core/v1"
	//"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/workqueue"

	apimgmtv3 "github.com/rancher/rancher/pkg/apis/management.cattle.io/v3"
	//v3 "github.com/rancher/rancher/pkg/generated/norman/management.cattle.io/v3"

	//"k8s.io/client-go/informers"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"

	//"k8s.io/client-go/dynamic/dynamiclister"

	"os/signal"
	"time"

	"encoding/json"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
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
	//clientset, err := kubernetes.NewForConfig(config)
	clientset, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	groupVersion := schema.GroupVersion{Group: "management.cattle.io", Version: "v3"}
	dynamicInformerFactory := dynamicinformer.NewDynamicSharedInformerFactory(clientset, time.Minute*5)

	informer := dynamicInformerFactory.ForResource(schema.GroupVersionResource{
		Group:    groupVersion.Group,
		Version:  groupVersion.Version,
		Resource: "nodes",
	}).Informer()

	informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
		// AddFunc: func(obj interface{}) {
		// 	fmt.Println("New custom resource added")
		// },
		// UpdateFunc: func(oldObj, newObj interface{}) {
		// 	fmt.Println("Custom resource updated")
		// },
		DeleteFunc: func(obj interface{}) {
			var tmp []byte
			var err error
			var v3Node interface{}
			tmp, err = obj.(*unstructured.Unstructured).MarshalJSON()
			if err != nil {
				return
			}
			err = json.Unmarshal(tmp, v3Node)
			log.Println(v3Node)
			v3Node2, ok := v3Node.(*apimgmtv3.Node)
			if !ok {
				fmt.Println("Failed to convert to CustomResourceDefinition")
				return
			}

			//log.Println(obj)
			// switch v := obj.(type) {
			// case *apimgmtv3.Node:
			// 	fmt.Println("obj의 타입은 apimgmtv3.Node입니다.")
			// 	// obj가 apimgmtv3.Node 타입인 경우에는 v를 사용하여 작업할 수 있습니다.
			// default:
			// 	fmt.Printf("obj의 타입은 %T입니다.\n", v)
			// }
			// crd, ok := obj.(*apimgmtv3.Node)
			// if !ok {
			// 	fmt.Println("Failed to convert to CustomResourceDefinition")
			// 	return
			// }
			fmt.Printf("Custom resource %s deleted", v3Node2.Name)
		},
	})

	// 인포머 시작
	stopCh := make(chan struct{})
	defer close(stopCh)
	dynamicInformerFactory.Start(stopCh)

	// 시그널 핸들링
	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-signalCh:
		fmt.Println("Received termination signal, shutting down...")
		close(stopCh) // 인포머를 멈추기 위해 stopCh 채널을 닫음
		<-stopCh      // 인포머가 종료될 때까지 대기
	case <-time.After(time.Minute): // 1분(60초) 대기 후 종료
		fmt.Println("Program has run for 1 minute, shutting down...")
		close(stopCh) // 인포머를 멈추기 위해 stopCh 채널을 닫음
	}

	// // 큐 생성
	// queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	// // 노드 이벤트 핸들러 생성
	// nodeEventHandler := &NodeEventHandler{
	// 	queue: queue,
	// }

	// // 노드 리소스를 감시하기 위한 채널을 생성합니다.
	// //watchlist := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "nodes", corev1.NamespaceAll, fields.Everything())
	// //watchlist := cache.NewListWatchFromClient(clientset.CoreV1().RESTClient(), "nodes.management.cattle.io", corev1.NamespaceAll, fields.Everything())

	// informerFactory := dynamicinformer.NewDynamicSharedInformerFactory(clientset, 0)
	// nodeInformer := informerFactory.Management().V3().Nodes().Informer()

	// // 노드 삭제 이벤트를 탐지하기 위한 채널을 생성합니다.
	// _, controller := cache.NewInformer(
	// 	watchlist,
	// 	&apimgmtv3.Cluster{},
	// 	0, // 캐시 크기
	// 	cache.ResourceEventHandlerFuncs{
	// 		DeleteFunc: nodeEventHandler.OnDelete,
	// 		// DeleteFunc: func(obj interface{}) {
	// 		// node := obj.(*corev1.Node)
	// 		// fmt.Printf("Node %s deleted\n", node.Name)
	// 		//},
	// 	},
	// )

	// crdInformerFactory := informers.NewSharedInformerFactory(clientset, 0)
	// crdInformer := crdInformerFactory.Apiextensions().V1().CustomResourceDefinitions().Informer()

	// crdInformer.AddEventHandler(cache.ResourceEventHandlerFuncs{
	// 	DeleteFunc: func(obj interface{}) {
	// 		crd, ok := obj.(*v1.CustomResourceDefinition)
	// 		if !ok {
	// 			fmt.Println("Failed to convert to CustomResourceDefinition")
	// 			return
	// 		}
	// 		fmt.Printf("CRD %s deleted\n", crd.Name)
	// 	},
	// })

	// 컨트롤러를 시작합니다.
	// stop := make(chan struct{})
	// defer close(stop)
	// go controller.Run(stop)

	// // 큐에서 노드 정보 처리
	// for {
	// 	key, quit := queue.Get()
	// 	if quit {
	// 		break
	// 	}
	// 	fmt.Printf("Node deleted: %v\n", key)
	// 	queue.Done(key)
	// }

	// // 앱이 종료되지 않도록 대기합니다.
	// select {}
}
