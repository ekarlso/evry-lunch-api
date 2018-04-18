node {
  ws {
    try {
      stage("scm") {
      // This is to allow for easier import paths in go.
        sh 'git config --global url."ssh://evry@vs-ssh.visualstudio.com:22/DefaultCollection/evry-iot/_ssh/".insteadOf "ssh://git@evry.visualstudio.com/reception/"'
        git 'ssh://git@evry.visualstudio.com/reception/visitors-api'
      }

      def registry = "evrybots.azurecr.io"
      def imgName = "evry-reception-api"

      // https://github.com/deis/workflow-cli/blob/master/Jenkinsfile
      gitCommit = sh(returnStdout: true, script: 'git rev-parse HEAD').trim()
      gitShortCommit = gitCommit[-8..-1]

      def imgTag = "${env.BUILD_NUMBER}-${gitShortCommit}"
      def imgFullName = "${registry}/${imgName}:${imgTag}"

      def goImage = "golang:1.10"
      def opaImage = "openpolicyagent/opa:0.7.1-alpine"

      stage("dependencies") {
        def workspace = pwd()

        docker.image(goImage).inside("-v ${workspace}:/src -w /src") {
          sh 'export GOCACHE=/src/.GOCACHE; V=1 make setup'
          sh 'export GOCACHE=/src/.GOCACHE; make deps'
        }
      }

      stage("test") {
        docker.image("postgres").withRun("-p 5432:5432 -e POSTGRES_PASSWORD=postgres") { c -> 
          docker.image(goImage).inside("-v ${workspace}:/src -w /src --link ${c.id}:db") {
            sh "export GOCACHE=/src/.GOCACHE; export DB_CONNECTION='host=db port=5432 dbname=postgres user=postgres password=postgres sslmode=disable': V=1 make test"
          }
        }

        docker.image(opaImage).inside("-v ${workspace}:/src -w /src --entrypoint ''") {
          sh "/opa test -v policies"
        }
      }

      stage("build") {
        docker.image(goImage).inside("-v ${workspace}:/src -w /src") {
          sh 'export GOCACHE=/src/.GOCACHE; V=1 make build'
        }
      }

      stage("publish") {
        def buildArgs = "."
        def img = docker.build(imgFullName, buildArgs)

        docker.withRegistry("https://${registry}", "registry-${registry}") {
          img.push()
          img.push('latest')
        }
      }

      stage("deploy") {
        def helmUrl = "https://storage.googleapis.com/kubernetes-helm"
        def helmVersion = "2.8.2"
        def helmTarball = "helm-v${helmVersion}-linux-amd64.tar.gz"

        sh("curl -O ${helmUrl}/${helmTarball}")
        sh("tar -zxvf ${helmTarball} -C $HOME && rm ${helmTarball} && mv /var/jenkins_home/linux-amd64/helm ~")
        switch(env.BRANCH_NAME) {
          case 'master':
            def clusterName = "svg-prod-cluster";
            sh("KUBECONFIG=/var/jenkins_home/.kube/config-${clusterName} ~/istioctl kube-inject --includeIPRanges=10.244.0.0/16,10.240.0.0/16 -f chart/templates/deployment.yaml > chart/templates/deployment2.yaml && rm chart/templates/deployment.yaml")
            sh("KUBECONFIG=/var/jenkins_home/.kube/config-${clusterName} ~/helm upgrade --set 'image.repository=${registry}/${imgName},image.tag=${imgTag}' --install ${imgName} --namespace=evry-reception chart")
          break
        }
      }
      
    } catch (InterruptedException e) {
      throw e

    // Catch all build failures and report them to Slack etc here.
    } catch (e) {
      throw e

    // Clean up the workspace before exiting. Wait for Jenkins' asynchronous
    // resource disposer to pick up before we close the connection to the worker
    // node.
    } finally {
      step([$class: 'WsCleanup'])
      sleep 10
    }
  }
}
