import React, { useState, useEffect } from 'react';
import {
    getAssets,
    setAssetManaged
} from '../../services/assetService';
import { useNavigate } from 'react-router-dom';
import './dashboard.css';

const Dashboard = () => {
    const navigate = useNavigate();
    const [activeTab, setActiveTab] = useState('all');
    const [assets, setAssets] = useState([]);
    const [stats, setStats] = useState({
        total: 0,
        managed: 0,
        unmanaged: 0
    });
    const [loading, setLoading] = useState(true);
    const [searchTerm, setSearchTerm] = useState('');
    const [searchType, setSearchType] = useState('name'); // 'name' or 'ip'
    const [assetTypeFilter, setAssetTypeFilter] = useState('all');
    const [selectedAsset, setSelectedAsset] = useState(null);

    useEffect(() => {
        fetchData();
    }, []);

    const fetchData = async () => {
        setLoading(true);
        try {
            // Get assets from both services
            let assetManagementData = await getAssets();
            if (!assetManagementData) assetManagementData =[]

            // Merge and deduplicate assets based on IP address
            const allAssets = assetManagementData;

            setAssets(allAssets);
            setStats({
                total: allAssets.length,
                managed: allAssets.filter(asset => asset.is_managed).length,
                unmanaged: allAssets.filter(asset => !asset.is_managed).length
            });
        } catch (error) {
            console.error('Error fetching data:', error);
        } finally {
            setLoading(false);
        }
    };

    // Merge assets from both services, prioritizing asset management service data
    const mergeAssets = (assetMgmtData, logAnalysisAssets) => {
        const assetMap = new Map();

        // Add asset management service data
        assetMgmtData.forEach(asset => {
            assetMap.set(asset.ip_address, {
                ...asset,
                source: 'asset_management'
            });
        });

        // Add log analysis assets that don't exist in asset management
        logAnalysisAssets.forEach(asset => {
            if (!assetMap.has(asset.ip_address)) {
                assetMap.set(asset.ip_address, {
                    ip_address: asset.ip_address,
                    mac_address: asset.mac_address || 'N/A',
                    hostname: asset.hostname || 'Unnamed Asset',
                    asset_type: asset.asset_type || 'unknown',
                    location: 'N/A',
                    operating_system: 'N/A',
                    first_seen: asset.first_seen,
                    last_seen: asset.last_seen,
                    is_managed: false,
                    source: 'log_analysis'
                });
            }
        });

        return Array.from(assetMap.values());
    };

    const handleSearch = (e) => {
        e.preventDefault();
        // Filter is handled in the filtered assets below
    };

    const handleManageAsset = async (ip, e) => {
        e.stopPropagation(); // Prevent opening the asset details
        try {
            await setAssetManaged(ip);
            // Update local state
            setAssets(assets.map(asset =>
                asset.ip_address === ip ? { ...asset, is_managed: true } : asset
            ));
            // Update stats
            setStats({
                ...stats,
                managed: stats.managed + 1,
                unmanaged: stats.unmanaged - 1
            });
        } catch (error) {
            console.error('Error managing asset:', error);
        }
    };

    const handleAssetClick = (asset) => {
        navigate(`/assets/${asset.ip_address}`);
    };

    const handleCloseModal = () => {
        setSelectedAsset(null);
    };

    // Get unique asset types for the filter dropdown
    const assetTypes = ['all', ...new Set(assets.map(asset => asset.asset_type).filter(Boolean))];

    // Filter assets based on search term, active tab, and asset type
    const filteredAssets = assets.filter(asset => {
        // Filter by tab
        if (activeTab === 'managed' && !asset.is_managed) return false;
        if (activeTab === 'unmanaged' && asset.is_managed) return false;

        // Filter by asset type (only in 'all' tab)
        if (activeTab === 'all' && assetTypeFilter !== 'all' && asset.asset_type !== assetTypeFilter) {
            return false;
        }

        // Filter by search term
        if (searchTerm) {
            if (searchType === 'name') {
                return asset.hostname?.toLowerCase().includes(searchTerm.toLowerCase());
            } else {
                return asset.ip_address?.includes(searchTerm);
            }
        }

        return true;
    });

    const formatDate = (dateString) => {
        if (!dateString) return 'N/A';
        const date = new Date(dateString);
        return `${date.getHours()}:${date.getMinutes()}:${date.getSeconds()} ${date.getDate()}/${date.getMonth()+1}/${date.getFullYear()}`;
    };

    return (
        <div className="dashboard-container">
            <div className="search-bar">
                <div className="search-input-container">
                    <i className="search-icon">🔍</i>
                    <input
                        type="text"
                        placeholder={`Tìm kiếm tài sản theo ${searchType === 'name' ? 'tên' : 'IP'}...`}
                        className="search-input"
                        value={searchTerm}
                        onChange={(e) => setSearchTerm(e.target.value)}
                    />
                </div>
                <div className="filter-container">
                    <select
                        className="filter-dropdown"
                        onChange={(e) => setSearchType(e.target.value)}
                        value={searchType}
                    >
                        <option value="name">Tìm theo tên</option>
                        <option value="ip">Tìm theo IP</option>
                    </select>
                    {activeTab === 'all' && (
                        <select
                            className="filter-dropdown"
                            onChange={(e) => setAssetTypeFilter(e.target.value)}
                            value={assetTypeFilter}
                        >
                            <option value="all">Tất cả loại</option>
                            {assetTypes.filter(type => type !== 'all').map(type => (
                                <option key={type} value={type}>{type}</option>
                            ))}
                        </select>
                    )}
                </div>
            </div>

            <div className="stats-container">
                <div className="stat-card">
                    <div className="stat-header">
                        <h3>Tổng tài sản</h3>
                        <button className="refresh-btn" onClick={fetchData}>⟳</button>
                    </div>
                    <div className="stat-content">
                        <div className="stat-count blue">{stats.total}</div>
                        <div className="stat-description">Số tài sản thuộc</div>
                    </div>
                </div>

                <div className="stat-card">
                    <div className="stat-header">
                        <h3>Managed Assets</h3>
                        <button className="refresh-btn" onClick={fetchData}>⟳</button>
                    </div>
                    <div className="stat-content">
                        <div className="stat-count green">{stats.managed}</div>
                        <div className="stat-description">Được quản lý</div>
                    </div>
                </div>

                <div className="stat-card">
                    <div className="stat-header">
                        <h3>Unmanaged Assets</h3>
                        <button className="refresh-btn" onClick={fetchData}>⟳</button>
                    </div>
                    <div className="stat-content">
                        <div className="stat-count red">{stats.unmanaged}</div>
                        <div className="stat-description">Chưa được quản lý</div>
                    </div>
                </div>
            </div>

            <div className="assets-container">
                <div className="tabs">
                    <button
                        className={`tab-btn ${activeTab === 'all' ? 'active' : ''}`}
                        onClick={() => setActiveTab('all')}
                    >
                        Tất cả ({stats.total})
                    </button>
                    <button
                        className={`tab-btn ${activeTab === 'managed' ? 'active' : ''}`}
                        onClick={() => setActiveTab('managed')}
                    >
                        Managed ({stats.managed})
                    </button>
                    <button
                        className={`tab-btn ${activeTab === 'unmanaged' ? 'active' : ''}`}
                        onClick={() => setActiveTab('unmanaged')}
                    >
                        Unmanaged ({stats.unmanaged})
                    </button>
                </div>

                {loading ? (
                    <div className="loading-container">
                        <p>Loading assets...</p>
                    </div>
                ) : (
                    <div className="assets-list">
                        {filteredAssets.length > 0 ? (
                            filteredAssets.map((asset, index) => (
                                <div
                                    key={`${asset.ip_address}-${index}`}
                                    className="asset-card"
                                    onClick={() => handleAssetClick(asset)}
                                >
                                    <div className="asset-header">
                                        <div className="asset-icon-container">
                                            {asset.asset_type === 'server' ? (
                                                <span className="asset-icon">💻</span>
                                            ) : asset.asset_type === 'network_device' ? (
                                                <span className="asset-icon">🔌</span>
                                            ) : (
                                                <span className="asset-icon">📱</span>
                                            )}
                                        </div>
                                        <h3>{asset.hostname || 'Unnamed Asset'}</h3>
                                        <div className={`asset-status ${asset.is_managed ? 'managed' : 'unmanaged'}`}>
                                            {asset.is_managed ? 'Managed' : 'Unmanaged'}
                                        </div>
                                    </div>

                                    <div className="asset-details">
                                        <div className="detail-row">
                                            <span className="detail-label">IP:</span>
                                            <span className="detail-value">{asset.ip_address}</span>
                                        </div>

                                        <div className="detail-row">
                                            <span className="detail-label">Type:</span>
                                            <span className="detail-value">{asset.asset_type || 'unknown'}</span>
                                        </div>

                                        {asset.operating_system && (
                                            <div className="detail-row">
                                                <span className="detail-label">OS:</span>
                                                <span className="detail-value">{asset.operating_system}</span>
                                            </div>
                                        )}

                                        <div className="detail-row">
                                            <span className="detail-label">Location:</span>
                                            <span className="detail-value">{asset.location || 'N/A'}</span>
                                        </div>

                                        {asset.last_seen && (
                                            <div className="detail-row">
                                                <span className="detail-label">Last Seen:</span>
                                                <span className="detail-value">{formatDate(asset.last_seen)}</span>
                                            </div>
                                        )}
                                    </div>
                                </div>
                            ))
                        ) : (
                            <div className="no-results">
                                <p>No assets found matching your criteria</p>
                            </div>
                        )}
                    </div>
                )}
            </div>

        </div>
    );
};

export default Dashboard;