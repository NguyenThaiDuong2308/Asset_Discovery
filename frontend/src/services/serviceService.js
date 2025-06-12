import { fetchApi, ASSET_SERVICE_URL } from './api';

// Get all services for an asset
export const getServices = async (ip) => {
    return fetchApi(`${ASSET_SERVICE_URL}/api/assets/${ip}/services`);
};

// Get a specific service
export const getService = async (ip, serviceId) => {
    // This isn't in the API docs, but we'd need this endpoint
    // For now, let's assume we'd find it in the list
    const services = await getServices(ip);
    return services.find(service => service.id === parseInt(serviceId));
};

// Create a new service for an asset
export const createService = async (ip, serviceData) => {
    return fetchApi(`${ASSET_SERVICE_URL}/api/assets/${ip}/services`, {
        method: 'POST',
        body: JSON.stringify(serviceData),
    });
};

// Update an existing service
export const updateService = async (ip, serviceId, serviceData) => {
    return fetchApi(`${ASSET_SERVICE_URL}/api/assets/${ip}/services/${serviceId}`, {
        method: 'PUT',
        body: JSON.stringify(serviceData),
    });
};

// Delete a service
export const deleteService = async (ip, serviceId) => {
    return fetchApi(`${ASSET_SERVICE_URL}/api/assets/${ip}/services/${serviceId}`, {
        method: 'DELETE',
    });
};

// Set a service as managed
export const setServiceManaged = async (ip, serviceId) => {
    return fetchApi(`${ASSET_SERVICE_URL}/api/assets/${ip}/services/${serviceId}/manage`, {
        method: 'PATCH',
    });
};