import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, act } from "@testing-library/react";
import { useUpdateProfile } from "./useUpdateProfile";
import type { CurrentUserResponse, UpdateProfileResponse } from "../api";

// SWR mutationをモック
const mockTrigger = vi.fn();
const mockReset = vi.fn();
const mockError = undefined;

vi.mock("@/features/common/hooks/useRequest", () => ({
  useRequestMutation: () => ({
    trigger: mockTrigger,
    isMutating: false,
    error: mockError,
    data: undefined,
    reset: mockReset,
  }),
}));

// VOファクトリをモック
vi.mock("@/domain/valueObjects", async () => {
  const actual = await vi.importActual("@/domain/valueObjects");
  return actual;
});

describe("useUpdateProfile", () => {
  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("初期化", () => {
    it("currentUserがundefinedの場合は空フォームで初期化される", () => {
      const { result } = renderHook(() => useUpdateProfile(undefined));

      expect(result.current.formState).toEqual({
        nickname: "",
        height: "",
        weight: "",
        activityLevel: "",
      });
      expect(result.current.errors).toEqual({
        nickname: null,
        height: null,
        weight: null,
        activityLevel: null,
      });
      expect(result.current.isValid).toBe(false);
    });

    it("currentUserがある場合はユーザー情報で初期化される", () => {
      const mockUser: CurrentUserResponse = {
        email: "test@example.com",
        nickname: "テストユーザー",
        weight: 70,
        height: 175,
        birthDate: "1990-01-01",
        gender: "male",
        activityLevel: "moderate",
      };

      const { result } = renderHook(() => useUpdateProfile(mockUser));

      expect(result.current.formState).toEqual({
        nickname: "テストユーザー",
        height: "175",
        weight: "70",
        activityLevel: "moderate",
      });
    });
  });

  describe("currentUser変更時の同期", () => {
    it("currentUserがundefinedからデータに変化したらフォーム状態が同期される", () => {
      const { result, rerender } = renderHook(
        ({ user }: { user: CurrentUserResponse | undefined }) => useUpdateProfile(user),
        { initialProps: { user: undefined as CurrentUserResponse | undefined } }
      );

      // 初期状態は空
      expect(result.current.formState.nickname).toBe("");

      // ユーザーデータを設定
      const mockUser: CurrentUserResponse = {
        email: "test@example.com",
        nickname: "新ユーザー",
        weight: 60,
        height: 160,
        birthDate: "1995-05-15",
        gender: "female",
        activityLevel: "light",
      };

      rerender({ user: mockUser });

      // フォーム状態が同期される
      expect(result.current.formState).toEqual({
        nickname: "新ユーザー",
        height: "160",
        weight: "60",
        activityLevel: "light",
      });
    });

    it("別のユーザーに変化したらフォーム状態が同期される", () => {
      const user1: CurrentUserResponse = {
        email: "user1@example.com",
        nickname: "ユーザー1",
        weight: 70,
        height: 175,
        birthDate: "1990-01-01",
        gender: "male",
        activityLevel: "moderate",
      };

      const { result, rerender } = renderHook(
        ({ user }: { user: CurrentUserResponse | undefined }) => useUpdateProfile(user),
        { initialProps: { user: user1 as CurrentUserResponse | undefined } }
      );

      expect(result.current.formState.nickname).toBe("ユーザー1");

      // 別ユーザーに変更
      const user2: CurrentUserResponse = {
        email: "user2@example.com",
        nickname: "ユーザー2",
        weight: 65,
        height: 170,
        birthDate: "1992-03-20",
        gender: "female",
        activityLevel: "veryActive",
      };

      rerender({ user: user2 });

      // フォーム状態が更新される
      expect(result.current.formState).toEqual({
        nickname: "ユーザー2",
        height: "170",
        weight: "65",
        activityLevel: "veryActive",
      });
    });

    it("同じ参照で再レンダーしても編集内容が維持される", () => {
      const mockUser: CurrentUserResponse = {
        email: "test@example.com",
        nickname: "テストユーザー",
        weight: 70,
        height: 175,
        birthDate: "1990-01-01",
        gender: "male",
        activityLevel: "moderate",
      };

      const { result, rerender } = renderHook(
        ({ user }: { user: CurrentUserResponse | undefined }) => useUpdateProfile(user),
        { initialProps: { user: mockUser as CurrentUserResponse | undefined } }
      );

      // ユーザーが編集
      act(() => {
        result.current.handleChange("nickname")("編集後");
      });

      expect(result.current.formState.nickname).toBe("編集後");

      // 同じ参照で再レンダー
      rerender({ user: mockUser });

      // 編集内容が維持される
      expect(result.current.formState.nickname).toBe("編集後");
    });
  });

  describe("バリデーション", () => {
    it("currentUserのデータが有効な場合isValid=trueになる", () => {
      const mockUser: CurrentUserResponse = {
        email: "test@example.com",
        nickname: "テストユーザー",
        weight: 70,
        height: 175,
        birthDate: "1990-01-01",
        gender: "male",
        activityLevel: "moderate",
      };

      const { result } = renderHook(() => useUpdateProfile(mockUser));

      expect(result.current.isValid).toBe(true);
    });

    it("フィールド変更時にバリデーションが実行される", () => {
      const mockUser: CurrentUserResponse = {
        email: "test@example.com",
        nickname: "テストユーザー",
        weight: 70,
        height: 175,
        birthDate: "1990-01-01",
        gender: "male",
        activityLevel: "moderate",
      };

      const { result } = renderHook(() => useUpdateProfile(mockUser));

      // 無効な体重を入力
      act(() => {
        result.current.handleChange("weight")("999");
      });

      // エラーが設定される
      expect(result.current.errors.weight).not.toBeNull();
      expect(result.current.isValid).toBe(false);
    });
  });

  describe("送信", () => {
    it("handleSubmitでAPIリクエストが送信される", async () => {
      const mockUser: CurrentUserResponse = {
        email: "test@example.com",
        nickname: "テストユーザー",
        weight: 70,
        height: 175,
        birthDate: "1990-01-01",
        gender: "male",
        activityLevel: "moderate",
      };

      const mockResponse: UpdateProfileResponse = {
        userId: "user-123",
        nickname: "更新後",
        weight: 75,
        height: 175,
        activityLevel: "moderate",
      };

      mockTrigger.mockResolvedValue(mockResponse);

      const { result } = renderHook(() => useUpdateProfile(mockUser));

      // ニックネームを変更
      act(() => {
        result.current.handleChange("nickname")("更新後");
        result.current.handleChange("weight")("75");
      });

      const mockEvent = {
        preventDefault: vi.fn(),
      } as unknown as React.FormEvent;

      await act(async () => {
        await result.current.handleSubmit(mockEvent);
      });

      expect(mockEvent.preventDefault).toHaveBeenCalled();
      expect(mockTrigger).toHaveBeenCalledWith({
        nickname: "更新後",
        height: 175,
        weight: 75,
        activityLevel: "moderate",
      });
    });

    it("成功時にonSuccessが呼び出される", async () => {
      const mockUser: CurrentUserResponse = {
        email: "test@example.com",
        nickname: "テストユーザー",
        weight: 70,
        height: 175,
        birthDate: "1990-01-01",
        gender: "male",
        activityLevel: "moderate",
      };

      const mockResponse: UpdateProfileResponse = {
        userId: "user-123",
        nickname: "更新後",
        weight: 70,
        height: 175,
        activityLevel: "moderate",
      };

      mockTrigger.mockResolvedValue(mockResponse);

      const onSuccess = vi.fn();
      const { result } = renderHook(() => useUpdateProfile(mockUser, onSuccess));

      const mockEvent = {
        preventDefault: vi.fn(),
      } as unknown as React.FormEvent;

      await act(async () => {
        await result.current.handleSubmit(mockEvent);
      });

      expect(onSuccess).toHaveBeenCalledWith(mockResponse);
    });
  });

  describe("リセット", () => {
    it("resetでcurrentUserの値に戻る", () => {
      const mockUser: CurrentUserResponse = {
        email: "test@example.com",
        nickname: "テストユーザー",
        weight: 70,
        height: 175,
        birthDate: "1990-01-01",
        gender: "male",
        activityLevel: "moderate",
      };

      const { result } = renderHook(() => useUpdateProfile(mockUser));

      // フィールドを変更
      act(() => {
        result.current.handleChange("nickname")("変更後");
        result.current.handleChange("weight")("80");
      });

      expect(result.current.formState.nickname).toBe("変更後");
      expect(result.current.formState.weight).toBe("80");

      // リセット
      act(() => {
        result.current.reset();
      });

      // 元の値に戻る
      expect(result.current.formState).toEqual({
        nickname: "テストユーザー",
        height: "175",
        weight: "70",
        activityLevel: "moderate",
      });
      expect(mockReset).toHaveBeenCalled();
    });
  });
});
