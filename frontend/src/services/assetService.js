import { fetchApi, ASSET_SERVICE_URL } from './api';

// Get all assets
export const getAssets = async () => {
    return fetchApi(`${ASSET_SERVICE_URL}/api/assets`);
};

// Get a specific asset by IP
export const getAssetByIp = async (ip) => {
    return fetchApi(`${ASSET_SERVICE_URL}/api/assets/${ip}`);
};

// Create a new asset
export const createAsset = async (assetData) => {
    return fetchApi(`${ASSET_SERVICE_URL}/api/assets`, {
        method: 'POST',
        body: JSON.stringify(assetData),
    });
};

// Update an existing asset
export const updateAsset = async (ip, assetData) => {
    return fetchApi(`${ASSET_SERVICE_URL}/api/assets/${ip}`, {
        method: 'PUT',
        body: JSON.stringify(assetData),
    });
};

// Delete an asset
export const deleteAsset = async (ip) => {
    return fetchApi(`${ASSET_SERVICE_URL}/api/assets/${ip}`, {
        method: 'DELETE',
    });
};

// Set an asset as managed
export const setAssetManaged = async (ip) => {
    return fetchApi(`${ASSET_SERVICE_URL}/api/assets/${ip}/manage`, {
        method: 'PATCH',
    });
};