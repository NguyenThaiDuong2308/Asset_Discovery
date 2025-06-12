import React from 'react';
import { Link } from 'react-router-dom';
import { deleteService, setServiceManaged } from '../../services/serviceService';

const ServiceCard = ({ service, assetIp }) => {
    const [error, setError] = React.useState(null);
    const [isManaged, setIsManaged] = React.useState(service.is_managed);

    const handleDelete = async () => {
        if (window.confirm('Are you sure you want to delete this service?')) {
            try {
                await deleteService(assetIp, service.id);
                // Remove from parent list or refresh
                window.location.reload();
            } catch (err) {
                setError('Failed to delete service');
            }
        }
    };

    const handleManage = async () => {
        try {
            await setServiceManaged(assetIp, service.id);
            setIsManaged(true);
        } catch (err) {
            setError('Failed to set service as managed');
        }
    };

    return (
        <div className="service-card">
            {error && <div className="error-message">{error}</div>}

            <div className="card-header">
                <h3>{service.name}</h3>
                <span className={`badge ${isManaged ? 'badge-success' : 'badge-warning'}`}>
          {isManaged ? 'Managed' : 'Unmanaged'}
        </span>
            </div>

            <div className="card-body">
                <p><strong>Port:</strong> {service.port}</p>
                <p><strong>Protocol:</strong> {service.protocol}</p>
                <p><strong>Description:</strong> {service.description}</p>
            </div>

            <div className="card-footer">
                <Link
                    to={`/assets/${assetIp}/services/${service.id}/edit`}
                    className="btn btn-secondary"
                >
                    Edit
                </Link>

                {!isManaged && (
                    <button onClick={handleManage} className="btn btn-success">
                        Set as Managed
                    </button>
                )}

                <button onClick={handleDelete} className="btn btn-danger">
                    Delete
                </button>
            </div>
        </div>
    );
};

export default ServiceCard;