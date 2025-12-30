import http.server
import socketserver

class SecurityHeadersHandler(http.server.SimpleHTTPRequestHandler):
    def end_headers(self):
        self.send_header("Cross-Origin-Opener-Policy", "same-origin")
        self.send_header("Cross-Origin-Embedder-Policy", "require-corp")
    
        self.send_header("Cache-Control", "no-cache, no-store, must-revalidate")
        self.send_header("Pragma", "no-cache")
        self.send_header("Expires", "0")
        
        super().end_headers()

PORT = 8000

with socketserver.TCPServer(("", PORT), SecurityHeadersHandler) as httpd:
    print(f"Serving at http://localhost:{PORT}")
    
    httpd.serve_forever()