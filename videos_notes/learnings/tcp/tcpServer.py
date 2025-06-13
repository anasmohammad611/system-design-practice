import socket
import time

HOST = '127.0.0.1'
PORT = 5500

print("Starting TCP server (before bind)")
server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

print("Binding to address...")
server_socket.bind((HOST, PORT))

print("Listening for connections...")
server_socket.listen()

print("Waiting for client to connect (handshake will happen now)...")
conn, addr = server_socket.accept()
print(f"[✓] Handshake complete. Client connected from {addr}")

while True:
    data = conn.recv(1024)
    if not data:
        break
    print("Received:", data.decode())
    conn.sendall(b"Hello from server")

print("Closing connection...")
conn.close()
