set -e

GINSIDER_DIR=logicrunner/goplugin/ginsider-cli
GINSIDER_BIN=$GINSIDER_DIR/ginsider-cli
TEST_DATA=testdata
INSOLARD=$TEST_DATA/functional/insolard
CONTRACT_STORAGE=contractstorage

stop_listening()
{
    ports="7777 19191 18181"
    echo "Stop listening..."
    for port in $ports
    do
        lsof -i :$port | grep LISTEN | awk '{print $2}' | xargs kill
    done
}

clear_dirs()
{
    echo "Cleaning directories ..."
    rm -rfv $CONTRACT_STORAGE/*
    rm -rfv data/*
}

create_required_dirs()
{
    mkdir -p $TEST_DATA/functional
    mkdir -p $CONTRACT_STORAGE
    mkdir -p data
}

prepare()
{
    stop_listening
    create_required_dirs
    clear_dirs
}

build_ginsider_cli()
{
    echo "Building ginsider-cli ..."
    go build -o $GINSIDER_BIN $GINSIDER_DIR/main.go
}

build_insolard()
{
    echo "Building insolard ..."
    go build -o $INSOLARD cmd/insolard/*.go 
}

check_binaries()
{
    echo "Check binaries ..."
    if [ ! -f $GINSIDER_BIN ]
    then
        build_ginsider_cli
    fi

    if [ ! -f $INSOLARD ]
    then
        build_insolard
    fi

}

rebuild_binaries()
{
    rm -rfv $GINSIDER_BIN
    rm -rfv $INSOLARD
    check_binaries
}

check_working_dir()
{
    if ! pwd | grep -q "go/src/github.com/insolar/insolar$"
    then
        echo "Run me from insolar root"
        exit 1
    fi
}

usage()
{
    echo "usage: $0 <clear,rebuild>"
}

process_input_params()
{
    param=$1
    if [  "$param" == "clear" ]
    then
        prepare
        exit 0
    fi

    if [ "$param" == "rebuild" ]
    then
        rebuild_binaries
        exit 0
    fi

    if [ "$param" == "help" ] || [ "$param" == "-h" ] || [ "$param" == "--help" ]
    then
        usage
        exit 0
    fi
}

param=$1
check_working_dir
process_input_params $param

prepare
check_binaries

$INSOLARD --config scripts/insolard/insolar.yaml

