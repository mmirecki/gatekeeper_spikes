package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/open-policy-agent/opa/rego"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

func getClient(kubeconfig *string) *kubernetes.Clientset {

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err.Error())
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}
	return clientset
}

func getParameters() (*string, *string) {
	var kubeconfig *string
	kubeconfigDefault := os.Getenv("KUBECONFIG")
	kubeconfig = flag.String("kubeconfig", kubeconfigDefault, "kubeconfig path")
	namespaceName := flag.String("namespace", "", "namespace name")

	flag.Parse()
	return kubeconfig, namespaceName
}

func getNamespace(nsName string, clientset *kubernetes.Clientset) *corev1.Namespace {
	ns, err := clientset.CoreV1().Namespaces().Get(context.Background(), nsName, metav1.GetOptions{})
	if err != nil {
		fmt.Printf("Error getting namespace: %+v", err)
		return nil
	}
	return ns
}

func readFile(file string) []byte {
	b, err := ioutil.ReadFile(file) // just pass the file name
	if err != nil {
		fmt.Print(err)
	}
	return b
}

func getRego() string {
	regoModule := string(readFile("namespace_rule.rego"))
	return regoModule
}

func main() {
	ctx := context.TODO()

	kubeconfig, namespaceName := getParameters()
	clientset := getClient(kubeconfig)
	ns := getNamespace(*namespaceName, clientset)

	regoModule := getRego()

	query, err := rego.New(
		rego.Query("x = data.regospike.p"),
		rego.Module("regospike", regoModule),
	).PrepareForEval(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	results, err := query.Eval(ctx, rego.EvalInput(ns))
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(results)
}
