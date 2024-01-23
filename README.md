# Rubrik

Rubrik is a simple key-value document store implemented in Go. It uses a 3D cube structure to organize documents with support for set, get, delete operations. The database stores documents on disk, avoiding the need to keep the entire cube in memory.

## Features

- **3D Cube Structure**: The database uses a 3D cube structure to organize documents based on their coordinates (x, y, z).

- **Disk Storage**: Documents are stored on disk, allowing for scalability and efficient use of resources.

- **Set, Get, Delete Operations**: Supports basic set, get, and delete operations to manipulate documents in the cube.
