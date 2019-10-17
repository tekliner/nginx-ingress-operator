node {
    checkout scm
    def branch = env.BRANCH_NAME.toLowerCase()
    def registry = "716309063777.dkr.ecr.us-east-1.amazonaws.com"
    def build = env.BUILD_NUMBER
    def image

    stage("Build image") {
        sh 'docker build -t nginx-ingress-operator .'
        image = docker.image("nginx-ingress-operator")
    }

    stage("Push image") {
        docker.withRegistry("https://"+registry+"/", 'ecr:us-east-1:3c5c323b-afed-4bf0-ae1a-3b19d1c904fe') {
            image.push("${branch}-${build}")
        }
    }

    if (branch == 'master') {

        stage('Generate production manifests') {
            writeFile file: 'k8s/kustomize/kustomization.yaml', text: """
images:
  - name: operator
    newName: 716309063777.dkr.ecr.us-east-1.amazonaws.com/nginx-ingress-operator
    newTag: ${branch}-${build}

resources:
- operator.yaml
- cluster_role.yaml
- cluster_role_binding.yaml
- service_account.yaml

patchesJson6902:
- target:
    group: apps
    version: v1
    kind: Deployment
    name: nginx-ingress-operator
  path: operator-patch.json
- target:
    group: rbac.authorization.k8s.io
    version: v1
    kind: ClusterRoleBinding
    name: nginx-ingress-operator
  path: cluster_role_binding-patch.json
"""
            writeFile file: 'k8s/kustomize/cluster_role_binding-patch.json', text: """
[
  {"op": "replace", "path": "/subjects/0/namespace", "value": "default"}
]
"""
            writeFile file: 'k8s/kustomize/operator-patch.json', text: """
[
  {"op": "replace", "path": "/spec/template/spec/containers/0/env/1/value", "value": ""}
]
"""
            sh "kubectl kustomize k8s/kustomize/ > deploy.yaml"
            archiveArtifacts: 'deploy.yaml'
        }

        stage('Wait for confirmation of build promotion') {
            input message: 'Is this build ready for production?', submitter: 'tekliner'
        }
        stage('Deploy to production') {
            sh "kubectl apply -f deploy/crds/app_v1alpha1_nginxingress_crd.yaml || true"
            sh "kubectl apply -f deploy.yaml -n default"
        }

    } else {
        stage('Generate sandbox manifests') {
            writeFile file: 'k8s/kustomize/kustomization.yaml', text: """
images:
  - name: operator
    newName: 716309063777.dkr.ecr.us-east-1.amazonaws.com/nginx-ingress-operator
    newTag: ${branch}-${build}

resources:
- operator.yaml
- role.yaml
- role_binding.yaml
- service_account.yaml

patchesJson6902:
- target:
    group: apps
    version: v1
    kind: Deployment
    name: nginx-ingress-operator
  path: operator-patch.json
"""
            writeFile file: 'k8s/kustomize/operator-patch.json', text: """
[
  {"op": "replace", "path": "/spec/template/spec/containers/0/env/1/value", "value": "nginx-ingress-operator-${branch}"}
]
"""
            sh "kubectl kustomize k8s/kustomize/ > deploy.yaml"
            archiveArtifacts: 'deploy.yaml'
        }

        stage('Deploy to sandbox') {
            sh "HOME=/root;KUBECONFIG=/root/.kube/sandbox.config kubectl create ns nginx-ingress-operator-${branch} || true"
            sh "HOME=/root;KUBECONFIG=/root/.kube/sandbox.config kubectl apply -f deploy/crds/app_v1alpha1_nginxingress_crd.yaml || true"
            sh "HOME=/root;KUBECONFIG=/root/.kube/sandbox.config kubectl apply -f deploy.yaml -n nginx-ingress-operator-${branch} || true"
        }
    }
}