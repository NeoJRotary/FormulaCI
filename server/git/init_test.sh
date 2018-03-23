#!/bin/bash
mkdir -p test
cd test
git init
git remote add rA https://github.com/NeoJRotary/FormulaCI
git remote add rB https://github.com/NeoJRotary/FormulaCI
git remote add rC https://github.com/NeoJRotary/FormulaCI
git checkout -b bA
git checkout -b bB
git checkout -b bC
