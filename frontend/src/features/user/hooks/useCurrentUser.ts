/**
 * useCurrentUser - プロフィール取得フック
 * 現在のユーザー情報を取得する
 */
import { useRequestGet } from "@/features/common/hooks";
import type { CurrentUserResponse } from "../api";

/**
 * useCurrentUser - プロフィール取得
 */
export function useCurrentUser() {
  const { data, error, isLoading, mutate } = useRequestGet<CurrentUserResponse>(
    "/api/v1/users/profile"
  );

  return {
    data,
    error,
    isLoading,
    refetch: mutate,
  };
}
