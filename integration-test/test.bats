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

@test "only works for build command" {
    run govvv doc
    echo "$output"
    [ "$status" -ne 0 ]
    [[ "$output" == *'only works with "build". try "go doc" instead'** ]] 
}

@test "fails on go tool failure and redirects output" {
    run govvv build -invalid-arg
    echo "$output"
    [ "$status" -ne 0 ]
    [[ "$output" == *'flag provided but not defined: -invalid-arg'** ]] 
}

@test "compiles program not without govvv variables" {
    tmp="${BATS_TMPDIR}/a.out"
    run govvv build -o "$tmp" ./integration-test/app-empty  
    echo "$output"
    [ "$status" -eq 0 ]

    run "$tmp"
    echo "$output"
    [ "$status" -eq 0 ]
    [[ "$output" == "Hello, world!" ]]
}

@test "compiles program using the govvv variables" {
    tmp="${BATS_TMPDIR}/a.out"
    run govvv build -o "$tmp" ./integration-test/app-example  
    echo "$output"
    [ "$status" -eq 0 ]

    run "$tmp"
    echo "$output"
    [ "$status" -eq 0 ]

    [[ "${lines[0]}" == "BuildDate="*Z ]]
    [[ "${lines[1]}" =~ ^GitCommit=[0-9a-f]{7}$ ]]
    [[ "${lines[2]}" =~ ^GitBranch=(.*)$ ]]
    [[ "${lines[3]}" =~ ^GitState=(clean|dirty)$ ]]
}


@test "existing -ldflags are preserved" {
    tmp="${BATS_TMPDIR}/a.out"
    run govvv build -o "$tmp" -ldflags="-X main.MyVariable=myValue" ./integration-test/app-extra-ldflags
    echo "$output"
    [ "$status" -eq 0 ]

    run "$tmp"
    echo "$output"
    [ "$status" -eq 0 ]
    [[ "${lines[0]}" == "MyVariable=myValue" ]]
    [[ "${lines[1]}" =~ ^GitCommit=[0-9a-f]{7}$ ]]
}

@test "Version is read from ./VERSION file" {
    tmp="${BATS_TMPDIR}/a.out"
    run bash -c "cd ${BATS_TEST_DIRNAME}/app-versioned && govvv build -o ${tmp} ."
    echo "$output"
    [ "$status" -eq 0 ]

    run "$tmp"
    echo "$output"
    [ "$status" -eq 0 ]
    [[ "$output" == "Version=2.0.1" ]]
}