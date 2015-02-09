#!/bin/bash

for f in results/*.py; do
    python $f
done

for f in data/*.py; do
    python $f
done
