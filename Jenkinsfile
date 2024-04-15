#!groovy

pipeline {
    agent { label 'dockerd' }

    stages {
        stage('Build') {
            steps {
                dockerBuildTagPush()
            }
        }
    }
}
