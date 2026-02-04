/**
 * useImageAnalysis - 画像解析フック
 * 画像から食品情報を解析する
 */
import { useCallback, useState } from "react";
import { AxiosError } from "axios";
import { apiClient, type ApiErrorResponse } from "@/lib/api";
import { getApiErrorMessage } from "@/features/common/helpers";
import type { ImageFile } from "@/domain/valueObjects/imageFile";

/** 解析された食品アイテム */
export type AnalyzedFoodItem = {
  name: string;
  calories: number;
};

/** 画像解析レスポンス型 */
export type ImageAnalysisResponse = {
  items: AnalyzedFoodItem[];
};

/** 画像解析リクエスト型 */
type ImageAnalysisRequest = {
  imageData: string;
  mimeType: string;
};

/** フックオプション */
type UseImageAnalysisOptions = {
  onSuccess?: (data: ImageAnalysisResponse) => void;
  onError?: (errorMessage: string) => void;
};

/** フック戻り値型 */
type UseImageAnalysisReturn = {
  analyze: (imageFile: ImageFile) => Promise<void>;
  data: ImageAnalysisResponse | undefined;
  error: string | null;
  isAnalyzing: boolean;
  reset: () => void;
};

/** APIエンドポイント */
const IMAGE_ANALYSIS_ENDPOINT = "/api/v1/analyze-image";

/** 画像解析用タイムアウト: 120秒（AI処理に時間がかかるため延長） */
const IMAGE_ANALYSIS_TIMEOUT_MS = 120000;

/**
 * useImageAnalysis - 画像解析フック
 * 画像解析はAI処理のため時間がかかるので、タイムアウトを120秒に設定
 */
export function useImageAnalysis(
  options?: UseImageAnalysisOptions
): UseImageAnalysisReturn {
  const [error, setError] = useState<string | null>(null);
  const [data, setData] = useState<ImageAnalysisResponse | undefined>(undefined);
  const [isAnalyzing, setIsAnalyzing] = useState(false);

  const analyze = useCallback(
    async (imageFile: ImageFile): Promise<void> => {
      setError(null);
      setIsAnalyzing(true);

      try {
        const response = await apiClient.post<ImageAnalysisResponse>(
          IMAGE_ANALYSIS_ENDPOINT,
          {
            imageData: imageFile.base64,
            mimeType: imageFile.mimeType,
          } satisfies ImageAnalysisRequest,
          {
            timeout: IMAGE_ANALYSIS_TIMEOUT_MS,
          }
        );

        setData(response.data);
        options?.onSuccess?.(response.data);
      } catch (e) {
        let apiError: ApiErrorResponse = {
          code: "INTERNAL_ERROR",
          message: "予期しないエラーが発生しました",
        };

        if (e instanceof AxiosError && e.response?.data) {
          apiError = e.response.data as ApiErrorResponse;
        } else if (e instanceof AxiosError && e.code === "ECONNABORTED") {
          // タイムアウトエラー
          apiError = {
            code: "IMAGE_ANALYSIS_FAILED",
            message: "画像解析がタイムアウトしました。しばらく経ってから再度お試しください。",
          };
        }

        const errorMessage = getApiErrorMessage(apiError.code);
        setError(errorMessage);
        options?.onError?.(errorMessage);
      } finally {
        setIsAnalyzing(false);
      }
    },
    [options]
  );

  const reset = useCallback(() => {
    setError(null);
    setData(undefined);
  }, []);

  return {
    analyze,
    data,
    error,
    isAnalyzing,
    reset,
  };
}
