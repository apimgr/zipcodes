#!/usr/bin/env bash
# Zipcodes API Server - Installation Script
# Installs zipcodes binary and sets up systemd service

set -e

PROJECTNAME="zipcodes"
INSTALL_DIR="/usr/local/bin"
SERVICE_USER="zipcodes"
SERVICE_GROUP="zipcodes"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo "========================================="
echo "  Zipcodes API Server - Installation"
echo "========================================="
echo ""

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}Error: This script must be run as root${NC}"
    echo "Usage: sudo ./scripts/install.sh"
    exit 1
fi

# Detect platform
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

BINARY_NAME="${PROJECTNAME}-${OS}-${ARCH}"
if [ "$OS" = "darwin" ]; then
    BINARY_NAME="${PROJECTNAME}-darwin-${ARCH}"
fi

echo "Detected platform: ${OS}/${ARCH}"
echo "Binary: ${BINARY_NAME}"
echo ""

# Check if binary exists in binaries/ directory
if [ ! -f "binaries/${BINARY_NAME}" ]; then
    echo -e "${RED}Error: Binary not found: binaries/${BINARY_NAME}${NC}"
    echo "Please run 'make build' first"
    exit 1
fi

# Install binary
echo -e "${GREEN}[1/5]${NC} Installing binary..."
cp "binaries/${BINARY_NAME}" "${INSTALL_DIR}/${PROJECTNAME}"
chmod +x "${INSTALL_DIR}/${PROJECTNAME}"
echo "  ✓ Installed to ${INSTALL_DIR}/${PROJECTNAME}"

# Create service user/group (Linux only)
if [ "$OS" = "linux" ]; then
    echo -e "${GREEN}[2/5]${NC} Creating service user..."
    if ! id -u "$SERVICE_USER" >/dev/null 2>&1; then
        useradd --system --no-create-home --shell /bin/false "$SERVICE_USER"
        echo "  ✓ Created user: ${SERVICE_USER}"
    else
        echo "  ✓ User already exists: ${SERVICE_USER}"
    fi
fi

# Create directories
echo -e "${GREEN}[3/5]${NC} Creating directories..."
mkdir -p /etc/${PROJECTNAME}
mkdir -p /var/lib/${PROJECTNAME}
mkdir -p /var/log/${PROJECTNAME}

if [ "$OS" = "linux" ]; then
    chown -R ${SERVICE_USER}:${SERVICE_GROUP} /etc/${PROJECTNAME}
    chown -R ${SERVICE_USER}:${SERVICE_GROUP} /var/lib/${PROJECTNAME}
    chown -R ${SERVICE_USER}:${SERVICE_GROUP} /var/log/${PROJECTNAME}
fi

echo "  ✓ Created /etc/${PROJECTNAME}"
echo "  ✓ Created /var/lib/${PROJECTNAME}"
echo "  ✓ Created /var/log/${PROJECTNAME}"

# Create systemd service (Linux only)
if [ "$OS" = "linux" ] && [ -d "/etc/systemd/system" ]; then
    echo -e "${GREEN}[4/5]${NC} Creating systemd service..."
    cat > /etc/systemd/system/${PROJECTNAME}.service << 'EOF'
[Unit]
Description=Zipcodes API Server
After=network.target

[Service]
Type=simple
User=zipcodes
Group=zipcodes
ExecStart=/usr/local/bin/zipcodes --port 8080
Restart=always
RestartSec=10
Environment=CONFIG_DIR=/etc/zipcodes
Environment=DATA_DIR=/var/lib/zipcodes
Environment=LOGS_DIR=/var/log/zipcodes

[Install]
WantedBy=multi-user.target
EOF

    systemctl daemon-reload
    echo "  ✓ Created systemd service"
    echo ""
    echo -e "${YELLOW}To start the service:${NC}"
    echo "  sudo systemctl enable ${PROJECTNAME}"
    echo "  sudo systemctl start ${PROJECTNAME}"
    echo "  sudo systemctl status ${PROJECTNAME}"
else
    echo -e "${YELLOW}[4/5]${NC} Skipped systemd service (not Linux or systemd not available)"
fi

# Final steps
echo ""
echo -e "${GREEN}[5/5]${NC} Installation complete!"
echo ""
echo "========================================="
echo "  Next Steps"
echo "========================================="
echo ""
echo "1. Start the service:"
if [ "$OS" = "linux" ] && [ -d "/etc/systemd/system" ]; then
    echo "   sudo systemctl start ${PROJECTNAME}"
    echo ""
    echo "2. Check status:"
    echo "   sudo systemctl status ${PROJECTNAME}"
    echo ""
    echo "3. View admin credentials:"
    echo "   sudo cat /etc/${PROJECTNAME}/admin_credentials"
else
    echo "   sudo ${INSTALL_DIR}/${PROJECTNAME} --port 8080"
    echo ""
    echo "2. View admin credentials after first run:"
    echo "   Check console output or /etc/${PROJECTNAME}/admin_credentials"
fi
echo ""
echo "4. Access the server:"
echo "   http://your-server:8080"
echo ""
echo "========================================="
