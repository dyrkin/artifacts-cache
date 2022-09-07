# Artifacts cache

1. Start server passing a data dir with `./repository <data_dir>`
2. Copy `artifacts` binary on the server where `gitlab-runner` is installed
3. For shell executor just add `artifacts` binary to the `PATH` env variable
4. For docker executor pass `artifacts` binary to every docker container:
   ```toml
   [runners.docker]
     volumes = ["/usr/local/bin/artifacts:/usr/local/bin/artifacts"]
   ```
5. Define `ARTIFACTS_SUBSET_ID` and `ARTIFACTS_REPOSITORIES` variables in a build pipeline
    ```yaml
    variables:
        ARTIFACTS_SUBSET_ID: "$CI_PROJECT_PATH-$CI_PIPELINE_ID"
        ARTIFACTS_REPOSITORIES: "http://repository1:8080, http://repository2:8080, http://repository3:8080"
    ```
6. Use
    ```yaml
    stages:
      - push
      - pull
        
        
    push:
      image: golang:1.19
      stage: push
      tags:
        - docker
      script:
        - artifacts push 'files/*.txt'
        
    pull:
      image: golang:1.19
      stage: pull
      tags:
        - docker
      script:
        - artifacts pull 'files/*.txt'
    ```

