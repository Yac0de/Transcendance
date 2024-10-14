export const API_BASE_URL = 'http://localhost:4000';

export async function apiRequest<T>(url: string, options: RequestInit): Promise<T> {
    const response = await fetch(`${API_BASE_URL}${url}`, options);
    const result = await response.json();
    
    if (!response.ok) {
        throw new Error(result.error || 'Request failed');
    }
    return result;
}
