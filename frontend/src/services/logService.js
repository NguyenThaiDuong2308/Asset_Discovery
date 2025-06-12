import { fetchApi, LOG_SERVICE_URL } from './api';

// Get all logs
export const getLogs = async () => {
    return fetchApi(`${LOG_SERVICE_URL}/api/logs`);
};

// Get assets from log service
export const getLogServiceAssets = async () => {
    return fetchApi(`${LOG_SERVICE_URL}/api/assets`);
};

// Get services from log service
export const getLogServiceServices = async () => {
    return fetchApi(`${LOG_SERVICE_URL}/api/services`);
};

// Get services for a specific IP from log service
export const getLogServiceAssetServices = async (ip) => {
    return fetchApi(`${LOG_SERVICE_URL}/api/services/${ip}`);
};

// Get specific asset details from log service
export const getLogServiceAssetDetails = async (ip) => {
    return fetchApi(`${LOG_SERVICE_URL}/api/assets/${ip}`);
};