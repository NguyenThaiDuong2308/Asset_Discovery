import React, { useState, useEffect } from 'react';
import { Link } from 'react-router-dom';
import { getAssets } from '../../services/assetService';
import AssetCard from './AssetCard';

const AssetList = () => {
    const [assets, setAssets] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchAssets = async () => {
            try {
                const data = await getAssets();

                // Đảm bảo data là mảng, nếu không thì fallback []
                if (Array.isArray(data)) {
                    setAssets(data);
                } else {
                    console.warn("getAssets() response is not an array:", data);
                    setAssets([]);
                }

                setLoading(false);
            } catch (err) {
                console.error("Failed to fetch assets:", err);
                setError('Failed to fetch assets');
                setAssets([]); // fallback để tránh null
                setLoading(false);
            }
        };

        fetchAssets();
    }, []);

    if (loading) return <div className="loader">Loading...</div>;
    if (error) return <div className="error-message">{error}</div>;

    return (
        <div className="asset-list-container">
            <div className="list-header">
                <h1>Assets</h1>
                <Link to="/assets/new" className="btn btn-primary">Add New Asset</Link>
            </div>

            {!assets || assets.length === 0 ? (
                <p>No assets found. Add one to get started.</p>
            ) : (
                <div className="asset-grid">
                    {assets.map(asset => (
                        <AssetCard key={asset.ip_address} asset={asset} />
                    ))}
                </div>
            )}
        </div>
    );
};

export default AssetList;
