import socket

HOST = '127.0.0.1'
PORT = 5001

# Create UDP socket
server_socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
server_socket.bind((HOST, PORT))

print(f"UDP server listening on {HOST}:{PORT}")

while True:
    print("Waiting to receive message...")
    data, addr = server_socket.recvfrom(1024)
    print(f"Got data from {addr}: {data.decode()}")
    print("Sending response back to client...")
    server_socket.sendto(b"Hello from UDP server", addr)
