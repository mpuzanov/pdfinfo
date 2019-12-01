#!/bin/bash
# подготавливаем структуру каталогов для программ

if [[ -z $1 ]]; then BINARY_DIR="bin"; else BINARY_DIR=$1; fi 

PATHWIN="${BINARY_DIR}/win"
PATHLINUX="${BINARY_DIR}/linux"

for dir in $PATHWIN $PATHLINUX
do
    if [[ !(-e $dir) ]] # проверяем на существование
    then
        echo "$dir not found. Create"
        mkdir -p $dir
        if [[ -e $dir ]]; then { echo "$dir - created."; } fi    
    fi
done
