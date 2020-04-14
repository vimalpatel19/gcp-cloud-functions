# GCP (Google Cloud Platform) - Cloud Functions 

This repository contains various Cloud Functions (written in `Go`) used across various projects that I am developing on my [GitHub account](https://github.com/vimalpatel19).

### Event/Trigger types
Currently, the following types of functions are found in this repository:
- HTTP trigger

### Testing locally
The functions created under the `api` directory can be run locally by running the main method found at `api/cmd/main.go`. This method runs a generic HTTP server on port `8080`. See that file to identify what path to hit/access to run the intended Cloud Function. As each new function is created, a new path should be added in this file to perform local testing against that function. 

### Deploying
- HTTP tigger functions can be executed by running the following CLI command:
```
gcloud functions deploy <FUNCTION_NAME> --entry-point <METHOD_NAME> --runtime go113 --trigger-http --set-env-vars <ENV1=VALUE1,ENV2=VALUE2,...>
```

### TODO
- [ ] Add documentation to README for deploying each Cloud Function
- [ ] Refactor code for resusability of common code across multiple functions