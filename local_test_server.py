#!/usr/bin/env python3
"""
Local test server for Library Management System frontend
This script serves the GitHub Pages files locally for testing
"""

import http.server
import socketserver
import os
import sys
from pathlib import Path

def main():
    # Change to docs directory
    docs_path = Path(__file__).parent / "docs"
    
    if not docs_path.exists():
        print(f"Error: {docs_path} does not exist!")
        sys.exit(1)
    
    os.chdir(docs_path)
    port = 8000
    
    # Parse command line argument for port
    if len(sys.argv) > 1:
        try:
            port = int(sys.argv[1])
        except ValueError:
            print("Usage: python local_test_server.py [port]")
            sys.exit(1)
    
    handler = http.server.SimpleHTTPRequestHandler
    httpd = socketserver.TCPServer(("", port), handler)
    
    print(f"Local test server started at http://localhost:{port}")
    print(f"Serving files from: {docs_path}")
    print("Press Ctrl+C to stop the server")
    
    try:
        httpd.serve_forever()
    except KeyboardInterrupt:
        print("\nServer stopped.")
        httpd.server_close()

if __name__ == "__main__":
    main()