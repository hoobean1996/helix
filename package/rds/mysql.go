package rds

import (
	"context"
	"fmt"

	"helix.io/helix/package/common"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/kubernetes"
)

type MysqlManager struct {
	clientset *kubernetes.Clientset
}

func NewMysqlManager(clientset *kubernetes.Clientset) *MysqlManager {

	manager := &MysqlManager{
		clientset: clientset,
	}
	return manager
}

func (m *MysqlManager) Deploy(spec *MySQLSpecification) {
	m.CreateMySQLDeployment(spec)
	m.CreateMySQLService(spec)
}

func (m *MysqlManager) CreateMySQLDeployment(spec *MySQLSpecification) {
	deployment := &appsv1.Deployment{
		// 直接引用资源
		// 命名空间中唯一
		// 不能修改
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("mysql-%s", spec.ProjectName),
			Namespace: "default",
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: common.Int32Ptr(1),
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
							Name:  "mysql",
							Image: "mysql:8.0",
							Ports: []corev1.ContainerPort{
								{
									ContainerPort: 3306,
								},
							},
							Env: []corev1.EnvVar{
								{Name: "MYSQL_ROOT_PASSWORD", Value: spec.Password},
								{Name: "MYSQL_DATABASE", Value: spec.Database},
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
							// 用途：检测容器是否还在运行
							LivenessProbe: &corev1.Probe{
								InitialDelaySeconds: 30,
								PeriodSeconds:       10,
								ProbeHandler: corev1.ProbeHandler{
									TCPSocket: &corev1.TCPSocketAction{
										Port: intstr.FromInt(3306),
									},
								},
							},
							// 用途：检测容器是否准备好接受流量
							ReadinessProbe: &corev1.Probe{
								InitialDelaySeconds: 5,
								PeriodSeconds:       2,
								ProbeHandler: corev1.ProbeHandler{
									Exec: &corev1.ExecAction{
										Command: []string{
											"mysqladmin",
											"ping",
										},
									},
								},
							},
							ImagePullPolicy: corev1.PullIfNotPresent,
						},
					},
				},
			},
		},
	}

	namespace := "default"
	if _, err := m.clientset.AppsV1().Deployments(namespace).Create(
		context.Background(),
		deployment,
		metav1.CreateOptions{},
	); err != nil {
		fmt.Printf("create deployment failed: err=%s", err.Error())
	}
}

func (m *MysqlManager) CreateMySQLService(spec *MySQLSpecification) {
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("mysql-%s", spec.ProjectName),
			Namespace: "default",
		},
		Spec: corev1.ServiceSpec{
			Selector: map[string]string{
				"project": spec.ProjectName,
				"env":     spec.Env,
			},
			Ports: []corev1.ServicePort{
				{
					Port:       3306,
					TargetPort: intstr.FromInt(3306),
				},
			},
		},
	}

	if _, err := m.clientset.CoreV1().Services("default").Create(context.Background(), service, metav1.CreateOptions{}); err != nil {
		fmt.Printf("create service failed, err=%s", err.Error())
	}
}
