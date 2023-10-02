#!/usr/bin/env bash
# -*- coding: utf-8 -*-

echo -e "\nCase 1: with env first, args later"
GOAT_TOKEN="GOAT-TOKEN-Latvia" security-goat -v --token gtokenrewrited --repo grepo --owner gowner --high 21 --critical 31 --low 41 --medium 100

echo -e "\nCase 1: with env first, config later"
GOAT_TOKEN="your_token_env" security-goat --config config_example.yaml

echo -e "\nCase 1: with only env"
export GOAT_GITHUB_TOKEN="your_token"
export GOAT_GITHUB_OWNER="your_username"
export GOAT_GITHUB_REPO="your_repository"
export GOAT_SEVERITY_LIMITS_CRITICAL=0
export GOAT_SEVERITY_LIMITS_HIGH=1
export GOAT_SEVERITY_LIMITS_MEDIUM=2
export GOAT_SEVERITY_LIMITS_LOW=10
security-goat
