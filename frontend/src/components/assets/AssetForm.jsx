import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { getAssetByIp, createAsset, updateAsset } from '../../services/assetService';

const AssetForm = () => {
    const { ip } = useParams();
    const navigate = useNavigate();
    const [loading, setLoading] = useState(ip ? true : false);
    const [error, setError] = useState(null);

    const [formData, setFormData] = useState({
        ip_address: '',
        mac_address: '',
        hostname: '',
        asset_type: '',
        location: '',
        operating_system: '',
        first_seen: new Date().toISOString(),
        last_seen: new Date().toISOString(),
        is_managed: false
    });

    useEffect(() => {
        const fetchAsset = async () => {
            try {
                const data = await getAssetByIp(ip);
                setFormData(data);
                setLoading(false);
            } catch (err) {
                setError('Failed to fetch asset details');
                setLoading(false);
            }
        };

        if (ip) {
            fetchAsset();
        }
    }, [ip]);

    const handleChange = (e) => {
        const { name, value, type, checked } = e.target;
        setFormData({
            ...formData,
            [name]: type === 'checkbox' ? checked : value
        });
    };

    const handleSubmit = async (e) => {
        e.preventDefault();
        try {
            if (ip) {
                await updateAsset(ip, formData);
            } else {
                await createAsset(formData);
            }
            navigate(ip ? `/assets/${ip}` : '/assets');
        } catch (err) {
            setError('Failed to save asset');
        }
    };

    if (loading) return <div className="loader">Loading...</div>;

    return (
        <div className="asset-form-container">
            <h1>{ip ? 'Edit Asset' : 'Add New Asset'}</h1>

            {error && <div className="error-message">{error}</div>}

            <form onSubmit={handleSubmit} className="asset-form">
                <div className="form-group">
                    <label htmlFor="ip_address">IP Address</label>
                    <input
                        type="text"
                        id="ip_address"
                        name="ip_address"
                        value={formData.ip_address}
                        onChange={handleChange}
                        required
                        readOnly={!!ip}
                    />
                </div>

                <div className="form-group">
                    <label htmlFor="mac_address">MAC Address</label>
                    <input
                        type="text"
                        id="mac_address"
                        name="mac_address"
                        value={formData.mac_address}
                        onChange={handleChange}
                    />
                </div>

                <div className="form-group">
                    <label htmlFor="hostname">Hostname</label>
                    <input
                        type="text"
                        id="hostname"
                        name="hostname"
                        value={formData.hostname}
                        onChange={handleChange}
                    />
                </div>

                <div className="form-group">
                    <label htmlFor="asset_type">Asset Type</label>
                    <select
                        id="asset_type"
                        name="asset_type"
                        value={formData.asset_type}
                        onChange={handleChange}
                        required
                    >
                        <option value="">Select type</option>
                        <option value="server">Server</option>
                        <option value="workstation">Workstation</option>
                        <option value="network">Network Device</option>
                        <option value="iot">IoT Device</option>
                        <option value="other">Other</option>
                    </select>
                </div>

                <div className="form-group">
                    <label htmlFor="location">Location</label>
                    <input
                        type="text"
                        id="location"
                        name="location"
                        value={formData.location}
                        onChange={handleChange}
                    />
                </div>

                <div className="form-group">
                    <label htmlFor="operating_system">Operating System</label>
                    <input
                        type="text"
                        id="operating_system"
                        name="operating_system"
                        value={formData.operating_system}
                        onChange={handleChange}
                    />
                </div>

                <div className="form-group checkbox">
                    <input
                        type="checkbox"
                        id="is_managed"
                        name="is_managed"
                        checked={formData.is_managed}
                        onChange={handleChange}
                    />
                    <label htmlFor="is_managed">Managed Asset</label>
                </div>

                <div className="form-actions">
                    <button type="submit" className="btn btn-primary">Save</button>
                    <button
                        type="button"
                        className="btn btn-secondary"
                        onClick={() => navigate(ip ? `/assets/${ip}` : '/assets')}
                    >
                        Cancel
                    </button>
                </div>
            </form>
        </div>
    );
};

export default AssetForm;