#!/bin/bash
set -e # make script fail on first error
SCRIPT_PATH="$(dirname "$0")"
source "$SCRIPT_PATH/../script.inc"
pushd "$SCRIPT_PATH"

DOTNET_CLI_TELEMETRY_OPTOUT=1

if [[ "$1" = "style" ]]
then
  INFO Check style of C# Benchmarks
  dotnet tool restore

  dotnet format whitespace Benchmarks.slnx --verify-no-changes
  dotnet format style Benchmarks.slnx --verify-no-changes
  dotnet format analyzers Benchmarks.slnx --verify-no-changes

  dotnet jb inspectcode Benchmarks.slnx --format=Text --output=-
else
  INFO Build C# Benchmarks
  dotnet build Benchmarks.slnx -c Release -p:DebugType=none
fi
