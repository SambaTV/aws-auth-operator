# sambatv/aws-auth-operator

This repository contains the Golang implementation of a [Kubernetes Operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
managing the `aws-auth` ConfigMap, built with [Kubebuilder](https://kubebuilder.io/)

## Custom Resource Definitions

This operator provides the following CRD kinds in the `aws-auth.samba.tv` API group.

- [MapRole](config/samples/maprole.yaml)
- [MapUser](config/samples/mapuser.yaml)

## External Resources

- [Kubebuilder documentation](https://book.kubebuilder.io/)
- [kubernetes Custom Resources documentation](https://kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)
- [Kubernetes Controllers documentation](https://kubernetes.io/docs/concepts/architecture/controller/)
- [Kubernetes Operator documentation](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/)
- [AWS EKS aws-auth ConfigMap documentation](https://docs.aws.amazon.com/eks/latest/userguide/add-user-role.html)

## Related work

- [ops42 aws-auth-operator](https://ops42.org/aws-auth-operator/) provides a
  way to map AWS IAM users to the `data.mapUsers` section of the `kube-system:aws-auth`
  ConfigMap.
- [rustrial aws-eks-iam-auth-controller](https://github.com/rustrial/aws-eks-iam-auth-controller)
  similarly provides a way to map AWS IAM users to the `data.mapUsers` section of the
 `kube-system:aws-auth` ConfigMap.
- [aws-auth](https://github.com/keikoproj/aws-auth) provides a CLI and Golang
  package enabling management of the `data.mapRoles` and `data.mapUsers`
  sections of the `kube-system:aws-auth` ConfigMap.
