import socket

HOST = '127.0.0.1'
PORT = 5001

client_socket = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)

print("Sending message to server...")
client_socket.sendto(b"Hello from UDP client", (HOST, PORT))
print("Message sent. Waiting for response...")

data, addr = client_socket.recvfrom(1024)

print("Received from server:", data.decode())

client_socket.close()
