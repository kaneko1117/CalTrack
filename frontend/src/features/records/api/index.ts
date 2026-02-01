/**
 * カロリー記録API
 * カロリー記録に関するAPI関数定義
 */

import { AxiosError } from "axios";
import { apiClient } from "@/lib/api";
import type {
  CreateRecordRequest,
  CreateRecordResponse,
} from "../types";
import type { ApiErrorResponse, ErrorCode } from "@/features/common";

/** エラーコード定数 */
const ERROR_CODE_INTERNAL_ERROR: ErrorCode = "INTERNAL_ERROR";

/** エラーメッセージ定数 */
const ERROR_MESSAGE_UNEXPECTED = "予期しないエラーが発生しました";

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
 * カロリー記録を作成
 * @param request - カロリー記録作成データ
 * @returns Promise<CreateRecordResponse> - 作成された記録
 * @throws ApiError - バリデーションエラー、認証エラー、サーバーエラー時
 */
export async function createRecord(
  request: CreateRecordRequest
): Promise<CreateRecordResponse> {
  try {
    const response = await apiClient.post<CreateRecordResponse>(
      "/api/v1/records",
      request,
      {
        withCredentials: true,
      }
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
