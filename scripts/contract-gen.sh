#!/usr/bin/env bash
set -e

CONTRACT_PATH='./contract'

## Generate abi and bin files from contract (force overwrite)
solc --overwrite --abi --bin $CONTRACT_PATH/killcord.sol -o $CONTRACT_PATH

## generate contract library in go for use with killcord contract
abigen \
--abi $CONTRACT_PATH/killcord.abi \
--pkg contract \
--type KillCord \
--out $CONTRACT_PATH/killcord.go \
--bin $CONTRACT_PATH/killcord.bin