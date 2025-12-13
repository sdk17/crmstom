#!/bin/bash
# Server setup script for CRM Stomatology
# Run this on the remote server to set up the environment

set -e

APP_DIR="/opt/crmstom"
REPO_URL="git@github.com:sdk17/crmstom.git"

echo "=== CRM Stomatology Server Setup ==="

# Install Docker if not present
if ! command -v docker &> /dev/null; then
    echo "Installing Docker..."
    curl -fsSL https://get.docker.com | sh
    systemctl enable docker
    systemctl start docker
fi

# Install Docker Compose plugin if not present
if ! docker compose version &> /dev/null; then
    echo "Installing Docker Compose plugin..."
    apt-get update
    apt-get install -y docker-compose-plugin
fi

# Create application directory
echo "Creating application directory..."
mkdir -p $APP_DIR
cd $APP_DIR

# Clone repository if not exists
if [ ! -d ".git" ]; then
    echo "Cloning repository..."
    git clone $REPO_URL .
else
    echo "Repository already exists, pulling latest..."
    git pull origin main
fi

# Create .env file if not exists
if [ ! -f ".env" ]; then
    echo "Creating .env file..."
    cat > .env << 'EOF'
# Database Configuration
DB_NAME=crmstom
DB_USER=crmstom_user
DB_PASSWORD=CHANGE_THIS_TO_SECURE_PASSWORD

# Application Configuration
APP_PORT=8080
EOF
    echo ""
    echo "IMPORTANT: Edit /opt/crmstom/.env and set a secure DB_PASSWORD!"
    echo ""
fi

echo "=== Setup Complete ==="
echo ""
echo "Next steps:"
echo "1. Edit /opt/crmstom/.env and set a secure DB_PASSWORD"
echo "2. Run: cd /opt/crmstom && docker compose -f docker-compose.prod.yml up -d"
echo ""
