#!/bin/sh
./gin -c config/b00m.config -stderrthreshold=INFO 2>&1 | funnel -app=b00m-gin &


