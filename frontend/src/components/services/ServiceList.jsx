import React, { useState, useEffect } from 'react';
import { useParams, Link } from 'react-router-dom';
import { getServices } from '../../services/serviceService';
import ServiceCard from './ServiceCard';

const ServiceList = () => {
    const { ip } = useParams();
    const [services, setServices] = useState([]);
    const [loading, setLoading] = useState(true);
    const [error, setError] = useState(null);

    useEffect(() => {
        const fetchServices = async () => {
            try {
                const data = await getServices(ip);

                if (Array.isArray(data)) {
                    setServices(data);
                } else {
                    console.warn("getServices() response is not an array:", data);
                    setServices([]);
                }

                setLoading(false);
            } catch (err) {
                console.error("Failed to fetch services:", err);
                setError('Failed to fetch services');
                setServices([]);
                setLoading(false);
            }
        };

        fetchServices();
    }, [ip]);

    if (loading) return <div className="loader">Loading...</div>;
    if (error) return <div className="error-message">{error}</div>;

    return (
        <div className="service-list-container">
            <div className="list-header">
                <h1>Services for {ip}</h1>
                <div>
                    <Link to={`/assets/${ip}`} className="btn btn-secondary">Back to Asset</Link>
                    <Link to={`/assets/${ip}/services/new`} className="btn btn-primary">Add New Service</Link>
                </div>
            </div>

            {!services || services.length === 0 ? (
                <p>No services found for this asset. Add one to get started.</p>
            ) : (
                <div className="service-grid">
                    {services.map(service => (
                        <ServiceCard
                            key={service.id}
                            service={service}
                            assetIp={ip}
                        />
                    ))}
                </div>
            )}
        </div>
    );
};

export default ServiceList;
