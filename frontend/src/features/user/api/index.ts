/**
 * User API 型定義
 * ユーザープロフィールの取得・更新に関するAPI型
 */
import type { ActivityLevelValue } from "@/domain/valueObjects";

/** プロフィール更新リクエスト */
export type UpdateProfileRequest = {
  nickname: string;
  height: number;
  weight: number;
  activityLevel: string;
};

/** プロフィール更新レスポンス */
export type UpdateProfileResponse = {
  userId: string;
  nickname: string;
  height: number;
  weight: number;
  activityLevel: string;
};

/** 現在のユーザー情報レスポンス */
export type CurrentUserResponse = {
  email: string;
  nickname: string;
  weight: number;
  height: number;
  birthDate: string;
  gender: string;
  activityLevel: ActivityLevelValue;
};
