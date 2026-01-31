/**
 * 認証API
 * ユーザー登録APIの関数定義
 */

import { AxiosError } from "axios";
import { apiClient } from "@/lib/api";
import type {
  RegisterUserRequest,
  RegisterUserResponse,
  ApiErrorResponse,
  ErrorCode,
} from "../types";
import {
  ERROR_CODE_INTERNAL_ERROR,
  ERROR_MESSAGE_UNEXPECTED,
} from "../types";

/**
 * カスタムAPIエラークラス
 */
export class ApiError extends Error {
  readonly code: ErrorCode;
  readonly details?: string[];
  readonly status: number;

  constructor(
    code: ErrorCode,
    message: string,
    status: number,
    details?: string[]
  ) {
    super(message);
    this.name = "ApiError";
    this.code = code;
    this.status = status;
    this.details = details;
  }

  static fromResponse(response: ApiErrorResponse, status: number): ApiError {
    return new ApiError(response.code, response.message, status, response.details);
  }
}

/**
 * 新規ユーザー登録
 * @param request - ユーザー登録データ
 * @returns Promise<RegisterUserResponse> - 作成されたユーザーID
 * @throws ApiError - バリデーションエラー、メール重複、サーバーエラー時
 */
export async function registerUser(
  request: RegisterUserRequest
): Promise<RegisterUserResponse> {
  try {
    const response = await apiClient.post<RegisterUserResponse>(
      "/api/v1/users",
      request
    );
    return response.data;
  } catch (error) {
    if (error instanceof AxiosError && error.response) {
      const errorData = error.response.data as ApiErrorResponse;
      throw ApiError.fromResponse(errorData, error.response.status);
    }
    // ネットワークエラーまたは予期しないエラー
    throw new ApiError(
      ERROR_CODE_INTERNAL_ERROR,
      ERROR_MESSAGE_UNEXPECTED,
      500
    );
  }
}
