[![](https://images.microbadger.com/badges/version/etejeda/github-akamai-purge.svg)](https://microbadger.com/images/etejeda/github-akamai-purge "Get your own version badge on microbadger.com")
![Docker Pulls](https://img.shields.io/docker/pulls/etejeda/github-akamai-purge)
![GitHub top language](https://img.shields.io/github/languages/top/enriquetejeda/github-akamai-purge)
![GitHub last commit](https://img.shields.io/github/last-commit/enriquetejeda/github-akamai-purge)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0)
# Github Akamai Purge
A solution for make fast purge in the akamai network all made with Go in a distroless container (only 12MB~)! :heart:

Visit also [docker hub repository](https://hub.docker.com/repository/docker/etejeda/github-akamai-purge).

## How works?

Basically is a golang app, which have integration with the api client from Akamai Network & Github.

The flow is simple, if the commit id (sha) is passed to the application search this commit exactly and if doesnt pass the id the last commit on the repo is picked, for later filter the assets and sends the array to perform the purge to akamai api.


<div>
    <img src="https://github.com/enriquetejeda/github-akamai-purge/raw/master/docs/mainflow.png" width="250" />
    <div>
        <em>Fig. 1: The Main Flow</em>
    </div>
</div>


## Requirements

* Docker Engine. :heart:
* Valid Github Token with this permissions. [(more info)](https://help.github.com/es/github/authenticating-to-github/creating-a-personal-access-token-for-the-command-line)
* Valid API Credentials for make api calls over the Akamai Network. [(more info)](https://developer.akamai.com/api/getting-started)

## Getting Started

### Standalone

You only run this command in your terminal:

```
docker run 
-e 'AKAMAI_HOST=' 
-e 'AKAMAI_ACCESS_TOKEN=****' 
-e 'AKAMAI_CLIENT_TOKEN=****' 
-e 'AKAMAI_CLIENT_SECRET=****'
-e 'AKAMAI_PURGE_METHOD=invalidate' 
-e 'AKAMAI_PURGE_NETWORK=production' 
-e 'AKAMAI_PURGE_HOSTNAME=www.foo.com' 
-e 'GITHUB_TOKEN=****' 
-e 'GITHUB_ORGANIZATION=EnriqueTejeda' 
-e 'GITHUB_REPOSITORY=FOO-REPO' 
-e 'GITHUB_BRANCH=master' etejeda/github-akamai-purge:latest
```

or just create a file naming `.env` and put the env vars inside:

```
AKAMAI_HOST=
AKAMAI_ACCESS_TOKEN=
AKAMAI_CLIENT_TOKEN=
AKAMAI_CLIENT_SECRET=
AKAMAI_PURGE_METHOD=invalidate
AKAMAI_PURGE_NETWORK=production
AKAMAI_PURGE_HOSTNAME=

GITHUB_COMMIT_SHA=
GITHUB_TOKEN=
GITHUB_ORGANIZATION=
GITHUB_REPOSITORY=
GITHUB_BRANCH=master
```

and then run the container in your terminal:

```
docker run --env-file YOUR-ENV-FILE etejeda/github-akamai-purge:latest
```
### Continuous Integration

#### Jenkins

For more security, first add the tokens like secrets (like secret text) and then you need add a step in your pipeline executing the container image, example:

```
#!groovy
pipeline {
    agent { node { label 'master' } }
    options { skipDefaultCheckout true }
    environment {}
    stages {
        stage('Build'){
            steps {
                checkout scm
            }
        }
        stage('Test'){
            steps {
                echo 'testing'
            }
        }
        stage('Deploy'){
            steps {
                script {
                    /*
                    * Before you need add all tokens like a secret for more security
                    */
                    withCredentials([
                            string(credentialsId: 'akamai-access-token', variable: 'akamai_access_token'),
                            string(credentialsId: 'akamai-client-token', variable: 'akamai_client_token'),
                            string(credentialsId: 'akamai-client-secret', variable: 'akamai_client_secret'),    
                            string(credentialsId: 'github-token', variable: 'github_token')
                    ]){
                        def command = "docker run " +
                        "-e AKAMAI_HOST=${env.AKAMAI_HOST} " +
                        "-e AKAMAI_ACCESS_TOKEN=${akamai_access_token} " +
                        "-e AKAMAI_CLIENT_TOKEN=${akamai_client_token} " +
                        "-e AKAMAI_CLIENT_SECRET=${akamai_client_secret} " +
                        "-e AKAMAI_PURGE_METHOD=invalidate " +
                        "-e AKAMAI_PURGE_NETWORK=production " +
                        "-e AKAMAI_PURGE_HOSTNAME=${env.AKAMAI_PURGE_HOSTNAME} " +
                        "-e GITHUB_TOKEN=${github_token} " +
                        "-e GITHUB_ORGANIZATION=ExperienciasXcaret " +
                        "-e GITHUB_REPOSITORY=${env.GITHUB_REPOSITORY} " +
                        "-e GITHUB_BRANCH=master "
                        if(env.GITHUB_COMMIT){
                            command += " -e GITHUB_COMMIT_SHA=${env.GIT_COMMIT} "
                        }
                
                        command += " etejeda/github-akamai-purge:latest"
                
                        sh command
                    } 
                }
            }
        }
    }   
}
```

## Development
### Building the container

I provided a makefile for do this job, only run this command:
```
make run build 
```

### Environment Variables 

For run this tool, you need a valid *Github Token* and valid API Access for interact with the akamai api, for more information with this you can visit the akamai official documentation [here](https://developer.akamai.com/api/getting-started).

| Name  | Description  | Default | Required |
| -- | -- | -- | -- |
| AKAMAI_HOST | The akamai base url for you credential | - | *yes* |
| AKAMAI_ACCESS_TOKEN | Access token for interact with the api | - | *yes* |
| AKAMAI_CLIENT_TOKEN | Client token provided when your create a api client | - | *yes* |
| AKAMAI_CLIENT_SECRET | Client secret token provided when your create a api client | - | *yes* |
| AKAMAI_PURGE_METHOD |  Purge method, `invalidate` or `delete` | `invalidate` | no |
| AKAMAI_PURGE_NETWORK | The current network for refresh the content, `staging` or `production` | `production` | no |
| AKAMAI_PURGE_HOSTNAME | Hostname, example: www.google.com | - | *yes* |
| GITHUB_COMMIT_SHA | Git SHA if you like point a specified commit | - | no |
| GITHUB_TOKEN | Github token with permissions for read commit in the repo | - | *yes*  |
| GITHUB_ORGANIZATION|  Github Org or User, example `EnriqueTejeda` | - | *yes*  |
| GITHUB_REPOSITORY | Name of the repo on Github | - | *yes*  |
| GITHUB_BRANCH | The branch for looking the commit, example `master`, `staging`, `develop` | `master` | no |

## How contribute? :rocket:

Please feel free to contribute to this project, please fork the repository and make a pull request!. :heart:

## Share the Love :heart:

Like this project? Please give it a ★ on [this GitHub](https://github.com/EnriqueTejeda/github-akamai-purge)! (it helps me a lot).

## License

[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)](https://opensource.org/licenses/Apache-2.0) 

See [LICENSE](LICENSE) for full details.

    Licensed to the Apache Software Foundation (ASF) under one
    or more contributor license agreements.  See the NOTICE file
    distributed with this work for additional information
    regarding copyright ownership.  The ASF licenses this file
    to you under the Apache License, Version 2.0 (the
    "License"); you may not use this file except in compliance
    with the License.  You may obtain a copy of the License at

      https://www.apache.org/licenses/LICENSE-2.0

    Unless required by applicable law or agreed to in writing,
    software distributed under the License is distributed on an
    "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
    KIND, either express or implied.  See the License for the
    specific language governing permissions and limitations
    under the License.

