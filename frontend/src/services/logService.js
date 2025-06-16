import { fetchApi, LOG_SERVICE_URL } from './api';

// Get all logs
export const getLogs = async () => {
    return fetchApi(`${LOG_SERVICE_URL}/api/logs`);
};

// Get assets from log service
export const getLogAnalysisAssets = async () => {
    return fetchApi(`${LOG_SERVICE_URL}/api/assets`);
};

// Get services from log service
export const getLogAnalysisServices = async () => {
    return fetchApi(`${LOG_SERVICE_URL}/api/services`);
};

// Get services for a specific IP from log service
export const getLogAnalysisServicesByIp = async (ip) => {
    return fetchApi(`${LOG_SERVICE_URL}/api/services/${ip}`);
};

// Get specific asset details from log service
export const getLogServiceAssetByIp = async (ip) => {
    return fetchApi(`${LOG_SERVICE_URL}/api/assets/${ip}`);
};