pipeline {
    agent any

    environment {
        GIT_REPO = 'https://github.com/DukeNicasio/cloud_final.git'
        COMPOSE_FILE = 'docker-compose.yml'
    }

    stages {
        stage('Checkout') {
            steps {
                git url: "${GIT_REPO}", branch: 'main'
            }
        }

        stage('Shutdown Existing Containers') {
            steps {
                sh "docker-compose -f ${COMPOSE_FILE} down"
            }
        }

        stage('Build Services') {
            steps {
                sh "docker-compose -f ${COMPOSE_FILE} build"
            }
        }

        stage('Restart Services') {
            steps {
                sh "docker-compose -f ${COMPOSE_FILE} up -d"
            }
        }
    }

    post {
        always {
            echo 'Pipeline finished.'
        }
    }
}
