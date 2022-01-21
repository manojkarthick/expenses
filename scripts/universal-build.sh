#!/usr/bin/env bash

ASSET_VERSION="${LATEST_VERSION:1}"
echo "Asset version is $ASSET_VERSION"

echo "Downloading tarballs..."
wget "https://github.com/manojkarthick/expenses/releases/download/${LATEST_VERSION}/expenses_${ASSET_VERSION}_darwin_amd64.tar.gz" -O expenses-amd64.tar.gz
wget "https://github.com/manojkarthick/expenses/releases/download/${LATEST_VERSION}/expenses_${ASSET_VERSION}_darwin_arm64.tar.gz" -O expenses-arm64.tar.gz
mkdir -pv amd64 arm64 universal
tar -xzvf expenses-amd64.tar.gz -C amd64
tar -xzvf expenses-arm64.tar.gz -C arm64

echo "Building universal binary..."
lipo -create -output expenses amd64/expenses arm64/expenses
tar -czvf expenses_darwin_universal.tar.gz ./expenses
