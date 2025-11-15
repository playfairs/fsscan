#!/bin/bash

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m'

BINARY_NAME="fs"
DEMO_BINARY="fs-demo"

echo -e "${BLUE}=== fsscan Build Script ===${NC}"

if ! command -v go &> /dev/null; then
    echo -e "${RED}Error: Go is not installed or not in PATH${NC}"
    exit 1
fi
echo -e "${BLUE}Go version:${NC} $(go version)"

echo -e "${YELLOW}Cleaning previous builds...${NC}"
rm -f "$BINARY_NAME" "$DEMO_BINARY"

echo -e "${YELLOW}Tidying Go modules...${NC}"
go mod tidy

echo -e "${YELLOW}Checking code with go vet...${NC}"
go vet ./...

echo -e "${YELLOW}Formatting code...${NC}"
go fmt ./...

echo -e "${YELLOW}Building main application...${NC}"
go build -ldflags="-s -w" -o "$BINARY_NAME" .

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Main application built successfully: $BINARY_NAME${NC}"
else
    echo -e "${RED}✗ Failed to build main application${NC}"
    exit 1
fi

echo -e "${YELLOW}Building demo application...${NC}"
go build -ldflags="-s -w" -o "$DEMO_BINARY" ./cmd/demo

if [ $? -eq 0 ]; then
    echo -e "${GREEN}✓ Demo application built successfully: $DEMO_BINARY${NC}"
else
    echo -e "${RED}✗ Failed to build demo application${NC}"
    exit 1
fi

if ls *_test.go 1> /dev/null 2>&1; then
    echo -e "${YELLOW}Running tests...${NC}"
    go test -v ./...
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓ All tests passed${NC}"
    else
        echo -e "${RED}✗ Some tests failed${NC}"
        exit 1
    fi
fi

echo -e "\n${BLUE}=== Build Results ===${NC}"

if [ -f "$BINARY_NAME" ]; then
    SIZE=$(ls -lh "$BINARY_NAME" | awk '{print $5}')
    echo -e "$BINARY_NAME (${SIZE})"
fi

if [ -f "$DEMO_BINARY" ]; then
    SIZE=$(ls -lh "$DEMO_BINARY" | awk '{print $5}')
    echo -e "$DEMO_BINARY (${SIZE})"
fi

echo -e "\n${BLUE}Usage:${NC}"
echo -e "  ${GREEN}Full system scan:${NC}     ./$BINARY_NAME"
echo -e "  ${GREEN}With sudo:${NC}           sudo ./$BINARY_NAME"
echo -e "  ${GREEN}Demo (current dir):${NC}  ./$DEMO_BINARY"
echo -e "  ${GREEN}Demo (custom path):${NC}  ./$DEMO_BINARY /path/to/scan"

echo -e "\n${GREEN}✓ Build completed successfully!${NC}"

chmod +x "$BINARY_NAME" "$DEMO_BINARY"

echo -e "\n${YELLOW}Note:${NC} For full system scan, you may need to run with sudo privileges"
echo -e "${YELLOW}Warning:${NC} Full system scan can take hours and use significant I/O resources"
