/**
 * useImageAnalysis フックのテスト
 */
import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, act } from "@testing-library/react";
import { useImageAnalysis } from "./useImageAnalysis";
import type { ImageFile } from "@/domain/valueObjects/imageFile";
import type { ApiErrorResponse } from "@/lib/api";

// モックの型定義
type MockTrigger = ReturnType<typeof vi.fn>;
type MockReset = ReturnType<typeof vi.fn>;
type MockOptions = {
  onSuccess?: (data: { items: { name: string; calories: number }[] }) => void;
  onError?: (error: ApiErrorResponse) => void;
};

let mockTrigger: MockTrigger;
let mockReset: MockReset;
let mockData: { items: { name: string; calories: number }[] } | undefined;
let mockIsMutating: boolean;
let mockOptions: MockOptions;

vi.mock("@/features/common/hooks", () => ({
  useRequestMutation: vi.fn(
    (
      _url: string,
      _method: string,
      options?: MockOptions
    ) => {
      mockOptions = options || {};
      return {
        trigger: mockTrigger,
        data: mockData,
        isMutating: mockIsMutating,
        reset: mockReset,
      };
    }
  ),
}));

vi.mock("@/features/common/helpers", () => ({
  getApiErrorMessage: vi.fn((code: string) => {
    const messages: Record<string, string> = {
      IMAGE_ANALYSIS_FAILED: "画像の解析に失敗しました。別の画像をお試しください",
      INVALID_IMAGE_FORMAT:
        "対応していない画像形式です。JPEG、PNG、WebPのいずれかを選択してください",
      IMAGE_TOO_LARGE:
        "画像サイズが大きすぎます。10MB以下の画像を選択してください",
      NO_FOOD_DETECTED:
        "食べ物を検出できませんでした。食べ物が写った画像を選択してください",
    };
    return messages[code] ?? "予期しないエラーが発生しました";
  }),
}));

// テスト用ImageFileの作成
const createMockImageFile = (overrides?: Partial<ImageFile>): ImageFile => ({
  base64: "dGVzdGltYWdlZGF0YQ==",
  mimeType: "image/jpeg",
  fileName: "test.jpg",
  fileSize: 1024,
  dataUrl: "data:image/jpeg;base64,dGVzdGltYWdlZGF0YQ==",
  equals: () => false,
  ...overrides,
});

