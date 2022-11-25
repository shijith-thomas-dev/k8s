package main

import (
	"context"
	"flag"
	"log"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func elementExists(s []string, element string) bool {
	for i := range s {
		if s[i] == element {
			return true
		}
	}
	return false
}
func getDuplicates(inp map[string][]string, n int) {

	data := make([]string, 0)
	for _, s := range inp {
		// log.Printf("Checking %v node for more than %d  deployment pods \n", node, n)
		result := make(map[string]int)
		for i := range s {
			if v, ok := result[s[i]]; ok {
				result[s[i]] = v + 1
			} else {
				result[s[i]] = 1
			}
		}
		for i, v := range result { //map[deploymentb:2 deploymentg:2 deploymentr:1]
			check := elementExists(data, i)
			if !check {
				if v >= n {
					data = append(data, i)
				}
			}
		}
	}
	log.Println(data)
}

func main() {
	var kubeconfig *string

	nodeMap := make(map[string][]string)
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	noOfPodsFlag := flag.Int("n", 2, "no of  deployment pods")

	flag.Parse()
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	// create the kubeClient
	kubeClient, err := kubernetes.NewForConfig(config)

	if err != nil {
		panic(err.Error())
	}

	nodeClient, _ := kubeClient.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}

	for _, nodes := range nodeClient.Items {
		nodeMap[nodes.Name] = []string{}
	}
	// log.Println(nodeMap)

	pods, err := kubeClient.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		panic(err.Error())
	}

	for _, pod := range pods.Items {
		if len(pod.OwnerReferences) == 0 {
			log.Printf("Pod %s has no owner", pod.Name)
			continue
		}

		var ownerName string

		switch pod.OwnerReferences[0].Kind {
		case "ReplicaSet":
			replica, repErr := kubeClient.AppsV1().ReplicaSets(pod.Namespace).Get(context.TODO(), pod.OwnerReferences[0].Name, metav1.GetOptions{})
			if repErr != nil {
				panic(repErr.Error())
			}

			ownerName = replica.OwnerReferences[0].Name

		case "DaemonSet", "StatefulSet":
			continue
		default:
			log.Printf("Could not find resource manager for type %s\n", pod.OwnerReferences[0].Kind)
			continue
		}
		// log.Printf("POD %s is managed by %s %s\n", pod.Name, ownerName, ownerKind)
		nodeMap[pod.Spec.NodeName] = append(nodeMap[pod.Spec.NodeName], ownerName)

	}
	// log.Println(nodeMap)
	getDuplicates(nodeMap, *noOfPodsFlag)
}
