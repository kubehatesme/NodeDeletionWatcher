package main

import (
	"context"
	"fmt"
	"os"
	"time"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/tools/clientcmd"
)

func test() {
	// Kubernetes 클러스터에 연결하기 위한 kubeconfig 파일 경로
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = clientcmd.RecommendedHomeFile // 기본 kubeconfig 파일 경로
	}

	// kubeconfig 파일로부터 클라이언트 설정 생성
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		panic(err)
	}

	// 동적 클라이언트 생성
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	// CRD 리소스 정의
	crdResource := schema.GroupVersionResource{
		Group:    "management.cattle.io",
		Version:  "v3",
		Resource: "customresourcedefinitions",
	}

	// CRD 삭제 이벤트 감지
	watcher, err := dynamicClient.Resource(crdResource).Watch(context.Background(), v1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for {
		select {
		case event, ok := <-watcher.ResultChan():
			if !ok {
				fmt.Println("Watcher channel closed")
				return
			}

			if event.Type == watch.Deleted {
				crd := event.Object.(*unstructured.Unstructured)
				// 객체를 맵으로 변환
				crdMap := crd.UnstructuredContent()
				if err != nil {
					fmt.Printf("Error converting CRD to map: %v\n", err)
					continue
				}
				// kind 필드 확인
				kind, ok := crdMap["kind"].(string)
				if !ok {
					fmt.Println("Kind field not found or not a string")
					continue
				}
				// kind가 "test"인 경우 삭제 이벤트 출력
				if kind == "Node" {
					fmt.Printf("CRD %s deleted\n", crd.GetName())
				}
			}
		case <-time.After(time.Minute * 10): // 10분 후에 프로그램 종료
			fmt.Println("Timeout reached, exiting...")
			return
		}
	}
}
