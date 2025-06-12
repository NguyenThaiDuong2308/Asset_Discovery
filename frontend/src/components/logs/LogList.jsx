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
                setLogs(data.logs);
                setLoading(false);
            } catch (err) {
                setError('Failed to fetch logs');
                setLoading(false);
            }
        };

        fetchLogs();
    }, []);

    if (loading) return <div className="loader">Loading...</div>;
    if (error) return <div className="error-message">{error}</div>;

    return (
        <div className="log-list-container">
            <h1>System Logs</h1>

            {logs.length === 0 ? (
                <p>No logs found.</p>
            ) : (
                <table className="log-table">
                    <thead>
                    <tr>
                        <th>ID</th>
                        <th>Time</th>
                        <th>Type</th>
                        <th>Action</th>
                        <th>Source IP</th>
                        <th>Destination IP</th>
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