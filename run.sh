#!/bin/bash
ffmpeg -i ./test/test.mp4 -i ./test/back.jpg -filter_complex "[0:v]chromakey=0x00FF00:0.1:0.2[ckout];[1:v][ckout]overlay[out]" -map "[out]" output.mp4
