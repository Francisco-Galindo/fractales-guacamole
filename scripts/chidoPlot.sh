#!/bin/bash

gnuplot << EOF
set term png
set output "$2"
set title "Métricas de rendimiento"
set xlabel "Número de hilos"
set ylabel ""
set yrange[0:]
set datafile separator ","
plot for[col=2:$(awk -F' ' '{print NF}' $1 | head -n 1)] "$1" using 1:col title columnheader(col) with lines
EOF
