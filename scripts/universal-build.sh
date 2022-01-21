#!/usr/bin/env bash

LATEST_VERSION=$(git tag | grep ^v | sort -V | tail -n 1)
ASSET_VERSION="${LATEST_VERSION:1}"

echo "Latest version is $LATEST_VERSION"
echo "Asset version is $ASSET_VERSION"

echo "Downloading tarballs..."
wget "https://github.com/manojkarthick/expenses/releases/download/${LATEST_VERSION}/expenses_${ASSET_VERSION}_darwin_amd64.tar.gz" -O expenses-amd64.tar.gz
wget "https://github.com/manojkarthick/expenses/releases/download/${LATEST_VERSION}/expenses_${ASSET_VERSION}_darwin_arm64.tar.gz" -O expenses-arm64.tar.gz

echo "Extracting..."
mkdir -pv amd64 arm64 universal
tar -xzvf expenses-amd64.tar.gz -C amd64
tar -xzvf expenses-arm64.tar.gz -C arm64

echo "Building universal binary..."
lipo -create -output expenses amd64/expenses arm64/expenses
mv expenses universal/
tar -czvf "universal/expenses_${ASSET_VERSION}_darwin_universal.tar.gz" ./universal/expenses

echo "Cleaning up"
rm -rf amd64 arm64
