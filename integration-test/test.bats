#!/usr/bin/env bats

@test "govvv exists" {
    command -v govvv
}

@test "checks not enough arguments" {
    run govvv
    echo "$output"
    [ "$status" -ne 0 ]
    [[ "$output" == *"not enough arguments"** ]]
}

@test "whitelists certain go commands" {
    run govvv doc
    echo "$output"
    [ "$status" -ne 0 ]
    [[ "$output" == *'only works with "build", "install" and "list". try "go doc" instead'** ]]
}

@test "fails on go tool failure and redirects output" {
    run govvv build -invalid-arg
    echo "$output"
    [ "$status" -ne 0 ]
    [[ "$output" == *'flag provided but not defined: -invalid-arg'** ]]
}

@test "govvv build - dry run" {
    run govvv build -v -print
    echo "$output"
    [ "$status" -eq 0 ]
    [[ "$output" == *"go build \\"* ]]
    [[ "$output" == *"-ldflags"* ]]
}

@test "govvv -flags" {
    run govvv -flags
    echo "$output"
    [ "$status" -eq 0 ]
    [[ "$output" =~ ^-X\ .* ]]
}

@test "govvv list" {
    run govvv list ./integration-test/app-empty
    echo "$output"
    [ "$status" -eq 0 ]
}

@test "govvv build - program with no compile-time variables" {
    tmp="${BATS_TMPDIR}/a.out"
    run govvv build -o "$tmp" ./integration-test/app-empty
    echo "$output"
    [ "$status" -eq 0 ]

    run "$tmp"
    echo "$output"
    [ "$status" -eq 0 ]
    [[ "$output" == "Hello, world!" ]]
}

@test "govvv install - works" {
    run govvv install ./integration-test/app-empty
    echo "$output"
    [ "$status" -eq 0 ]

    run app-empty
    echo "$output"
    [ "$status" -eq 0 ]
    [[ "$output" == "Hello, world!" ]]
}

@test "govvv build - program with compile-time variables" {
    tmp="${BATS_TMPDIR}/a.out"
    run govvv build -o "$tmp" ./integration-test/app-example
    echo "$output"
    [ "$status" -eq 0 ]

    run "$tmp"
    echo "$output"
    [ "$status" -eq 0 ]

    [[ "${lines[0]}" == "Version=untouched" ]]
    [[ "${lines[1]}" == "BuildDate="*Z ]]
    [[ "${lines[2]}" =~ ^GitCommit=[0-9a-f]{4,15}$ ]]
    [[ "${lines[3]}" =~ ^GitBranch=(.*)$ ]]
    [[ "${lines[4]}" =~ ^GitState=(clean|dirty)$ ]]
    [[ "${lines[5]}" =~ ^GitSummary=(.*)$ ]]
}

@test "govvv build - compile-time variables in different package" {
    tmp="${BATS_TMPDIR}/a.out"

    run bash -c "cd ${BATS_TEST_DIRNAME}/app-different-package && govvv build -pkg github.com/ahmetb/govvv/integration-test/app-different-package/mypkg -o $tmp"
    echo "$output"
    [ "$status" -eq 0 ]

    run "$tmp"
    echo "$output"
    [ "$status" -eq 0 ]

    [[ "${lines[0]}" == "Version=2.0.1-app-different-package" ]]
    [[ "${lines[1]}" == "BuildDate="*Z ]]
    [[ "${lines[2]}" =~ ^GitCommit=[0-9a-f]{4,15}$ ]]
    [[ "${lines[3]}" =~ ^GitBranch=(.*)$ ]]
    [[ "${lines[4]}" =~ ^GitState=(clean|dirty)$ ]]
    [[ "${lines[5]}" =~ ^GitSummary=(.*)$ ]]
}

@test "govvv -flags and -pkg" {

    run bash -c "cd ${BATS_TEST_DIRNAME}/app-different-package && govvv -flags -pkg github.com/ahmetb/govvv/integration-test/app-different-package/mypkg"
    echo "$output"
    [ "$status" -eq 0 ]

    [[ "$output" =~ -X\ github.com/ahmetb/govvv/integration-test/app-different-package/mypkg\.Version=2.0.1-app-different-package ]]
}

@test "govvv build - preserves given -ldflags" {
    tmp="${BATS_TMPDIR}/a.out"
    run govvv build -o "$tmp" -ldflags="-X main.MyVariable=myValue" ./integration-test/app-extra-ldflags
    echo "$output"
    [ "$status" -eq 0 ]

    run "$tmp"
    echo "$output"
    [ "$status" -eq 0 ]
    [[ "${lines[0]}" == "MyVariable=myValue" ]]
    [[ "${lines[1]}" =~ ^GitCommit=[0-9a-f]{4,15}$ ]]
}

@test "govvv build - reads Version from ./VERSION file" {
    tmp="${BATS_TMPDIR}/a.out"
    run bash -c "cd ${BATS_TEST_DIRNAME}/app-versioned && govvv build -o ${tmp} ."
    echo "$output"
    [ "$status" -eq 0 ]

    run "$tmp"
    echo "$output"
    [ "$status" -eq 0 ]
    [[ "$output" == "Version=2.0.1-app-versioned" ]]
}

@test "govvv build - reads Version from -version option" {
    tmp="/Users/acabrera/a.out"
    run bash -c "cd ${BATS_TEST_DIRNAME}/app-example && govvv build -o ${tmp} -version 1.2.3-command-line"
    [ "$status" -eq 0 ]
    echo "$output"
    [ "$status" -eq 0 ]

    run "$tmp"
    [ "$status" -eq 0 ]

    [[ "${lines[0]}" == "1.2.3-command-line" ]]
    [[ "${lines[1]}" == "BuildDate="*Z ]]
    [[ "${lines[2]}" =~ ^GitCommit=[0-9a-f]{4,15}$ ]]
    [[ "${lines[3]}" =~ ^GitBranch=(.*)$ ]]
    [[ "${lines[4]}" =~ ^GitState=(clean|dirty)$ ]]
    [[ "${lines[5]}" =~ ^GitSummary=(.*)$ ]]
}

@test "govvv build - ./VERSION file overridden by -version option" {
    tmp="${BATS_TMPDIR}/a.out"
    run bash -c "cd ${BATS_TEST_DIRNAME}/app-versioned && govvv build -o ${tmp} -version 1.2.3-command-line"
    echo "$output"
    [ "$status" -eq 0 ]

    run "$tmp"
    echo "$output"
    [ "$status" -eq 0 ]
    [[ "$output" == "Version=1.2.3-command-line" ]]
}

@test "govvv compiled with govvv" {
    touch main.go
    run govvv install
    echo "$output"
    [ "$status" -eq 0 ]

    run govvv
    echo "$output"
    [[ "${lines[1]}" =~ ^version:\ (.*)@[0-9a-f]{4,15}-(dirty|clean)$ ]]
}
