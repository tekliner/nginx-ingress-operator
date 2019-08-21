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

    stage('Generate manifests') {
        writeFile file: 'k8s/kustomize/kustomization.yaml', text: """
images: 
  - name: operator
    newName: 716309063777.dkr.ecr.us-east-1.amazonaws.com/nginx-ingress-operator
    newTag: ${branch}-${build}

resources:
- deployment.yaml
"""
        sh "kubectl kustomize k8s/kustomize/ > deploy.yaml"
        archiveArtifacts: 'deploy.yaml'
    }

    if (branch == 'master') {
        stage ('Wait for confirmation of build promotion') {
            input message: 'Is this build ready for production?', submitter: 'tekliner'
        }
	stage ('Deploy to production') {
            sh "kubectl apply -f deploy.yaml -n default"
        }
    } else {
            sh "HOME=/root;KUBECONFIG=/root/.kube/sandbox.config kubectl create ns nginx-ingress-operator-${branch} || true"
            sh "HOME=/root;KUBECONFIG=/root/.kube/sandbox.config kubectl apply -f deploy/service_account.yaml -n nginx-ingress-operator-${branch} || true"
            sh "HOME=/root;KUBECONFIG=/root/.kube/sandbox.config kubectl apply -f deploy/role_binding.yaml -n nginx-ingress-operator-${branch} || true"
            sh "HOME=/root;KUBECONFIG=/root/.kube/sandbox.config kubectl apply -f deploy/role.yaml -n nginx-ingress-operator-${branch} || true"
            sh "HOME=/root;KUBECONFIG=/root/.kube/sandbox.config kubectl apply -f deploy.yaml -n nginx-ingress-operator-${branch} || true"
    }
}