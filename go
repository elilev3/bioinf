#!/bin/bash

PIP_VERSION='1.3.1'
PROJECT=bioinf

cd $(dirname ${0})

CELLAR=/usr/local/Cellar
XCODE_PATH=/Applications/XCode.app/Contents/Developer
if [ "Darwin" == `uname` ]; then
    if [ ! -x /usr/bin/clang ]; then
        echo "You need to install the XCode Command Line Tools before you can continue."
        exit -1
    fi
    if [ ! -x /usr/bin/xcode-select ]; then
        echo "You need to install the Apple Developer Tools before you can continue."
        exit -1
    elif [ "${XCODE_PATH}" != `/usr/bin/xcode-select -print-path | tr -d '\n'` ]; then
        echo "We need to switch up your xcode, fool. This requires system permissions to proceed..."
        sudo /usr/bin/xcode-select -switch ${XCODE_PATH}
    fi

    VIRTUALENV=$(which virtualenv)
    if [ "$VIRTUALENV" == "" ]; then
        echo "VirtualEnv requires system permissions to install..."
        sudo easy_install distribute
        sudo easy_install virtualenv
    fi

    if [ ! -d env ]; then
        virtualenv --clear --verbose --prompt "[$PROJECT] " --distribute -p python2.7 env
    fi
fi

echo "Activating virtualenv"
. env/bin/activate

mkdir -p data

echo Installing any necessary packages into your virtualenv...
PIP_VERSION_INSTALLED=$(pip --version | sed 's/ from.*$//g' | sed 's/pip //')
if [ "${PIP_VERSION}" \> "${PIP_VERSION_INSTALLED}" ]; then
     pip install --upgrade pip==${PIP_VERSION}
fi
if [ -s conf/dependencies.conf ]; then
    pip install --requirement=conf/dependencies.conf
fi

if [ setup.py -nt exterminator.egg-info ]; then
     python setup.py develop
fi

case $1 in
    'test' )
        if [ -s conf/dependencies-test.conf ]; then
            pip install --requirement=conf/dependencies-test.conf
        fi
        env/bin/nosetests --with-coverage --cover-package bioinf --with-spec --spec-color	
    ;;

    * )
        echo "Run something!"
    ;;
esac
