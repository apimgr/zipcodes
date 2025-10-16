pipeline {
    agent none
    
    environment {
        PROJECTNAME = 'zipcodes'
        PROJECTORG = 'apimgr'
        REGISTRY = 'ghcr.io'
        VERSION = "${env.BRANCH_NAME == 'main' ? readFile('release.txt').trim() : 'dev'}"
    }
    
    stages {
        stage('Build') {
            parallel {
                stage('Build AMD64') {
                    agent {
                        label 'amd64'
                    }
                    steps {
                        echo "Building ${PROJECTNAME} ${VERSION} for AMD64..."
                        sh 'make build'
                        stash includes: 'binaries/**', name: 'binaries-amd64'
                    }
                }
                
                stage('Build ARM64') {
                    agent {
                        label 'arm64'
                    }
                    steps {
                        echo "Building ${PROJECTNAME} ${VERSION} for ARM64..."
                        sh 'make build'
                        stash includes: 'binaries/**', name: 'binaries-arm64'
                    }
                }
            }
        }
        
        stage('Test') {
            parallel {
                stage('Test AMD64') {
                    agent {
                        label 'amd64'
                    }
                    steps {
                        echo "Running tests on AMD64..."
                        sh 'make test'
                    }
                }
                
                stage('Test ARM64') {
                    agent {
                        label 'arm64'
                    }
                    steps {
                        echo "Running tests on ARM64..."
                        sh 'make test'
                    }
                }
            }
        }
        
        stage('Docker Build & Push') {
            when {
                branch 'main'
            }
            agent {
                label 'amd64'
            }
            steps {
                echo "Building and pushing Docker images..."
                script {
                    withCredentials([usernamePassword(credentialsId: 'ghcr-token', usernameVariable: 'DOCKER_USER', passwordVariable: 'DOCKER_PASS')]) {
                        sh '''
                            echo $DOCKER_PASS | docker login ghcr.io -u $DOCKER_USER --password-stdin
                            make docker
                        '''
                    }
                }
            }
        }
        
        stage('Release') {
            when {
                branch 'main'
            }
            agent {
                label 'amd64'
            }
            steps {
                echo "Creating GitHub release ${VERSION}..."
                unstash 'binaries-amd64'
                unstash 'binaries-arm64'
                script {
                    withCredentials([string(credentialsId: 'github-token', variable: 'GITHUB_TOKEN')]) {
                        sh '''
                            export GITHUB_TOKEN=${GITHUB_TOKEN}
                            make release
                        '''
                    }
                }
            }
        }
    }
    
    post {
        success {
            echo "✓ Pipeline completed successfully!"
        }
        failure {
            echo "✗ Pipeline failed!"
        }
        cleanup {
            cleanWs()
        }
    }
}