describe("useImageAnalysis", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockTrigger = vi.fn().mockResolvedValue({ items: [] });
    mockReset = vi.fn();
    mockData = undefined;
    mockIsMutating = false;
    mockOptions = {};
  });

  describe("初期状態", () => {
    it("初期状態でdataはundefined", () => {
      const { result } = renderHook(() => useImageAnalysis());
      expect(result.current.data).toBeUndefined();
    });

    it("初期状態でerrorはnull", () => {
      const { result } = renderHook(() => useImageAnalysis());
      expect(result.current.error).toBeNull();
    });

    it("初期状態でisAnalyzingはfalse", () => {
      const { result } = renderHook(() => useImageAnalysis());
      expect(result.current.isAnalyzing).toBe(false);
    });

    it("analyze関数が返される", () => {
      const { result } = renderHook(() => useImageAnalysis());
      expect(typeof result.current.analyze).toBe("function");
    });

    it("reset関数が返される", () => {
      const { result } = renderHook(() => useImageAnalysis());
      expect(typeof result.current.reset).toBe("function");
    });
  });

  describe("analyze関数", () => {
    it("ImageFileを渡すとtriggerが正しいデータで呼ばれる", async () => {
      const { result } = renderHook(() => useImageAnalysis());
      const mockImageFile = createMockImageFile();

      await act(async () => {
        await result.current.analyze(mockImageFile);
      });

      expect(mockTrigger).toHaveBeenCalledWith({
        imageData: "dGVzdGltYWdlZGF0YQ==",
        mimeType: "image/jpeg",
      });
    });

    it("PNG画像でもtriggerが正しく呼ばれる", async () => {
      const { result } = renderHook(() => useImageAnalysis());
      const mockImageFile = createMockImageFile({
        base64: "cG5naW1hZ2VkYXRh",
        mimeType: "image/png",
        fileName: "test.png",
      });

      await act(async () => {
        await result.current.analyze(mockImageFile);
      });

      expect(mockTrigger).toHaveBeenCalledWith({
        imageData: "cG5naW1hZ2VkYXRh",
        mimeType: "image/png",
      });
    });

    it("WebP画像でもtriggerが正しく呼ばれる", async () => {
      const { result } = renderHook(() => useImageAnalysis());
      const mockImageFile = createMockImageFile({
        base64: "d2VicGltYWdlZGF0YQ==",
        mimeType: "image/webp",
        fileName: "test.webp",
      });

      await act(async () => {
        await result.current.analyze(mockImageFile);
      });

      expect(mockTrigger).toHaveBeenCalledWith({
        imageData: "d2VicGltYWdlZGF0YQ==",
        mimeType: "image/webp",
      });
    });

    it("analyze呼び出し時にerrorがnullにリセットされる", async () => {
      const { result } = renderHook(() => useImageAnalysis());
      const mockImageFile = createMockImageFile();

      await act(async () => {
        await result.current.analyze(mockImageFile);
      });

      expect(result.current.error).toBeNull();
    });
  });

  describe("成功時のコールバック", () => {
    it("onSuccessコールバックが呼ばれる", async () => {
      const onSuccess = vi.fn();
      const mockResponse = {
        items: [{ name: "りんご", calories: 100 }],
      };

      renderHook(() => useImageAnalysis({ onSuccess }));

      // onSuccessを直接呼び出してテスト
      act(() => {
        mockOptions.onSuccess?.(mockResponse);
      });

      expect(onSuccess).toHaveBeenCalledWith(mockResponse);
    });

    it("成功時にerrorがnullにセットされる", async () => {
      const mockResponse = {
        items: [{ name: "りんご", calories: 100 }],
      };

      const { result } = renderHook(() => useImageAnalysis());

      act(() => {
        mockOptions.onSuccess?.(mockResponse);
      });

      expect(result.current.error).toBeNull();
    });
  });

  describe("エラー時のコールバック", () => {
    it("onErrorコールバックがエラーメッセージ付きで呼ばれる", () => {
      const onError = vi.fn();

      renderHook(() => useImageAnalysis({ onError }));

      act(() => {
        mockOptions.onError?.({
          code: "IMAGE_ANALYSIS_FAILED",
          message: "Analysis failed",
        });
      });

      expect(onError).toHaveBeenCalledWith(
        "画像の解析に失敗しました。別の画像をお試しください"
      );
    });

    it("エラー時にerror状態がセットされる", () => {
      const { result } = renderHook(() => useImageAnalysis());

      act(() => {
        mockOptions.onError?.({
          code: "IMAGE_ANALYSIS_FAILED",
          message: "Analysis failed",
        });
      });

      expect(result.current.error).toBe(
        "画像の解析に失敗しました。別の画像をお試しください"
      );
    });

    it("INVALID_IMAGE_FORMATエラーが正しく処理される", () => {
      const { result } = renderHook(() => useImageAnalysis());

      act(() => {
        mockOptions.onError?.({
          code: "INVALID_IMAGE_FORMAT",
          message: "Invalid format",
        });
      });

      expect(result.current.error).toBe(
        "対応していない画像形式です。JPEG、PNG、WebPのいずれかを選択してください"
      );
    });

    it("IMAGE_TOO_LARGEエラーが正しく処理される", () => {
      const { result } = renderHook(() => useImageAnalysis());

      act(() => {
        mockOptions.onError?.({
          code: "IMAGE_TOO_LARGE",
          message: "Image too large",
        });
      });

      expect(result.current.error).toBe(
        "画像サイズが大きすぎます。10MB以下の画像を選択してください"
      );
    });

    it("NO_FOOD_DETECTEDエラーが正しく処理される", () => {
      const { result } = renderHook(() => useImageAnalysis());

      act(() => {
        mockOptions.onError?.({
          code: "NO_FOOD_DETECTED",
          message: "No food detected",
        });
      });

      expect(result.current.error).toBe(
        "食べ物を検出できませんでした。食べ物が写った画像を選択してください"
      );
    });

    it("未知のエラーコードはデフォルトメッセージを返す", () => {
      const { result } = renderHook(() => useImageAnalysis());

      act(() => {
        mockOptions.onError?.({
          code: "UNKNOWN_ERROR" as "IMAGE_ANALYSIS_FAILED",
          message: "Unknown error",
        });
      });

      expect(result.current.error).toBe("予期しないエラーが発生しました");
    });
  });

  describe("reset関数", () => {
    it("resetを呼ぶとerrorがnullになる", () => {
      const { result } = renderHook(() => useImageAnalysis());

      // エラー状態をセット
      act(() => {
        mockOptions.onError?.({
          code: "IMAGE_ANALYSIS_FAILED",
          message: "Error",
        });
      });

      expect(result.current.error).not.toBeNull();

      // リセット
      act(() => {
        result.current.reset();
      });

      expect(result.current.error).toBeNull();
    });

    it("resetを呼ぶとresetMutationが呼ばれる", () => {
      const { result } = renderHook(() => useImageAnalysis());

      act(() => {
        result.current.reset();
      });

      expect(mockReset).toHaveBeenCalled();
    });
  });

  describe("isAnalyzing状態", () => {
    it("isMutatingがtrueの場合isAnalyzingがtrue", () => {
      mockIsMutating = true;
      const { result } = renderHook(() => useImageAnalysis());
      expect(result.current.isAnalyzing).toBe(true);
    });

    it("isMutatingがfalseの場合isAnalyzingがfalse", () => {
      mockIsMutating = false;
      const { result } = renderHook(() => useImageAnalysis());
      expect(result.current.isAnalyzing).toBe(false);
    });
  });

  describe("data状態", () => {
    it("解析結果がdataとして返される", () => {
      mockData = {
        items: [
          { name: "りんご", calories: 100 },
          { name: "バナナ", calories: 80 },
        ],
      };
      const { result } = renderHook(() => useImageAnalysis());
      expect(result.current.data).toEqual({
        items: [
          { name: "りんご", calories: 100 },
          { name: "バナナ", calories: 80 },
        ],
      });
    });

    it("空の結果も正しく返される", () => {
      mockData = { items: [] };
      const { result } = renderHook(() => useImageAnalysis());
      expect(result.current.data).toEqual({ items: [] });
    });
  });
});
