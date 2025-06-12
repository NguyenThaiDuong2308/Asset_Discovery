import React, { useState, useEffect } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { getService, createService, updateService } from '../../services/serviceService';

const ServiceForm = () => {
    const { ip, serviceId } = useParams();
    const navigate = useNavigate();
    const [loading, setLoading] = useState(serviceId ? true : false);
    const [error, setError] = useState(null);

    const [formData, setFormData] = useState({
        name: '',
        port: '',
        protocol: 'TCP',
        description: '',
        is_managed: false
    });

    useEffect(() => {
        const fetchService = async () => {
            try {
                const data = await getService(ip, serviceId);
                setFormData(data);
                setLoading(false);
            } catch (err) {
                setError('Failed to fetch service details');
                setLoading(false);
            }
        };

        if (serviceId) {
            fetchService();
        }
    }, [ip, serviceId]);

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
            if (serviceId) {
                await updateService(ip, serviceId, formData);
            } else {
                await createService(ip, formData);
            }
            navigate(`/assets/${ip}/services`);
        } catch (err) {
            setError('Failed to save service');
        }
    };

    if (loading) return <div className="loader">Loading...</div>;

    return (
        <div className="service-form-container">
            <h1>{serviceId ? 'Edit Service' : 'Add New Service'}</h1>
            <h2>Asset IP: {ip}</h2>

            {error && <div className="error-message">{error}</div>}

            <form onSubmit={handleSubmit} className="service-form">
                <div className="form-group">
                    <label htmlFor="name">Service Name</label>
                    <input
                        type="text"
                        id="name"
                        name="name"
                        value={formData.name}
                        onChange={handleChange}
                        required
                    />
                </div>

                <div className="form-group">
                    <label htmlFor="port">Port</label>
                    <input
                        type="number"
                        id="port"
                        name="port"
                        value={formData.port}
                        onChange={handleChange}
                        required
                        min="1"
                        max="65535"
                    />
                </div>

                <div className="form-group">
                    <label htmlFor="protocol">Protocol</label>
                    <select
                        id="protocol"
                        name="protocol"
                        value={formData.protocol}
                        onChange={handleChange}
                        required
                    >
                        <option value="TCP">TCP</option>
                        <option value="UDP">UDP</option>
                    </select>
                </div>

                <div className="form-group">
                    <label htmlFor="description">Description</label>
                    <textarea
                        id="description"
                        name="description"
                        value={formData.description}
                        onChange={handleChange}
                        rows="3"
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
                    <label htmlFor="is_managed">Managed Service</label>
                </div>

                <div className="form-actions">
                    <button type="submit" className="btn btn-primary">Save</button>
                    <button
                        type="button"
                        className="btn btn-secondary"
                        onClick={() => navigate(`/assets/${ip}/services`)}
                    >
                        Cancel
                    </button>
                </div>
            </form>
        </div>
    );
};

export default ServiceForm;