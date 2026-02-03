/**
 * API Client Configuration
 */
import axios, { AxiosError, AxiosResponse } from "axios";
import { router } from "@/routes";

/** ログイン画面のパス */
const LOGIN_PATH = "/";

/** 401エラーのステータスコード */
const HTTP_STATUS_UNAUTHORIZED = 401;

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
 * 401エラー時にログイン画面へリダイレクト
 */
function redirectToLogin(): void {
  if (router.state.location.pathname !== LOGIN_PATH) {
    router.navigate(LOGIN_PATH);
  }
}

/**
 * レスポンスインターセプター
 */
apiClient.interceptors.response.use(
  (response: AxiosResponse) => response,
  (error: AxiosError) => {
    if (error.response?.status === HTTP_STATUS_UNAUTHORIZED) {
      redirectToLogin();
    }
    return Promise.reject(error);
  }
);

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

export const get = <T>(url: string) => request(apiClient.get<T>(url));
export const post = <T, D = unknown>(url: string, data?: D) =>
  request(apiClient.post<T>(url, data));
export const put = <T, D = unknown>(url: string, data?: D) =>
  request(apiClient.put<T>(url, data));
export const del = <T>(url: string) => request(apiClient.delete<T>(url));
