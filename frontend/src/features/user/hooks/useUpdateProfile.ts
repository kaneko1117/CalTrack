/**
 * useUpdateProfile - プロフィール更新フック
 * プロフィールの更新処理を行う
 */
import { useMemo } from "react";
import { useForm } from "@/features/common/hooks";
import {
  newNickname,
  newHeight,
  newWeight,
  newActivityLevel,
} from "@/domain/valueObjects";
import type { UpdateProfileResponse, CurrentUserResponse, UpdateProfileRequest } from "../api";

/** フォームフィールド型 */
type ProfileField = "nickname" | "height" | "weight" | "activityLevel";

/** VOファクトリ設定 */
const formConfig = {
  nickname: newNickname,
  height: newHeight,
  weight: newWeight,
  activityLevel: newActivityLevel,
};

/** 空のフォーム状態 */
const emptyFormState: Record<ProfileField, string> = {
  nickname: "",
  height: "",
  weight: "",
  activityLevel: "",
};

/** エラーの初期状態 */
const initialErrors: Record<ProfileField, string | null> = {
  nickname: null,
  height: null,
  weight: null,
  activityLevel: null,
};

/**
 * currentUserから初期フォーム状態を構築
 */
function buildFormState(user: CurrentUserResponse | undefined): Record<ProfileField, string> {
  if (!user) return emptyFormState;
  return {
    nickname: user.nickname,
    height: String(user.height),
    weight: String(user.weight),
    activityLevel: user.activityLevel,
  };
}

/**
 * useUpdateProfile - プロフィール更新フォーム
 * @param currentUser - 現在のユーザー情報
 * @param onSuccess - 成功時のコールバック
 */
export function useUpdateProfile(
  currentUser: CurrentUserResponse | undefined,
  onSuccess?: (result: UpdateProfileResponse) => void
) {
  const initialFormState = useMemo(
    () => buildFormState(currentUser),
    [currentUser],
  );

  return useForm<ProfileField, UpdateProfileResponse, UpdateProfileRequest>({
    config: formConfig,
    initialFormState,
    initialErrors,
    url: "/api/v1/users/profile",
    method: "PATCH",
    transformData: (data) => ({
      nickname: data.nickname,
      height: parseFloat(data.height),
      weight: parseFloat(data.weight),
      activityLevel: data.activityLevel,
    }),
    onSuccess,
  });
}
