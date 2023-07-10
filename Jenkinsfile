#!/usr/bin/env groovy
def GO111MODULE="on"
def GOPROXY="https://goproxy.cn,direct"
def POL_ENV="dev"
def GONOPROXY="*.bjx.cloud"
def GONOSUMDB="*.bjx.cloud"
def GOSUMDB="off"
def GOPRIVATE="*.bjx.cloud"
pipeline{
	agent {
		label 'jenkins-slave'
	}
    options {
        disableConcurrentBuilds()
        skipDefaultCheckout()
        timeout(time: 1, unit: 'HOURS')
        timestamps()
    }
	parameters {
		choice(choices: ['roll-update', 'roll-back'], description: '版本发布或者回滚', name: 'action')
		choice(choices: ['all', 'polaris-app', 'polaris-callsvc', 'polaris-commonsvc', 'polaris-front-proxy-inside', 'polaris-front-proxy-outside', 'polaris-idsvc', 'polaris-msgsvc', 'polaris-orgsvc', 'polaris-processsvc', 'polaris-projectsvc', 'polaris-resourcesvc', 'polaris-rolesvc', 'polaris-schedule', 'polaris-trendssvc', 'polaris-websitesvc'], description: '发布或回滚的服务', name: 'service')
		string(name: 'rollBackImage', defaultValue: '例:v1.4.3.43.release', description: '回滚的版本号')
		choice(choices: ['test', 'test1'], description: '部署或回滚的环境', name: 'deployEnv')
		booleanParam(name: 'runSonar', defaultValue: false, description: '是否执行sonar-scanner代码扫描')
		booleanParam(name: 'runBuild', defaultValue: true, description: '是否执行构建')
		booleanParam(name: 'runDeploy', defaultValue: true, description: '是否需要部署')
		//booleanParam(name: 'runTest', defaultValue: false, description: '是否执行自动化测试')
	}
	environment{
		GIT_URL = "https://github.com/star-table/polaris-backend.git"
		SONAR_APP = "polaris-backend"
		SONAR_PASSWD = "password=admin"
		CREDENTIALS_ID = "989050c7-b02d-4f9f-a4b1-d081667e4b56"
		ACTION_TO = sh(returnStdout: true,script: 'echo ${action}').trim()
		SERVICE_TO = sh(returnStdout: true,script: 'echo ${service}').trim()
		ROLL_BACK_IMAGE = sh(returnStdout: true,script: 'echo ${rollBackImage}').trim()
		DEPLOY_ENV = sh(returnStdout: true,script: 'echo ${deployEnv}').trim()
		RUN_SONAR = sh(returnStdout: true,script: 'echo ${runSonar}').trim()
		RUN_BUILD = sh(returnStdout: true,script: 'echo ${runBuild}').trim()
		RUN_DEPLOY = sh(returnStdout: true,script: 'echo ${runDeploy}').trim()
		//RUN_TEST = sh(returnStdout: true,script: 'echo ${runTest}').trim()
	}
	stages {
		stage('roll back') {
			when {
				environment name:'ACTION_TO',value:'roll-back'
			}
			steps {
				sh '''
				if [ ${SERVICE_TO} = "all" ];then
				svcs=$(ls /data/package/polaris-backend/release/${ROLL_BACK_IMAGE})
				for svccm in ${svcs}
				do
				kubectl -n${DEPLOY_ENV} apply -f /data/package/polaris-backend/release/${ROLL_BACK_IMAGE}/${svccm}
				svc=$(echo ${svccm}|awk -F '-configmap' '{print $1}')
				svcimage="registry-vpc.cn-shanghai.aliyuncs.com/polaris-team/"${svc}":"${ROLL_BACK_IMAGE}
				kubectl -n${DEPLOY_ENV} set image deployment/${svc} ${svc}=${svcimage} --record
				done
				else
				kubectl -n${DEPLOY_ENV} apply -f /data/package/polaris-backend/release/${ROLL_BACK_IMAGE}/${SERVICE_TO}-configmap.yaml
				svcimage="registry-vpc.cn-shanghai.aliyuncs.com/polaris-team/"${SERVICE_TO}":"${ROLL_BACK_IMAGE}
				kubectl -n${DEPLOY_ENV} set image deployment/${SERVICE_TO} ${SERVICE_TO}=${svcimage} --record
				fi
				'''
			}
		}
		stage('git checkout') {
			when {
				environment name:'ACTION_TO',value:'roll-update'
			}
			steps {
				git branch: "${BRANCH_NAME}", credentialsId: "${CREDENTIALS_ID}", url: "${GIT_URL}"
			}
		}
		stage('run sonar') {
			when {
				environment name:'ACTION_TO',value:'roll-update'
				environment name:'RUN_SONAR',value:'true'
			}
			steps {
				sh '''
				BRANCH_TAG=$(echo ${BRANCH_NAME} | awk -F '/' '{print $1"-"$2}')
				sed -i "s/${SONAR_APP}/${SONAR_APP}-${BRANCH_TAG}/g" sonar-project.properties
				sed -i "s/${SONAR_PASSWD}/password=runx@123/g" sonar-project.properties
				cat sonar-project.properties
				sonar-scanner
				'''
			}
		}
		stage('go build') {
			when {
				environment name:'ACTION_TO',value:'roll-update'
				environment name:'RUN_BUILD',value:'true'
			}
			steps {
				sh '''
				BRANCH_TAG=$(echo ${BRANCH_NAME} | awk -F '/' '{print $NF}')
				/data/package/polaris-backend/install.sh ${BRANCH_NAME} ${BRANCH_TAG}.${BUILD_ID} dev
				/data/package/polaris-backend/release/scripts/getconfigfile.sh
				'''
			}
		}
		stage('deploy') {
			when {
				environment name:'ACTION_TO',value:'roll-update'
				environment name:'RUN_BUILD',value:'true'
				environment name:'RUN_DEPLOY',value:'true'
			}
            //input {
            //    message "确定发布?"
            //    ok "Yes, continue."
            //}
			steps {
				sh '''
				DEPLOY_RELEASE=$(echo ${BRANCH_NAME} | awk -F '/' '{print $NF}').${BUILD_ID}.$(echo ${BRANCH_NAME} | awk -F '/' '{print $1}')
				if [ ${SERVICE_TO} = "all" ];then
				svcs=$(ls /data/package/polaris-backend/release/${DEPLOY_RELEASE})
				for svccm in ${svcs}
				do
				kubectl -n${DEPLOY_ENV} apply -f /data/package/polaris-backend/release/${DEPLOY_RELEASE}/${svccm}
				svc=$(echo ${svccm}|awk -F '-configmap' '{print $1}')
				svcimage="registry-vpc.cn-shanghai.aliyuncs.com/polaris-team/"${svc}":"${DEPLOY_RELEASE}
				kubectl -n${DEPLOY_ENV} set image deployment/${svc} ${svc}=${svcimage} --record
				done
				else
				kubectl -n${DEPLOY_ENV} apply -f /data/package/polaris-backend/release/${DEPLOY_RELEASE}/${SERVICE_TO}-configmap.yaml
				svcimage="registry-vpc.cn-shanghai.aliyuncs.com/polaris-team/"${SERVICE_TO}":"${DEPLOY_RELEASE}
				kubectl -n${DEPLOY_ENV} set image deployment/${SERVICE_TO} ${SERVICE_TO}=${svcimage} --record
				fi
				'''
			}
		}
		//stage('automated testing') {
		//	when {
		//		environment name:'ACTION_TO',value:'roll-update'
		//		environment name:'RUN_BUILD',value:'true'
		//		environment name:'RUN_TEST',value:'true'
		//	}
		//	steps {
		//		sh "python -V"
		//	}
		//}
	}
}
