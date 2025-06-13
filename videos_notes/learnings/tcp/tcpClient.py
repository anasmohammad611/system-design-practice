import socket
import time

HOST = '127.0.0.1'
PORT = 5500

print("Creating socket...")
client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)

time.sleep(2)
print("Connecting to server (handshake starts now)...")
client_socket.connect((HOST, PORT))
print("[✓] Connected — handshake done.")

time.sleep(1)
print("Sending message...")
client_socket.sendall(b"Hello from client")

# client_socket.sendall(b"Hello dsaadsasd dsasd")

print("Waiting for response...")
data = client_socket.recv(1024)
print("Received:", data.decode())

print("Closing client socket...")
client_socket.close()
