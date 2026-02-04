/**
 * ImageInput - シンプルな画像選択コンポーネント
 */
import { useCallback, useRef, useState } from "react";
import { Button } from "@/components/ui/button";
import { Skeleton } from "@/components/ui/skeleton";
import {
  newImageFile,
  formatFileSize,
  getAllowedMimeTypes,
  type ImageFile,
} from "@/domain/valueObjects/imageFile";

export type ImageInputProps = {
  onImageSelect: (imageFile: ImageFile) => void;
  isAnalyzing?: boolean;
  error?: string | null;
  previewUrl?: string | null;
  disabled?: boolean;
  onClear?: () => void;
};

function CameraIcon({ className }: { className?: string }) {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={className} aria-hidden="true">
      <path d="M14.5 4h-5L7 7H4a2 2 0 0 0-2 2v9a2 2 0 0 0 2 2h16a2 2 0 0 0 2-2V9a2 2 0 0 0-2-2h-3l-2.5-3z" />
      <circle cx="12" cy="13" r="3" />
    </svg>
  );
}

function LoaderIcon({ className }: { className?: string }) {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={className} aria-hidden="true">
      <path d="M21 12a9 9 0 1 1-6.219-8.56" />
    </svg>
  );
}

function XIcon({ className }: { className?: string }) {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={className} aria-hidden="true">
      <line x1="18" y1="6" x2="6" y2="18" />
      <line x1="6" y1="6" x2="18" y2="18" />
    </svg>
  );
}

function AlertCircleIcon({ className }: { className?: string }) {
  return (
    <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round" className={className} aria-hidden="true">
      <circle cx="12" cy="12" r="10" />
      <line x1="12" y1="8" x2="12" y2="12" />
      <line x1="12" y1="16" x2="12.01" y2="16" />
    </svg>
  );
}

export function ImageInput({
  onImageSelect,
  isAnalyzing = false,
  error,
  previewUrl,
  disabled = false,
  onClear,
}: ImageInputProps) {
  const [validationError, setValidationError] = useState<string | null>(null);
  const inputRef = useRef<HTMLInputElement>(null);

  const allowedMimeTypes = getAllowedMimeTypes();
  const acceptAttribute = allowedMimeTypes.join(",");

  const handleChange = useCallback(
    async (e: React.ChangeEvent<HTMLInputElement>) => {
      const files = e.target.files;
      if (files && files.length > 0) {
        setValidationError(null);
        const result = await newImageFile(files[0]);
        if (!result.ok) {
          setValidationError(result.error.message);
          return;
        }
        onImageSelect(result.value);
      }
      e.target.value = "";
    },
    [onImageSelect]
  );

  const handleClear = useCallback(() => {
    setValidationError(null);
    onClear?.();
  }, [onClear]);

  const handleButtonClick = useCallback(() => {
    inputRef.current?.click();
  }, []);

  const displayError = error || validationError;
  const isDisabled = disabled || isAnalyzing;

  return (
    <div className="space-y-3">
      <input
        ref={inputRef}
        type="file"
        accept={acceptAttribute}
        onChange={handleChange}
        className="hidden"
        disabled={isDisabled}
      />

      {isAnalyzing && (
        <div className="flex items-center gap-3 p-4 border rounded-lg bg-muted/30">
          <LoaderIcon className="w-5 h-5 text-primary animate-spin" />
          <div className="flex-1 space-y-2">
            <Skeleton className="h-4 w-3/4" />
            <Skeleton className="h-4 w-1/2" />
          </div>
        </div>
      )}

      {!isAnalyzing && previewUrl && (
        <div className="relative">
          <img
            src={previewUrl}
            alt="選択された画像のプレビュー"
            className="w-full h-32 object-contain rounded-lg border bg-muted/30"
          />
          {onClear && !disabled && (
            <button
              type="button"
              onClick={handleClear}
              className="absolute top-2 right-2 p-1.5 bg-background/90 rounded-full hover:bg-background transition-colors border shadow-sm"
              aria-label="画像をクリア"
            >
              <XIcon className="w-4 h-4" />
            </button>
          )}
        </div>
      )}

      {!isAnalyzing && !previewUrl && (
        <Button
          type="button"
          variant="outline"
          onClick={handleButtonClick}
          disabled={isDisabled}
          className="w-full h-12"
        >
          <CameraIcon className="w-5 h-5 mr-2" />
          画像を選択
        </Button>
      )}

      {!isAnalyzing && previewUrl && (
        <Button
          type="button"
          variant="outline"
          onClick={handleButtonClick}
          disabled={isDisabled}
          className="w-full"
          size="sm"
        >
          別の画像を選択
        </Button>
      )}

      {!previewUrl && !isAnalyzing && (
        <p className="text-xs text-muted-foreground text-center">
          JPEG, PNG, WebP ({formatFileSize(10 * 1024 * 1024)}まで)
        </p>
      )}

      {displayError && (
        <p className="flex items-center gap-1.5 text-sm text-destructive">
          <AlertCircleIcon className="w-4 h-4 flex-shrink-0" />
          <span>{displayError}</span>
        </p>
      )}
    </div>
  );
}
