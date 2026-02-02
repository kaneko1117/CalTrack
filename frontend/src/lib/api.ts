/**
 * API Client Configuration
 * Shared axios instance for all API calls
 */
import axios, { AxiosError, AxiosResponse } from "axios";

/** バックエンドからのエラーコード */
export type ErrorCode =
  | "INVALID_REQUEST"
  | "VALIDATION_ERROR"
  | "UNAUTHORIZED"
  | "EMAIL_ALREADY_EXISTS"
  | "INVALID_CREDENTIALS"
  | "INTERNAL_ERROR";

/** APIエラーレスポンス */
export type ApiErrorResponse = {
  code: ErrorCode;
  message: string;
  details?: string[];
};

export const apiClient = axios.create({
  baseURL: "http://localhost:8080",
  timeout: 10000,
  headers: {
    "Content-Type": "application/json",
  },
  withCredentials: true,
});

/**
 * APIエラーをApiErrorResponse形式に変換
 */
function handleApiError(error: unknown): ApiErrorResponse {
  if (error instanceof AxiosError && error.response?.data) {
    return error.response.data as ApiErrorResponse;
  }
  return {
    code: "INTERNAL_ERROR",
    message: "予期しないエラーが発生しました",
  };
}

/**
 * request - 汎用APIリクエスト
 * @param request - apiClientのメソッド呼び出し
 * @returns Promise<T>
 */
export async function request<T>(
  request: Promise<AxiosResponse<T>>
): Promise<T> {
  try {
    const response = await request;
    return response.data;
  } catch (error) {
    throw handleApiError(error);
  }
}

/**
 * GET
 */
export const get = <T>(url: string) => request(apiClient.get<T>(url));

/**
 * POST
 */
export const post = <T, D = unknown>(url: string, data?: D) =>
  request(apiClient.post<T>(url, data));

/**
 * PUT
 */
export const put = <T, D = unknown>(url: string, data?: D) =>
  request(apiClient.put<T>(url, data));

/**
 * DELETE
 */
export const del = <T>(url: string) => request(apiClient.delete<T>(url));
