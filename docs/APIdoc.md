# Asset Management Service API Documentation

- Global Error Response
    ```json
    {
        "error": "error message description"
    }
    ```

## Asset Management Service
- PORT: 8081
### Asset Operations
- `GET /api/assets`: Get list of all assets
    - Response
    ```json
    [
        {
            "ip_address": "192.168.1.100",
            "mac_address": "00:1B:44:11:3A:B7",
            "hostname": "server-01",
            "asset_type": "server",
            "location": "Data Center A",
            "operating_system": "Ubuntu 20.04",
            "first_seen": "2024-01-15T08:30:00Z",
            "last_seen": "2024-01-20T14:25:00Z",
            "is_managed": true
        }
    ]
    ```

- `GET /api/assets/{ip}`: Get specific asset by IP address
    - Path Parameters
        - `ip`: IP address of the asset
    - Response
    ```json
    {
        "ip_address": "192.168.1.100",
        "mac_address": "00:1B:44:11:3A:B7",
        "hostname": "server-01",
        "asset_type": "server",
        "location": "Data Center A",
        "operating_system": "Ubuntu 20.04",
        "first_seen": "2024-01-15T08:30:00Z",
        "last_seen": "2024-01-20T14:25:00Z",
        "is_managed": true
    }
    ```

- `POST /api/assets`: Create a new asset
    - Request
    ```json
    {
        "ip_address": "192.168.1.101",
        "mac_address": "00:1B:44:11:3A:B8",
        "hostname": "server-02",
        "asset_type": "server",
        "location": "Data Center B",
        "operating_system": "CentOS 8",
        "first_seen": "2024-01-20T09:00:00Z",
        "last_seen": "2024-01-20T09:00:00Z",
        "is_managed": false
    }
    ```
    - Response
    ```json
    {
        "ip_address": "192.168.1.101",
        "mac_address": "00:1B:44:11:3A:B8",
        "hostname": "server-02",
        "asset_type": "server",
        "location": "Data Center B",
        "operating_system": "CentOS 8",
        "first_seen": "2024-01-20T09:00:00Z",
        "last_seen": "2024-01-20T09:00:00Z",
        "is_managed": false
    }
    ```

- `PUT /api/assets/{ip}`: Update existing asset
    - Path Parameters
        - `ip`: IP address of the asset to update
    - Request
    ```json
    {
        "ip_address": "192.168.1.100",
        "mac_address": "00:1B:44:11:3A:B7",
        "hostname": "server-01-updated",
        "asset_type": "server",
        "location": "Data Center A - Rack 5",
        "operating_system": "Ubuntu 22.04",
        "first_seen": "2024-01-15T08:30:00Z",
        "last_seen": "2024-01-20T15:30:00Z",
        "is_managed": true
    }
    ```
    - Response
    ```json
    {
        "ip_address": "192.168.1.100",
        "mac_address": "00:1B:44:11:3A:B7",
        "hostname": "server-01-updated",
        "asset_type": "server",
        "location": "Data Center A - Rack 5",
        "operating_system": "Ubuntu 22.04",
        "first_seen": "2024-01-15T08:30:00Z",
        "last_seen": "2024-01-20T15:30:00Z",
        "is_managed": true
    }
    ```

- `DELETE /api/assets/{ip}`: Delete an asset
    - Path Parameters
        - `ip`: IP address of the asset to delete
    - Response
    ```json
    {
        "message": "Asset deleted successfully"
    }
    ```

- `PATCH /api/assets/{ip}/manage`: Set asset as managed
    - Path Parameters
        - `ip`: IP address of the asset to manage
    - Response
    ```json
    {
        "message": "Asset managed successfully"
    }
    ```

## Service Management
### Service Operations
- `GET /api/assets/{ip}/services`: Get all services for a specific asset
    - Path Parameters
        - `ip`: IP address of the asset
    - Response
    ```json
    [
        {
            "id": 1,
            "asset_ip": "192.168.1.100",
            "name": "HTTP Server",
            "port": 80,
            "protocol": "TCP",
            "description": "Apache HTTP Server",
            "is_managed": true
        },
        {
            "id": 2,
            "asset_ip": "192.168.1.100",
            "name": "SSH Server",
            "port": 22,
            "protocol": "TCP",
            "description": "OpenSSH Server",
            "is_managed": false
        }
    ]
    ```

- `POST /api/assets/{ip}/services`: Add a new service to an asset
    - Path Parameters
        - `ip`: IP address of the asset
    - Request
    ```json
    {
        "name": "MySQL Database",
        "port": 3306,
        "protocol": "TCP",
        "description": "MySQL Server 8.0",
        "is_managed": false
    }
    ```
    - Response
    ```json
    {
        "id": 3,
        "asset_ip": "192.168.1.100",
        "name": "MySQL Database",
        "port": 3306,
        "protocol": "TCP",
        "description": "MySQL Server 8.0",
        "is_managed": false
    }
    ```

