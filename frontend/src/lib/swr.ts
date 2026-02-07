/**
 * SWR設定
 */
import { AxiosError } from "axios";
import { apiClient, type ApiErrorResponse } from "./api";

/** GET用fetcher */
export const fetcher = async <T>(url: string): Promise<T> => {
  try {
    const res = await apiClient.get<T>(url);
    return res.data;
  } catch (e) {
    if (e instanceof AxiosError && e.response?.data) {
      throw e.response.data as ApiErrorResponse;
    }
    throw { code: "INTERNAL_ERROR", message: "予期しないエラーが発生しました" };
  }
};

type MutationMethod = "POST" | "PUT" | "PATCH" | "DELETE";

/** Mutation用fetcher */
export const mutate = async <T>(
  url: string,
  { arg }: { arg: { method: MutationMethod; data?: unknown } }
): Promise<T> => {
  try {
    const res = await apiClient.request<T>({ url, method: arg.method, data: arg.data });
    return res.data;
  } catch (e) {
    if (e instanceof AxiosError && e.response?.data) {
      throw e.response.data as ApiErrorResponse;
    }
    throw { code: "INTERNAL_ERROR", message: "予期しないエラーが発生しました" };
  }
};
