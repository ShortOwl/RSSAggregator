import { clearSession, getToken } from "./session";
export class ApiError extends Error {
  constructor(
    message: string,
    public status: number,
    public retryAfter = 1,
  ) {
    super(message);
    this.name = "ApiError";
  }
}
export const API_URL = (
  process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080/v1"
).replace(/\/$/, "");
export function queryString(params: object) {
  const query = new URLSearchParams();
  Object.entries(params).forEach(([key, value]) => {
    if (value !== "" && value !== undefined && value !== false)
      query.set(key, String(value));
  });
  return query.toString();
}
export async function api<T>(
  path: string,
  options: RequestInit = {},
  authenticated = true,
): Promise<T> {
  const token = authenticated ? getToken() : null;
  if (authenticated && !token) {
    clearSession();
    throw new ApiError("Your session has ended. Please sign in again.", 401);
  }
  const headers = new Headers(options.headers);
  if (options.body) headers.set("Content-Type", "application/json");
  if (token) headers.set("Authorization", `Bearer ${token}`);
  let response: Response;
  try {
    response = await fetch(`${API_URL}${path}`, {
      ...options,
      headers,
      cache: "no-store",
    });
  } catch (error) {
    if (options.signal?.aborted) throw error;
    throw new ApiError(
      "Cannot reach your RSS server. Check your connection and try again.",
      0,
    );
  }
  if (!response.ok) {
    const payload = await response.json().catch(() => ({}));
    if (response.status === 401 && authenticated && getToken() === token)
      clearSession();
    throw new ApiError(
      payload.error || "Something went wrong. Please try again.",
      response.status,
      Number(response.headers.get("Retry-After")) || 1,
    );
  }
  if (response.status === 204) return undefined as T;
  return response.json() as Promise<T>;
}
