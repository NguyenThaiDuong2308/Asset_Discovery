// Base URLs for the services
const ASSET_SERVICE_URL = 'http://localhost:8081';
const LOG_SERVICE_URL = 'http://localhost:8080';

// Helper for making API requests
export const fetchApi = async (url, options = {}) => {
    const response = await fetch(url, {
        ...options,
        headers: {
            'Content-Type': 'application/json',
            ...options.headers,
        },
    });

    if (!response.ok) {
        const error = await response.json();
        throw new Error(error.error || 'Something went wrong');
    }

    return response.json();
};

export { ASSET_SERVICE_URL, LOG_SERVICE_URL };