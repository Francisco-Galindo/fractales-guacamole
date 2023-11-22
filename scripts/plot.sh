#!/bin/bash

gnuplot << EOF
set term png
set output "$2"
set title "Rendimiento en función del número de hilos"
set xlabel "Número de hilos"
set ylabel "Tiempo de ejecución (segundos)"
set yrange[0:]
set datafile separator "\t"
plot for[col=2:$(awk -F' ' '{print NF}' $1 | head -n 1)] "$1" using 1:col title columnheader(col) with lines
EOF
