package main

import (
	"context"
	"fmt"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	_ "k8s.io/client-go/plugin/pkg/client/auth/oidc"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

func main() {
	//Load the kube config file
	fmt.Println("Getting Kube config file...")
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", filepath.Join(homedir.HomeDir(), ".kube", "config"))
	if err != nil {
		panic(err)
	}
	//Connect to Kubernetes
	fmt.Println("Connecting to Kubernetes...")
	clientset, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		panic(err)
	}

	//Create Configmap
	configmap, err := CreateCustomConfigMap(clientset, "nginx-configmap", map[string]string{"sample.txt": "Hello Team!"})
	if err != nil {
		panic(err)
	}
	fmt.Println("Created Configmap:\n", configmap.ObjectMeta.Name)

	//Create Secret
	secret, err := CreateCustomSecret(clientset, "nginx-secret", map[string][]byte{
		"username": []byte("c3VkaGFyc2hhbgo="),
		"password": []byte("MTIzNDU2Cg=="),
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Created Secret:\n", secret.ObjectMeta.Name)

	//Create Pods
	pod1, err := CreateCustomPodWithConfigandSecret(clientset, "nginx-pod", "nginx-container", "nginx", "nginx-configmap", "nginx-secret")
	if err != nil {
		panic(err)
	}
	fmt.Println("Created POD:\n", pod1.ObjectMeta.Name)
}

// Function for creating custom pod
func CreateCustomPod(clientset *kubernetes.Clientset, podname string, containername string, imagename string) (*v1.Pod, error) {
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podname,
			Namespace: "default",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  containername,
					Image: imagename,
				},
			},
			RestartPolicy: v1.RestartPolicyNever,
		},
	}
	createdpod, err := clientset.CoreV1().Pods("default").Create(context.Background(), pod, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	return createdpod, err
}

// Function for creating custom pod
func CreateCustomPodWithConfigandSecret(clientset *kubernetes.Clientset, podname string, containername string, imagename string, configmap string, secret string) (*v1.Pod, error) {
	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      podname,
			Namespace: "default",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  containername,
					Image: imagename,
					VolumeMounts: []v1.VolumeMount{
						{
							Name:      "configmap-volume",
							MountPath: "/sudharshan/configmapfolder",
						},
						{
							Name:      "secret-volume",
							MountPath: "/sudharshan/secret",
						},
					},
				},
			},
			Volumes: []v1.Volume{
				{
					Name: "configmap-volume",
					VolumeSource: v1.VolumeSource{
						ConfigMap: &v1.ConfigMapVolumeSource{
							LocalObjectReference: v1.LocalObjectReference{
								Name: configmap,
							},
						},
					},
				},
				{
					Name: "secret-volume",
					VolumeSource: v1.VolumeSource{
						Secret: &v1.SecretVolumeSource{
							SecretName: secret,
						},
					},
				},
			},
			RestartPolicy: v1.RestartPolicyNever,
		},
	}
	createdpod, err := clientset.CoreV1().Pods("default").Create(context.Background(), pod, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	return createdpod, err
}

// Function for creating custom configmap
func CreateCustomConfigMap(clienset *kubernetes.Clientset, configmapname string, configmapdata map[string]string) (*v1.ConfigMap, error) {
	configmap := &v1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      configmapname,
			Namespace: "default",
		},
		Data: configmapdata,
	}
	createdconfigmap, err := clienset.CoreV1().ConfigMaps("default").Create(context.Background(), configmap, metav1.CreateOptions{})
	return createdconfigmap, err
}

// Function for creating custom secret
func CreateCustomSecret(clienset *kubernetes.Clientset, secretname string, secretdata map[string][]byte) (*v1.Secret, error) {
	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      secretname,
			Namespace: "default",
		},
		Data: secretdata,
		Type: v1.SecretTypeOpaque,
	}
	createdSecret, err := clienset.CoreV1().Secrets("default").Create(context.Background(), secret, metav1.CreateOptions{})
	return createdSecret, err
}
