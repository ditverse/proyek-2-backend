#!/bin/bash
# Air wrapper script - auto-detects OS and uses appropriate config

if [[ "$OSTYPE" == "linux-gnu"* ]] || [[ "$OSTYPE" == "darwin"* ]]; then
    echo "🐧 Detected Linux/macOS - using .air.linux.toml"
    air -c .air.linux.toml
elif [[ "$OSTYPE" == "msys" ]] || [[ "$OSTYPE" == "cygwin" ]] || [[ "$OSTYPE" == "win32" ]]; then
    echo "🪟 Detected Windows - using .air.windows.toml"
    air -c .air.windows.toml
else
    echo "⚠️ Unknown OS: $OSTYPE - defaulting to .air.toml"
    air
fi
