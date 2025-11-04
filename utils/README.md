# Utils

This folder contains shared code and utility files. This also contains the modified go version that is instrumented to gve timing data.

## Creating the instance of GO

In order to create the instance of go with all the instrumentation changes you need to do these steps.

1. cd into go-instrumented/src
2. Run ./make.bash
3. For now, to use the local version of go you have to call it using ./utils/go-instrumented/bin/go and then any relevant functions eg. ./utils/go-instrumented/bin/go verison == go verison


Before running any programs ensure that the go submodule (local version of go) is up to date with its github counterpart:

git submodule update --init --recursive