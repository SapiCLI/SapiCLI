#!/bin/bash

# Check if running as root
if [[ $EUID -ne 0 ]]; then
   echo "Error: This script must be run as root."
   exit 1
fi

# Function to check if a command is available
command_exists() {
    command -v $1 >/dev/null 2>&1
}

# Function to install packages using apt package manager
install_packages_apt() {
    echo "Installing missing dependencies using apt: ${missing_dependencies[*]}"
    apt-get update
    apt-get install -y "${missing_dependencies[@]}"
}

# Function to install packages using yum package manager
install_packages_yum() {
    echo "Installing missing dependencies using yum: ${missing_dependencies[*]}"
    yum install -y "${missing_dependencies[@]}"
}

# Check package manager and set installation function
if command_exists apt-get; then
    package_manager="apt"
    install_packages_function="install_packages_apt"
elif command_exists yum; then
    package_manager="yum"
    install_packages_function="install_packages_yum"
else
    echo "Error: Neither apt nor yum package manager found."
    exit 1
fi

# Check dependencies
echo "Checking dependencies..."
dependencies=("curl" "wget" "screen")
missing_dependencies=()

for dependency in "${dependencies[@]}"; do
    if ! command_exists $dependency; then
        missing_dependencies+=($dependency)
    fi
done

if [[ ${#missing_dependencies[@]} -ne 0 ]]; then
    $install_packages_function
fi

# Download and install scli
echo "Downloading scli..."
wget -O scli "https://github.com/forkyyy/SapiCLI/raw/main/scli"
chmod +x scli
mv scli /usr/bin/
mkdir -p /home/root/

# Display a note
echo "--------------------------------------------------------------"
echo "Note:"
echo "With this script, you can use  screen command unlimitedly  "
echo "(if you add a way to send X concurrents to X targets using "
echo "X method etc etc every X seconds).You can also use it like "
echo "a normal client. For detailed explanation: "
echo "https://github.com/forkyyy/SapiCLI."
echo "sript run: scli"
echo "Coded By The old ones know it as T13R :)"
echo "--------------------------------------------------------------"

scli
