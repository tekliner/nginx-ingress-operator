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
        stage ('Wait for confirmation of build promotion') {
	    input message: 'Is this build ready for production?', submitter: 'tekliner'
        }
	stage ('Deploy to production') {
            writeFile file: 'k8s/kustomize/kustomization.yaml', text: """
commonLabels:
  version: ${branch}-${build}

images: 
  - name: operator
    newName: 716309063777.dkr.ecr.us-east-1.amazonaws.com/nginx-ingress-operator
    newTag: ${branch}-${build}

resources:
- deployment.yaml
"""
            sh "kubectl kustomize k8s/kustomize/ > deploy.yaml"
            archiveArtifacts: 'deploy.yaml'
            sh "kubectl apply -f deploy.yaml -n default"
        }
    }
}