#!/usr/bin/env bash

declare modules=('ahk' 'aqi' 'automation' 'bravia.tv' 'calendar' 'conbee.lights' 'conbee.sensors' 'config/config' 'config/strcrypt' 'config/pubconf' 'flux' 'presence' 'samsung.tv' 'shout' 'suncalc' 'weather' 'wemo' 'yee')

export GO_HOME_KEY=2D4B6150645267552D4B615064526755

goversion=$(go version)
printf "$goversion\n"

gohomepath=$(pwd)
for module in "${modules[@]}" ; do 
    printf "Building: $module \n"

    # Compile each module
    cd "$module"
    go build
    cd "$gohomepath"
done

