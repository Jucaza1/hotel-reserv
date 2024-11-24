#!/bin/sh

# Execute seed binary
echo "Running first binary..."
./bin/seed
if [ $? -ne 0 ]; then
    echo "Error: First binary failed."
    exit 1
fi

# Execute api binary
echo "Running second binary..."
./bin/api
if [ $? -ne 0 ]; then
    echo "Error: Second binary failed."
    exit 1
fi

echo "Both binaries executed successfully."
