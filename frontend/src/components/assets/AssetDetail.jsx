import React, { useState, useEffect } from 'react';
import { useParams, Link, useNavigate } from 'react-router-dom';
import { getAssetByIp, deleteAsset, setAssetManaged } from '../../services/assetService';

const AssetDetail = () => {
    const { ip } = useParams();
    const navigate = useNavigate();
    const [asset, setAsset] = useState(null);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchAsset = async () => {
            try {
                const data = await getAssetByIp(ip);
                setAsset(data);
                setLoading(false);
            } catch (err) {
                setError('Failed to fetch asset details');
                setLoading(false);
            }
        };

        fetchAsset();
    }, [ip]);

    const handleDelete = async () => {
        if (window.confirm('Are you sure you want to delete this asset?')) {
            try {
                await deleteAsset(ip);
                navigate('/assets');
            } catch (err) {
                setError('Failed to delete asset');
            }
        }
    };

    const handleManage = async () => {
        try {
            await setAssetManaged(ip);
            setAsset({ ...asset, is_managed: true });
        } catch (err) {
            setError('Failed to set asset as managed');
        }
    };

    if (loading) return <div className="loader">Loading...</div>;
    if (error) return <div className="error-message">{error}</div>;
    if (!asset) return <div className="not-found">Asset not found</div>;

    return (
        <div className="asset-detail-container">
            <div className="detail-header">
                <h1>{asset.hostname || 'Unnamed Asset'}</h1>
                <div className="action-buttons">
                    <Link to={`/assets/${ip}/edit`} className="btn btn-primary">Edit</Link>
                    <Link to={`/assets/${ip}/services`} className="btn btn-info">Services</Link>
                    {!asset.is_managed && (
                        <button onClick={handleManage} className="btn btn-success">Set as Managed</button>
                    )}
                    <button onClick={handleDelete} className="btn btn-danger">Delete</button>
                </div>
            </div>

            <div className="detail-content">
                <div className="detail-row">
                    <div className="detail-label">IP Address</div>
                    <div className="detail-value">{asset.ip_address}</div>
                </div>
                <div className="detail-row">
                    <div className="detail-label">MAC Address</div>
                    <div className="detail-value">{asset.mac_address}</div>
                </div>
                <div className="detail-row">
                    <div className="detail-label">Asset Type</div>
                    <div className="detail-value">{asset.asset_type}</div>
                </div>
                <div className="detail-row">
                    <div className="detail-label">Location</div>
                    <div className="detail-value">{asset.location}</div>
                </div>
                <div className="detail-row">
                    <div className="detail-label">Operating System</div>
                    <div className="detail-value">{asset.operating_system}</div>
                </div>
                <div className="detail-row">
                    <div className="detail-label">First Seen</div>
                    <div className="detail-value">{new Date(asset.first_seen).toLocaleString()}</div>
                </div>
                <div className="detail-row">
                    <div className="detail-label">Last Seen</div>
                    <div className="detail-value">{new Date(asset.last_seen).toLocaleString()}</div>
                </div>
                <div className="detail-row">
                    <div className="detail-label">Status</div>
                    <div className="detail-value">
            <span className={`badge ${asset.is_managed ? 'badge-success' : 'badge-warning'}`}>
              {asset.is_managed ? 'Managed' : 'Unmanaged'}
            </span>
                    </div>
                </div>
            </div>
        </div>
    );
};

export default AssetDetail;