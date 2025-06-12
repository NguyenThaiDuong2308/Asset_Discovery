import React from 'react';
import { Link } from 'react-router-dom';

const AssetCard = ({ asset }) => {
    return (
        <div className="asset-card">
            <div className="card-header">
                <h3>{asset.hostname || 'Unnamed Asset'}</h3>
                <span className={`badge ${asset.is_managed ? 'badge-success' : 'badge-warning'}`}>
          {asset.is_managed ? 'Managed' : 'Unmanaged'}
        </span>
            </div>
            <div className="card-body">
                <p><strong>IP:</strong> {asset.ip_address}</p>
                <p><strong>Type:</strong> {asset.asset_type}</p>
                <p><strong>Location:</strong> {asset.location}</p>
                <p><strong>OS:</strong> {asset.operating_system}</p>
                <p><strong>Last Seen:</strong> {new Date(asset.last_seen).toLocaleString()}</p>
            </div>
            <div className="card-footer">
                <Link to={`/assets/${asset.ip_address}`} className="btn btn-info">Details</Link>
                <Link to={`/assets/${asset.ip_address}/services`} className="btn btn-secondary">Services</Link>
            </div>
        </div>
    );
};

export default AssetCard;