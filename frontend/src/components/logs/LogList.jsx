import React, { useState, useEffect } from 'react';
import { getLogs } from '../../services/logService';

const LogList = () => {
    const [logs, setLogs] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchLogs = async () => {
            try {
                const data = await getLogs();
                setLogs(data?.logs || []); // fallback nếu logs null
            } catch (err) {
                setError('Không thể lấy danh sách logs');
            } finally {
                setLoading(false);
            }
        };

        fetchLogs();
    }, []);

    if (loading) return <div className="loader">Đang tải log...</div>;
    if (error) return <div className="error-message">{error}</div>;

    return (
        <div className="log-list-container">
            <h1>System Logs</h1>

            {logs.length === 0 ? (
                <p>Không có log nào được tìm thấy.</p>
            ) : (
                <table className="log-table">
                    <thead>
                    <tr>
                        <th>ID</th>
                        <th>Thời gian</th>
                        <th>Loại</th>
                        <th>Hành động</th>
                        <th>IP nguồn</th>
                        <th>IP đích</th>
                    </tr>
                    </thead>
                    <tbody>
                    {logs.map(log => (
                        <tr key={log.id}>
                            <td>{log.id}</td>
                            <td>{new Date(log.time).toLocaleString()}</td>
                            <td>{log.type}</td>
                            <td>{log.action}</td>
                            <td>{log.source_ip}</td>
                            <td>{log.dest_ip || '-'}</td>
                        </tr>
                    ))}
                    </tbody>
                </table>
            )}
        </div>
    );
};

export default LogList;
