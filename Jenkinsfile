def service_name        = "mb-atm-branch-location"
def service_version     = "v2.1.0"
def url_repo            = "https://gitlab.mncbank.co.id/fatahmu/${service_name}.git"
def registry_auth       = "dockerRegistryAuth"
def registry            = "docker.mncbank.co.id:5000"
def kube_credential     = "kubeconfig"
def kube_namespace      = "testing"
def cluster_name        = "rancher-${env.GET_TRIGGER}"
def unitTest_standard   = "0.0%"

pipeline{
    agent any
    stages{
        stage("Unit Test"){
            when{
                expression { env.GET_TRIGGER == "dev" }
            }
            steps{
                script{
                    def golang = tool name: 'go1.16.4', type: 'go'
                    withEnv(["GOROOT=/usr/lib/golang", "PATH+GO=/usr/lib/golang/bin", "GOSUMDB=off"]) {                        
                        unitTest()

                        def unitTestGetValue = sh(returnStdout: true, script: 'go tool cover -func=coverage.out | grep total | sed "s/[[:blank:]]*$//;s/.*[[:blank:]]//"')
                        unitTest_score = "Your score is ${unitTestGetValue}"
                        echo "${unitTest_score}"

                        if (unitTestGetValue >= unitTest_standard){
                            echo "Unit Test fulfill standar value with score ${unitTestGetValue}/${unitTest_standard}"
                        } else {
                            currentBuild.result = 'ABORTED'
                            error("Sorry your unit test score not fulfill standard score ${unitTestGetValue}/${unitTest_standard}")
                        }
                    }
                }   
            }
        }
        stage('Building image') {
            when{
                expression { env.GET_TRIGGER == "dev" }
            }
            steps{
                script {
                    dockerImage = docker.build registry + "/${service_name}" + ":${service_version}"
                }
            }
        }
        stage('Push Image') {
            when{
                expression { env.GET_TRIGGER == "dev" }
            }
            steps{
                script {
                    docker.withRegistry( "https://${registry}", registry_auth ) {
                    dockerImage.push()
                    }
                }
            }
        }
    }
    post{
        success{
            echo "success"
        }
        failure{
            echo "error"
        }
        aborted{
            echo "error"
        }
    }
}

def unitTest() {
    sh "go test ./... -covermode=count -coverprofile coverage.out"
    sh "go tool cover -func=coverage.out"
}
