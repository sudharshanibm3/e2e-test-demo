package main

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func TestCreatePod(t *testing.T) {
	//Load the kube config file
	fmt.Println("Getting Kube config file...")
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		t.Error(err)
	}
	//Connect to Kubernetes
	fmt.Println("Connecting to Kubernetes...")
	clientset, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		t.Error(err)
	}
	podname := "test-pod"
	contianername := "test-container"
	imagname := "quay.io/sudharshanibm3/test-images:testuser"
	pod, err := CreateCustomPod(clientset, podname, contianername, imagname)
	assert.Nil(t, err)
	assert.Equal(t, pod.ObjectMeta.Name, podname)
	assert.Equal(t, pod.Spec.Containers[0].Name, contianername)
	assert.Equal(t, pod.Spec.Containers[0].Image, imagname)
	err = clientset.CoreV1().Pods("default").Delete(context.Background(), podname, v1.DeleteOptions{})
	if err != nil {
		t.Error(err)
	}

}
func TestCreateConfigmap(t *testing.T) {
	//Load the kube config file
	fmt.Println("Getting Kube config file...")
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		t.Error(err)
	}
	//Connect to Kubernetes
	fmt.Println("Connecting to Kubernetes...")
	clientset, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		t.Error(err)
	}
	configmapname := "test-configmap"
	data := map[string]string{
		"testfile.txt": "hello test!",
	}
	configmap, err := CreateCustomConfigMap(clientset, configmapname, data)
	assert.Nil(t, err)
	assert.Equal(t, configmap.ObjectMeta.Name, configmapname)
	assert.Equal(t, configmap.Data, data)
	err = clientset.CoreV1().ConfigMaps("default").Delete(context.Background(), configmapname, v1.DeleteOptions{})
	if err != nil {
		t.Error(err)
	}

}
func TestCreateSecret(t *testing.T) {
	//Load the kube config file
	fmt.Println("Getting Kube config file...")
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		t.Error(err)
	}
	//Connect to Kubernetes
	fmt.Println("Connecting to Kubernetes...")
	clientset, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		t.Error(err)
	}
	secretname := "test-secret"
	data := map[string][]byte{
		"username": []byte("c3VkaGFyc2hhbgo="),
		"password": []byte("MTIzNDU2Cg=="),
	}
	secret, err := CreateCustomSecret(clientset, secretname, data)
	assert.Nil(t, err)
	assert.Equal(t, secret.ObjectMeta.Name, secretname)
	assert.Equal(t, secret.Data, data)
	err = clientset.CoreV1().Secrets("default").Delete(context.Background(), secretname, v1.DeleteOptions{})
	if err != nil {
		t.Error(err)
	}

}
