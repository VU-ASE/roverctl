# Use Ubuntu 22.04 as the base image
FROM ubuntu:22.04

# Set environment variables to prevent interactive prompts during package installation
ARG DEBIAN_FRONTEND=noninteractive

# Update package list and install necessary tools and dependencies
RUN apt-get update && \
    apt-get install -y \
    git \
    curl \
    gnupg \
    make \
    build-essential \
    pkg-config \
    xz-utils \
    && rm -rf /var/lib/apt/lists/*

# Install Go 1.21 based on the system architecture
RUN ARCH=$(dpkg --print-architecture) && \
    if [ "$ARCH" = "amd64" ]; then \
    GO_ARCH="go1.21.0.linux-amd64.tar.gz"; \
    elif [ "$ARCH" = "arm64" ]; then \
    GO_ARCH="go1.21.0.linux-arm64.tar.gz"; \
    else \
    echo "Unsupported architecture"; exit 1; \
    fi && \
    curl -fsSL "https://go.dev/dl/$GO_ARCH" -o "$GO_ARCH" && \
    tar -C /usr/local -xzf "$GO_ARCH" && \
    rm "$GO_ARCH" && \
    # Install Go language server
    /usr/local/go/bin/go install golang.org/x/tools/gopls@latest

# Install golangci-lint
RUN curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b /usr/local/bin v1.54.2

# Set up Go environment variables
ENV PATH="/usr/local/go/bin:/usr/local/bin:${PATH}"

# Ensure CGO is enabled for the zmq4 package
ENV CGO_ENABLED=1

# Install Node.js and NPM based on the system architecture
RUN ARCH=$(dpkg --print-architecture) && \
    if [ "$ARCH" = "amd64" ]; then \
    NODE_ARCH="node-v20.0.0-linux-x64.tar.xz"; \
    elif [ "$ARCH" = "arm64" ]; then \
    NODE_ARCH="node-v20.0.0-linux-arm64.tar.xz"; \
    else \
    echo "Unsupported architecture"; exit 1; \
    fi && \
    curl -fsSL "https://nodejs.org/dist/v20.0.0/$NODE_ARCH" -o "$NODE_ARCH" && \
    tar -C /usr/local --strip-components=1 -xf "$NODE_ARCH" && \
    rm "$NODE_ARCH"

# Ensure Node.js and NPM are in the PATH
ENV PATH="/usr/local/bin:${PATH}"

# Install Java based on architecture
RUN ARCH=$(dpkg --print-architecture) && \
    if [ "$ARCH" = "amd64" ]; then \
    JAVA_ARCH="openjdk-17-jdk"; \
    elif [ "$ARCH" = "arm64" ]; then \
    JAVA_ARCH="openjdk-17-jdk"; \
    else \
    echo "Unsupported architecture"; exit 1; \
    fi && \
    apt-get update && \
    apt-get install -y $JAVA_ARCH && \
    rm -rf /var/lib/apt/lists/*

# Verify Java installation
RUN java -version

# Verify Node.js and NPM installation
RUN node -v && npm -v

# Install openapi-generator-cli globally
RUN npm install -g @openapitools/openapi-generator-cli

# Verify openapi-generator-cli installation
RUN openapi-generator-cli version
