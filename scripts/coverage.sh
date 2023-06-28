#!/bin/bash
#
# Code coverage generation

COVERAGE_DIR="${COVERAGE_DIR:-coverage}"
COVERAGE_PKG_DIR="${COVERAGE_PKG_DIR:-coverage/pkg}"
PKG_LIST=$(go list ./... | grep -v /vendor/)

# Create the coverage files directory
mkdir -p "$COVERAGE_PKG_DIR";

# Create a coverage file for each package
for package in ${PKG_LIST}; do
    go test \
        -covermode=count \
        -coverprofile "${COVERAGE_PKG_DIR}/${package##*/}.cov" \
        "$package" ;
done ;

# Merge the coverage profile files
echo 'mode: count' > "${COVERAGE_DIR}"/coverage.cov ;
tail -q -n +2 "${COVERAGE_PKG_DIR}"/*.cov >> "${COVERAGE_DIR}"/coverage.cov ;

# Display the global code coverage
go tool cover -func="${COVERAGE_DIR}"/coverage.cov ;

# If needed, generate HTML report
if [ "$1" == "html" ]; then
    go tool cover -html="${COVERAGE_DIR}"/coverage.cov -o "${COVERAGE_DIR}"/coverage.html ;
fi
