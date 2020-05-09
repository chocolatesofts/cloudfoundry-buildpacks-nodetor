# Cloud Foundry NodeTor build pack 
### based on Node.js Buildpack (https://github.com/cloudfoundry/nodejs-buildpack)

A Cloud Foundry [buildpack](http://docs.cloudfoundry.org/buildpacks/) for Node based apps with added Tor support
Supplement for libraries like [tor-request](https://www.npmjs.com/package/tor-request)

### Buildpack User Documentation

Official buildpack documentation can be found at [node buildpack docs](http://docs.cloudfoundry.org/buildpacks/node/index.html).

### Building the Buildpack

To build this buildpack, run the following commands from the buildpack's directory:

1. Source the .envrc file in the buildpack directory.

   ```bash
   source .envrc
   ```
   To simplify the process in the future, install [direnv](https://direnv.net/) which will automatically source .envrc when you change directories.

1. Install buildpack-packager

    ```bash
     go install github.com/cloudfoundry/libbuildpack/packager/buildpack-packager
    ```

1. Build the buildpack

    ```bash
    buildpack-packager build -stack [STACK] [ --cached=(true|false) ]
    ```

1. Use in Cloud Foundry

   Upload the buildpack to your Cloud Foundry and optionally specify it by name

    ```bash
    cf create-buildpack [BUILDPACK_NAME] [BUILDPACK_ZIP_FILE_PATH] 1
    cf push my_app [-b BUILDPACK_NAME]
    ```
### Contributing

Feel free to contribute here.

### Reporting Issues

Open an issue on this project.

### Acknowledgements

Inspired by the [Heroku buildpack](https://github.com/heroku/heroku-buildpack-nodejs).
