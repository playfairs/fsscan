#!/bin/bash

echo "Installing fs command..."
go install .

if [ -d "$HOME/.nix-profile/bin" ]; then
    echo "Nix environment detected"
    echo "Copying to ~/.nix-profile/bin/fs..."
    if [ -w "$HOME/.nix-profile/bin" ]; then
        cp "$BINARY_PATH" "$HOME/.nix-profile/bin/fs"
        chmod +x "$HOME/.nix-profile/bin/fs"
        echo "fs command installed to Nix profile!"
        echo "You can now run 'fs' from anywhere"
    else
        echo "Need sudo to copy to Nix profile"
        sudo cp "$BINARY_PATH" "$HOME/.nix-profile/bin/fs"
        sudo chmod +x "$HOME/.nix-profile/bin/fs"
        echo "fs command installed to Nix profile (with sudo)!"
        echo "You can now run 'fs' from anywhere"
    fi
else
    echo "Traditional environment detected"
    
    if [[ ":$PATH:" != *":$HOME/go/bin:"* ]]; then
        echo "GOPATH/bin is not in your PATH"
        echo "Add one of these to your shell profile (~/.zshrc, ~/.bashrc, etc.):"
        echo ""
        echo "export PATH=\"\$HOME/go/bin:\$PATH\""
        echo "# OR"
        echo "export PATH=\"\$HOME/go/bin:\$PATH\""
        echo ""
        echo "Then restart your shell or run: source ~/.zshrc"
        echo ""
        
        if [ -w "/usr/local/bin" ]; then
            ln -sf "$BINARY_PATH" "/usr/local/bin/fs"
            echo "Created symlink: /usr/local/bin/fs"
        else
            echo "You can also manually create a symlink:"
            echo "sudo ln -sf \"$BINARY_PATH\" /usr/local/bin/fs"
        fi
    else
        echo "fs command installed successfully!"
        echo "You can now run 'fs' from anywhere"
        echo ""
    fi
fi

echo ""
echo "Testing installation..."
if command -v fs &> /dev/null; then
    echo "fs command is available!"
    fs --help | head -n 3
else
    echo "fs command not found in PATH"
    echo "You may need to restart your shell or update your PATH"
fi
