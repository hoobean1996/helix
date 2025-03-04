package blob

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

type BlobManager struct {
	clientset *kubernetes.Clientset
}

func NewBlobManager(clientset *kubernetes.Clientset) *BlobManager {
	manager := &BlobManager{
		clientset: clientset,
	}
	return manager
}

func (b *BlobManager) Deploy(spec *BlobSpecification) {
	b.CreateMinioPVC(spec)
	b.CreateBlobDeployment(spec)
	b.CreateBlockService(spec)
}

func (b *BlobManager) CreateMinioPVC(spec *BlobSpecification) {
	pvc := &corev1.PersistentVolumeClaim{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("minio-pvc-%s", spec.ProjectName),
			Namespace: "default",
			Labels: map[string]string{
				"project": spec.ProjectName,
				"env":     spec.Env,
			},
		},
		Spec: corev1.PersistentVolumeClaimSpec{
			AccessModes: []corev1.PersistentVolumeAccessMode{
				corev1.ReadWriteOnce,
			},
			Resources: corev1.VolumeResourceRequirements{
				Requests: corev1.ResourceList{
					corev1.ResourceStorage: resource.MustParse("10Gi"), // 请求 10GB 存储
				},
			},
		},
	}
	if _, err := b.clientset.CoreV1().PersistentVolumeClaims("default").Create(
		context.Background(),
		pvc,
		metav1.CreateOptions{},
	); err != nil {
		fmt.Printf("create pvc failed, err=%s\n", err.Error())
	}
}

func (b *BlobManager) CreateBlobDeployment(spec *BlobSpecification) {
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("minio-%s", spec.ProjectName),
			Namespace: "default",
		},
		Spec: appsv1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: map[string]string{
					"project": spec.ProjectName,
					"env":     spec.Env,
				},
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"project": spec.ProjectName,
						"env":     spec.Env,
					},
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:  "minio",
							Image: "minio/minio:latest",
							Args: []string{
								"server",
								"/data",
								"--console-address",
								":9001",
							},
							Env: []corev1.EnvVar{
								{Name: "MINIO_ACCESS_KEY", Value: spec.AccessKey},
								{Name: "MINIO_SECRET_KEY", Value: spec.SecretKey},
							},
							Ports: []corev1.ContainerPort{
								{ContainerPort: 9000}, // API 端口
								{ContainerPort: 9001}, // Console 端口
							},
							Resources: corev1.ResourceRequirements{
								Requests: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse(spec.RequestCPU),
									corev1.ResourceMemory: resource.MustParse(spec.RequestMemory),
								},
								Limits: corev1.ResourceList{
									corev1.ResourceCPU:    resource.MustParse(spec.LimitCPU),
									corev1.ResourceMemory: resource.MustParse(spec.LimitMemory),
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name:      "data",
									MountPath: "/data",
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "data",
							VolumeSource: corev1.VolumeSource{
								PersistentVolumeClaim: &corev1.PersistentVolumeClaimVolumeSource{
									ClaimName: fmt.Sprintf("minio-pvc-%s", spec.ProjectName),
								},
							},
						},
					},
				},
			},
		},
	}

	namespace := "default"
	if _, err := b.clientset.AppsV1().Deployments(namespace).Create(
		context.Background(),
		deployment,
		metav1.CreateOptions{},
	); err != nil {
		fmt.Printf("create deployment failed: err=%s\n", err.Error())
	}
}

func (b *BlobManager) CreateBlockService(spec *BlobSpecification) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("blob-%s", spec.ProjectName),
			Namespace: "default",
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"project": spec.ProjectName,
				"env":     spec.Env,
			},
			Ports: []corev1.ServicePort{
				{Name: "api", Port: 9000, TargetPort: intstr.FromInt(9000)},
				{Name: "console", Port: 9001, TargetPort: intstr.FromInt(9001)},
			},
		},
	}
	if _, err := b.clientset.CoreV1().Services("default").Create(
		context.Background(),
		service,
		metav1.CreateOptions{},
	); err != nil {
		fmt.Printf("create service failed, err=%s", err.Error())
	}
}
