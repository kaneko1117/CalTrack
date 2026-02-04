/**
 * useImageAnalysis - 画像解析フック
 * 画像から食品情報を解析する
 */
import { useCallback, useState } from "react";
import { useRequestMutation } from "@/features/common/hooks";
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

/**
 * useImageAnalysis - 画像解析フック
 */
export function useImageAnalysis(
  options?: UseImageAnalysisOptions
): UseImageAnalysisReturn {
  const [error, setError] = useState<string | null>(null);

  const {
    trigger,
    data,
    isMutating,
    reset: resetMutation,
  } = useRequestMutation<ImageAnalysisResponse, ImageAnalysisRequest>(
    IMAGE_ANALYSIS_ENDPOINT,
    "POST",
    {
      onSuccess: (responseData) => {
        setError(null);
        options?.onSuccess?.(responseData);
      },
      onError: (apiError) => {
        const errorMessage = getApiErrorMessage(apiError.code);
        setError(errorMessage);
        options?.onError?.(errorMessage);
      },
    }
  );

  const analyze = useCallback(
    async (imageFile: ImageFile): Promise<void> => {
      setError(null);
      await trigger({
        imageData: imageFile.base64,
        mimeType: imageFile.mimeType,
      });
    },
    [trigger]
  );

  const reset = useCallback(() => {
    setError(null);
    resetMutation();
  }, [resetMutation]);

  return {
    analyze,
    data,
    error,
    isAnalyzing: isMutating,
    reset,
  };
}
