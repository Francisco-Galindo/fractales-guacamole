#!/bin/bash

frames=128
size=128
iterations=5

while getopts "n:s:i" opt; do
	case $opt in
		n) frames=$OPTARG   ;;
		i) iterations=$OPTARG   ;;
		s) size=$OPTARG   ;;
		*) echo 'Script: error al leer opciones' >&2
			exit 1
	esac
done


printf "No. de hilos\tTiempo promedio\n"
i=1
while [ $i -le 32 ]
do
	total=0
	for j in $(seq $iterations)
	do
		./main -t $i -n $frames -s $size >/dev/null 2>timefile
		elapsed=$(cat timefile)
		total=$(echo "$total+$elapsed" | bc -l)
	done

	rm timefile

	avg=$(echo "$total / $iterations" | bc -l)
	printf "$i\t$avg\n"

	((i=i+1))
done
