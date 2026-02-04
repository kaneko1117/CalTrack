/**
 * useImageAnalysis フックのテスト
 */
import { describe, it, expect, vi, beforeEach } from "vitest";
import { renderHook, act, waitFor } from "@testing-library/react";
import { useImageAnalysis } from "./useImageAnalysis";
import type { ImageFile } from "@/domain/valueObjects/imageFile";
import { apiClient } from "@/lib/api";
import { AxiosError, type AxiosResponse } from "axios";

// apiClientをモック
vi.mock("@/lib/api", () => ({
  apiClient: {
    post: vi.fn(),
  },
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

const mockApiClient = apiClient as unknown as {
  post: ReturnType<typeof vi.fn>;
};

describe("useImageAnalysis", () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockApiClient.post.mockResolvedValue({ data: { items: [] } });
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
    it("ImageFileを渡すとapiClient.postが正しいデータで呼ばれる", async () => {
      const { result } = renderHook(() => useImageAnalysis());
      const mockImageFile = createMockImageFile();

      await act(async () => {
        await result.current.analyze(mockImageFile);
      });

      expect(mockApiClient.post).toHaveBeenCalledWith(
        "/api/v1/analyze-image",
        {
          imageData: "dGVzdGltYWdlZGF0YQ==",
          mimeType: "image/jpeg",
        },
        {
          timeout: 120000,
        }
      );
    });

    it("PNG画像でもapiClient.postが正しく呼ばれる", async () => {
      const { result } = renderHook(() => useImageAnalysis());
      const mockImageFile = createMockImageFile({
        base64: "cG5naW1hZ2VkYXRh",
        mimeType: "image/png",
        fileName: "test.png",
      });

      await act(async () => {
        await result.current.analyze(mockImageFile);
      });

      expect(mockApiClient.post).toHaveBeenCalledWith(
        "/api/v1/analyze-image",
        {
          imageData: "cG5naW1hZ2VkYXRh",
          mimeType: "image/png",
        },
        {
          timeout: 120000,
        }
      );
    });

    it("WebP画像でもapiClient.postが正しく呼ばれる", async () => {
      const { result } = renderHook(() => useImageAnalysis());
      const mockImageFile = createMockImageFile({
        base64: "d2VicGltYWdlZGF0YQ==",
        mimeType: "image/webp",
        fileName: "test.webp",
      });

      await act(async () => {
        await result.current.analyze(mockImageFile);
      });

      expect(mockApiClient.post).toHaveBeenCalledWith(
        "/api/v1/analyze-image",
        {
          imageData: "d2VicGltYWdlZGF0YQ==",
          mimeType: "image/webp",
        },
        {
          timeout: 120000,
        }
      );
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
      mockApiClient.post.mockResolvedValue({ data: mockResponse });

      const { result } = renderHook(() => useImageAnalysis({ onSuccess }));
      const mockImageFile = createMockImageFile();

      await act(async () => {
        await result.current.analyze(mockImageFile);
      });

      expect(onSuccess).toHaveBeenCalledWith(mockResponse);
    });

    it("成功時にdataがセットされる", async () => {
      const mockResponse = {
        items: [{ name: "りんご", calories: 100 }],
      };
      mockApiClient.post.mockResolvedValue({ data: mockResponse });

      const { result } = renderHook(() => useImageAnalysis());
      const mockImageFile = createMockImageFile();

      await act(async () => {
        await result.current.analyze(mockImageFile);
      });

      expect(result.current.data).toEqual(mockResponse);
    });

    it("成功時にerrorがnullにセットされる", async () => {
      const mockResponse = {
        items: [{ name: "りんご", calories: 100 }],
      };
      mockApiClient.post.mockResolvedValue({ data: mockResponse });

      const { result } = renderHook(() => useImageAnalysis());
      const mockImageFile = createMockImageFile();

      await act(async () => {
        await result.current.analyze(mockImageFile);
      });

      expect(result.current.error).toBeNull();
    });
  });

  describe("エラー時のコールバック", () => {
    const createAxiosError = (code: string, message: string): AxiosError => {
      const error = new AxiosError(message);
      error.response = {
        data: { code, message },
        status: 400,
        statusText: "Bad Request",
        headers: {},
        config: {} as AxiosResponse["config"],
      };
      return error;
    };

    it("onErrorコールバックがエラーメッセージ付きで呼ばれる", async () => {
      const onError = vi.fn();
      mockApiClient.post.mockRejectedValue(
        createAxiosError("IMAGE_ANALYSIS_FAILED", "Analysis failed")
      );

      const { result } = renderHook(() => useImageAnalysis({ onError }));
      const mockImageFile = createMockImageFile();

      await act(async () => {
        await result.current.analyze(mockImageFile);
      });

      expect(onError).toHaveBeenCalledWith(
        "画像の解析に失敗しました。別の画像をお試しください"
      );
    });

    it("エラー時にerror状態がセットされる", async () => {
      mockApiClient.post.mockRejectedValue(
        createAxiosError("IMAGE_ANALYSIS_FAILED", "Analysis failed")
      );

      const { result } = renderHook(() => useImageAnalysis());
      const mockImageFile = createMockImageFile();

      await act(async () => {
        await result.current.analyze(mockImageFile);
      });

      expect(result.current.error).toBe(
        "画像の解析に失敗しました。別の画像をお試しください"
      );
    });

    it("INVALID_IMAGE_FORMATエラーが正しく処理される", async () => {
      mockApiClient.post.mockRejectedValue(
        createAxiosError("INVALID_IMAGE_FORMAT", "Invalid format")
      );

      const { result } = renderHook(() => useImageAnalysis());
      const mockImageFile = createMockImageFile();

      await act(async () => {
        await result.current.analyze(mockImageFile);
      });

      expect(result.current.error).toBe(
        "対応していない画像形式です。JPEG、PNG、WebPのいずれかを選択してください"
      );
    });

    it("IMAGE_TOO_LARGEエラーが正しく処理される", async () => {
      mockApiClient.post.mockRejectedValue(
        createAxiosError("IMAGE_TOO_LARGE", "Image too large")
      );

      const { result } = renderHook(() => useImageAnalysis());
      const mockImageFile = createMockImageFile();

      await act(async () => {
        await result.current.analyze(mockImageFile);
      });

      expect(result.current.error).toBe(
        "画像サイズが大きすぎます。10MB以下の画像を選択してください"
      );
    });

    it("NO_FOOD_DETECTEDエラーが正しく処理される", async () => {
      mockApiClient.post.mockRejectedValue(
        createAxiosError("NO_FOOD_DETECTED", "No food detected")
      );

      const { result } = renderHook(() => useImageAnalysis());
      const mockImageFile = createMockImageFile();

      await act(async () => {
        await result.current.analyze(mockImageFile);
      });

      expect(result.current.error).toBe(
        "食べ物を検出できませんでした。食べ物が写った画像を選択してください"
      );
    });

    it("未知のエラーコードはデフォルトメッセージを返す", async () => {
      mockApiClient.post.mockRejectedValue(
        createAxiosError("UNKNOWN_ERROR", "Unknown error")
      );

      const { result } = renderHook(() => useImageAnalysis());
      const mockImageFile = createMockImageFile();

      await act(async () => {
        await result.current.analyze(mockImageFile);
      });

      expect(result.current.error).toBe("予期しないエラーが発生しました");
    });

    it("タイムアウトエラーが正しく処理される", async () => {
      const timeoutError = new AxiosError("timeout");
      timeoutError.code = "ECONNABORTED";
      mockApiClient.post.mockRejectedValue(timeoutError);

      const { result } = renderHook(() => useImageAnalysis());
      const mockImageFile = createMockImageFile();

      await act(async () => {
        await result.current.analyze(mockImageFile);
      });

      expect(result.current.error).toBe(
        "画像の解析に失敗しました。別の画像をお試しください"
      );
    });
  });

  describe("reset関数", () => {
    it("resetを呼ぶとerrorがnullになる", async () => {
      mockApiClient.post.mockRejectedValue(
        new AxiosError("Error", undefined, undefined, undefined, {
          data: { code: "IMAGE_ANALYSIS_FAILED", message: "Error" },
          status: 400,
          statusText: "Bad Request",
          headers: {},
          config: {} as AxiosResponse["config"],
        })
      );

      const { result } = renderHook(() => useImageAnalysis());
      const mockImageFile = createMockImageFile();

      // エラー状態をセット
      await act(async () => {
        await result.current.analyze(mockImageFile);
      });

      expect(result.current.error).not.toBeNull();

      // リセット
      act(() => {
        result.current.reset();
      });

      expect(result.current.error).toBeNull();
    });

    it("resetを呼ぶとdataがundefinedになる", async () => {
      const mockResponse = {
        items: [{ name: "りんご", calories: 100 }],
      };
      mockApiClient.post.mockResolvedValue({ data: mockResponse });

      const { result } = renderHook(() => useImageAnalysis());
      const mockImageFile = createMockImageFile();

      // データをセット
      await act(async () => {
        await result.current.analyze(mockImageFile);
      });

      expect(result.current.data).not.toBeUndefined();

      // リセット
      act(() => {
        result.current.reset();
      });

      expect(result.current.data).toBeUndefined();
    });
  });

  describe("isAnalyzing状態", () => {
    it("analyze実行中はisAnalyzingがtrue", async () => {
      let resolvePromise: (value: { data: { items: [] } }) => void;
      const pendingPromise = new Promise<{ data: { items: [] } }>((resolve) => {
        resolvePromise = resolve;
      });
      mockApiClient.post.mockReturnValue(pendingPromise);

      const { result } = renderHook(() => useImageAnalysis());
      const mockImageFile = createMockImageFile();

      // 非同期で実行開始
      act(() => {
        result.current.analyze(mockImageFile);
      });

      // 実行中はtrue
      await waitFor(() => {
        expect(result.current.isAnalyzing).toBe(true);
      });

      // Promiseを解決
      await act(async () => {
        resolvePromise!({ data: { items: [] } });
      });

      // 完了後はfalse
      await waitFor(() => {
        expect(result.current.isAnalyzing).toBe(false);
      });
    });

    it("analyze完了後はisAnalyzingがfalse", async () => {
      mockApiClient.post.mockResolvedValue({ data: { items: [] } });

      const { result } = renderHook(() => useImageAnalysis());
      const mockImageFile = createMockImageFile();

      await act(async () => {
        await result.current.analyze(mockImageFile);
      });

      expect(result.current.isAnalyzing).toBe(false);
    });
  });

  describe("data状態", () => {
    it("解析結果がdataとして返される", async () => {
      const mockResponse = {
        items: [
          { name: "りんご", calories: 100 },
          { name: "バナナ", calories: 80 },
        ],
      };
      mockApiClient.post.mockResolvedValue({ data: mockResponse });

      const { result } = renderHook(() => useImageAnalysis());
      const mockImageFile = createMockImageFile();

      await act(async () => {
        await result.current.analyze(mockImageFile);
      });

      expect(result.current.data).toEqual({
        items: [
          { name: "りんご", calories: 100 },
          { name: "バナナ", calories: 80 },
        ],
      });
    });

    it("空の結果も正しく返される", async () => {
      mockApiClient.post.mockResolvedValue({ data: { items: [] } });

      const { result } = renderHook(() => useImageAnalysis());
      const mockImageFile = createMockImageFile();

      await act(async () => {
        await result.current.analyze(mockImageFile);
      });

      expect(result.current.data).toEqual({ items: [] });
    });
  });
});
