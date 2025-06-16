
````
                  ┌─────────────────┐
                  │  Client         │
                  │ Application     │
                  └────────┬────────┘
                           │
                           │
                  ┌────────▼────────┐
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

[5] Người dùng gửi yêu cầu ──▶  API Gateway
└─ API Gateway định tuyến đến Asset Management Service
(truy vấn, cập nhật, đánh dấu asset)

[6] Người dùng có thể:
├─ Theo dõi tổng quan tài sản
├─ Lọc theo trạng thái (managed/unmanaged)
└─ Thống kê theo IP, dịch vụ, mốc thời gian...

````