- `PUT /api/assets/{ip}/services/{service_id}`: Update an existing service
    - Path Parameters
        - `ip`: IP address of the asset
        - `service_id`: ID of the service to update
    - Request
    ```json
    {
        "name": "MySQL Database Updated",
        "port": 3306,
        "protocol": "TCP",
        "description": "MySQL Server 8.0 - Production",
        "is_managed": true
    }
    ```
    - Response
    ```json
    {
        "id": 3,
        "asset_ip": "192.168.1.100",
        "name": "MySQL Database Updated",
        "port": 3306,
        "protocol": "TCP",
        "description": "MySQL Server 8.0 - Production",
        "is_managed": true
    }
    ```

- `DELETE /api/assets/{ip}/services/{service_id}`: Delete a service
    - Path Parameters
        - `ip`: IP address of the asset
        - `service_id`: ID of the service to delete
    - Response
    ```json
    {
        "message": "Service deleted successfully"
    }
    ```

- `PATCH /api/assets/{ip}/services/{service_id}/manage`: Set service as managed
    - Path Parameters
        - `ip`: IP address of the asset
        - `service_id`: ID of the service to manage
    - Response
    ```json
    {
        "message": "Service managed successfully"
    }
    ```
  

## Log-analysis-service
- PORT: 8080
- `GET /api/assets`: Get asset information
- Response
```json
{
    "assets": [
        {
            "asset_type": "client",
            "first_seen": "2025-06-10T00:20:49Z",
            "hostname": "",
            "id": 2,
            "ip_address": "192.168.1.10",
            "last_seen": "2025-06-10T00:20:49Z",
            "mac_address": null
        }
    ],
    "count": 1
}
```

- `GET /api/services`: Get services information
- Response
```json
{
  "count": 1,
  "services": [
    {
      "asset_ip": "192.168.1.2",
      "id": 1,
      "name": "dns",
      "port": 53,
      "protocol": ""
    }
  ]
}
```
- `GET /api/services/:ip`: Get service information in a asset
- Response
```json
{
  "count": 1,
  "services": [
    {
      "name": "dns",
      "port": 53,
      "protocol": ""
    }
  ]
}
```
- `GET /api/assets/:ip`: Get details asset information 
- Response
```json
{
  "asset": {
    "asset_type": "server",
    "first_seen": "2025-06-10T02:37:11Z",
    "hostname": "",
    "id": 7,
    "ip_address": "192.168.1.2",
    "last_seen": "2025-06-10T02:37:11Z",
    "mac_address": null,
    "metadata": "null"
  }
}
```

- `GET /api/logs`: Get parsed-log information
- Response
```json
{
  "count": 2,
  "logs": [
    {
      "action": "mapping",
      "dest_ip": null,
      "id": 33,
      "source_ip": "192.168.1.1",
      "time": "2025-06-10T02:37:11Z",
      "type": "mac_ip_mapping"
    },
    {
      "action": "mapping",
      "dest_ip": null,
      "id": 32,
      "source_ip": "192.168.1.12",
      "time": "2025-06-10T02:37:11Z",
      "type": "mac_ip_mapping"
    }
  ]
}
```


````
                  ┌─────────────────┐
                  │  Client         │
                  │ Application     │
                  └─────────┬───────┘
                            │
                            │
                  ┌─────────▼───────┐
                  │ API Gateway     │
                  │ (REST API)      │
                  └────────┬────────┘
                           │
          ┐────────────────├─────────────────┐
          │                                  │
┌─────────▼───────┐                ┌─────────▼────────┐
│ Log Analysis    │                │ Asset Management │
│ Service         │                │ Service          │
└─────────┬───────┘                └─────────┬────────┘
          │                                  │
          │                                  │
┌─────────▼───────┐                ┌─────────▼────────┐
│ PostgreSQL      │                │ PostgreSQL       │
│ (Logs DB)       │                │ (Assets DB)      │
└─────────────────┘                └──────────────────┘
````
````
[1] Hệ thống gửi log  ───────────────▶  Log Analysis Service

[2] Log Analysis Service:
├─ Phân tích log
└─ Phát hiện tài sản (Asset) và dịch vụ (Service)

[3] Gửi dữ liệu phân tích ─────────▶  Asset Management Service

[4] Asset Management Service:
├─ Lưu thông tin tài sản & dịch vụ vào PostgreSQL
└─ Đánh dấu trạng thái mặc định: "unmanaged"

[5] Người dùng (Admin/User) gửi yêu cầu ──▶  API Gateway
└─ API Gateway định tuyến đến Asset Management Service
(truy vấn, cập nhật, đánh dấu asset)

[6] Khi tất cả dịch vụ liên quan đã "managed":
└─ Asset được cập nhật trạng thái: "managed"

[7] Admin có thể:
├─ Theo dõi tổng quan tài sản
├─ Lọc theo trạng thái (managed/unmanaged)
└─ Thống kê theo IP, dịch vụ, mốc thời gian...

````