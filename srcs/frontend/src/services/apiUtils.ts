import { getBaseHost } from '../utils/fetch'

export const API_BASE_URL = `${getBaseHost()}${import.meta.env.PROD ? '/api' : ':4000'}`;

//export const API_BASE_URL = import.meta.env.PROD 
//  ? 'https://localhost:8443/api'
//  : 'http://localhost:4000' 

export async function apiRequest<T>(url: string, options: RequestInit): Promise<T> {
    const response = await fetch(`${API_BASE_URL}${url}`, options);
    let result;
    
    try {
        result = await response.json();
    } catch (error) {
        throw new Error("Failed to parse JSON response.");
    }

    if (!response.ok) {
        throw { error: result.error || 'Request failed', status: response.status };
    }
    result.status = response.status;
    return result;
}

