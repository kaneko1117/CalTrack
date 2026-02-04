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
  | "INTERNAL_ERROR"
  | "IMAGE_ANALYSIS_FAILED"
  | "INVALID_IMAGE_FORMAT"
  | "IMAGE_TOO_LARGE"
  | "NO_FOOD_DETECTED";

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
