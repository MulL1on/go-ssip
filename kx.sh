#!/bin/bash

THRIFT_FILE="$1"


kitex -module go-ssip -gen-path app/common/kitex_gen $THRIFT_FILE
