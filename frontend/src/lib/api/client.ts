import { PUBLIC_API_URL } from '$env/static/public';

export class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public data?: unknown
  ) {
    super(message);
    this.name = 'ApiError';
  }
}

interface RequestOptions extends RequestInit {
  fetch?: typeof fetch;
}

async function apiRequest<T>(
  endpoint: string,
  options: RequestOptions = {}
): Promise<T> {
  const { fetch: customFetch, ...fetchOptions } = options;
  const fetchFn = customFetch || fetch;

  const headers = {
    'Content-Type': 'application/json',
    ...(fetchOptions.headers as Record<string, string>),
  };

  const response = await fetchFn(`${PUBLIC_API_URL}${endpoint}`, {
    ...fetchOptions,
    headers,
  });

  if (!response.ok) {
    const errorData = await response.json().catch(() => ({}));
    throw new ApiError(
      errorData.message || `API Error: ${response.statusText}`,
      response.status,
      errorData
    );
  }

  if (response.status === 204 || response.headers.get('content-length') === '0') {
    return undefined as T;
  }

  return response.json();
}

export const api = {
  get: <T>(endpoint: string, customFetch?: typeof fetch) =>
    apiRequest<T>(endpoint, { method: 'GET', fetch: customFetch }),

  post: <T>(endpoint: string, body?: unknown, customFetch?: typeof fetch) =>
    apiRequest<T>(endpoint, { method: 'POST', body: JSON.stringify(body), fetch: customFetch }),

  put: <T>(endpoint: string, body?: unknown, customFetch?: typeof fetch) =>
    apiRequest<T>(endpoint, { method: 'PUT', body: JSON.stringify(body), fetch: customFetch }),

  patch: <T>(endpoint: string, body?: unknown, customFetch?: typeof fetch) =>
    apiRequest<T>(endpoint, { method: 'PATCH', body: JSON.stringify(body), fetch: customFetch }),

  delete: <T = void>(endpoint: string, customFetch?: typeof fetch) =>
    apiRequest<T>(endpoint, { method: 'DELETE', fetch: customFetch }),
};