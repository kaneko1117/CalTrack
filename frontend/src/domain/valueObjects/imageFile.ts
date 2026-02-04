/**
 * 画像ファイルのValue Object
 * 画像解析APIに送信するためのファイル情報を保持する
 */

// 許可されるMIMEタイプ
const ALLOWED_MIME_TYPES = ['image/jpeg', 'image/png', 'image/webp'] as const;
export type AllowedMimeType = (typeof ALLOWED_MIME_TYPES)[number];

// 最大ファイルサイズ（10MB）
const MAX_FILE_SIZE_BYTES = 10 * 1024 * 1024;

// エラーコード
export const IMAGE_FILE_ERROR_CODE = {
  EMPTY_FILE: 'EMPTY_FILE',
  INVALID_MIME_TYPE: 'INVALID_MIME_TYPE',
  FILE_TOO_LARGE: 'FILE_TOO_LARGE',
  READ_ERROR: 'READ_ERROR',
} as const;

export type ImageFileErrorCode =
  (typeof IMAGE_FILE_ERROR_CODE)[keyof typeof IMAGE_FILE_ERROR_CODE];

// エラーメッセージ
export const IMAGE_FILE_ERROR_MESSAGE: Record<ImageFileErrorCode, string> = {
  [IMAGE_FILE_ERROR_CODE.EMPTY_FILE]: 'ファイルが選択されていません',
  [IMAGE_FILE_ERROR_CODE.INVALID_MIME_TYPE]:
    '対応していない画像形式です。JPEG、PNG、WebPのいずれかを選択してください',
  [IMAGE_FILE_ERROR_CODE.FILE_TOO_LARGE]:
    'ファイルサイズが大きすぎます。10MB以下のファイルを選択してください',
  [IMAGE_FILE_ERROR_CODE.READ_ERROR]: 'ファイルの読み込みに失敗しました',
};

// エラー型
export type ImageFileError = {
  code: ImageFileErrorCode;
  message: string;
};

// Result型
export type Result<T, E> =
  | { ok: true; value: T }
  | { ok: false; error: E };

// ImageFile型（イミュータブル）
export type ImageFile = Readonly<{
  base64: string;
  mimeType: AllowedMimeType;
  fileName: string;
  fileSize: number;
  dataUrl: string;
  equals: (other: ImageFile) => boolean;
}>;

/**
 * MIMEタイプが許可されているか検証
 */
const isAllowedMimeType = (mimeType: string): mimeType is AllowedMimeType => {
  return ALLOWED_MIME_TYPES.includes(mimeType as AllowedMimeType);
};

/**
 * FileをBase64文字列に変換
 */
const readFileAsBase64 = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => {
      const result = reader.result as string;
      // data:image/jpeg;base64,xxxxx の形式からbase64部分を抽出
      const base64 = result.split(',')[1];
      resolve(base64);
    };
    reader.onerror = () => reject(reader.error);
    reader.readAsDataURL(file);
  });
};

/**
 * FileをDataURL文字列に変換
 */
const readFileAsDataUrl = (file: File): Promise<string> => {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => resolve(reader.result as string);
    reader.onerror = () => reject(reader.error);
    reader.readAsDataURL(file);
  });
};

/**
 * ImageFileエラーを作成
 */
const createError = (code: ImageFileErrorCode): ImageFileError => ({
  code,
  message: IMAGE_FILE_ERROR_MESSAGE[code],
});

/**
 * ImageFileファクトリ関数（非同期）
 *
 * FileオブジェクトからImageFile VOを生成する。
 *
 * 注意: この関数が非同期なのは、ブラウザのFileReader APIが非同期でファイルを
 * 読み込むためです。Fileオブジェクトはバイナリデータへの参照に過ぎず、
 * 実際のバイナリ読み込みはI/O操作となるため、UIスレッドをブロックしないよう
 * 非同期処理となっています。
 *
 * @param file - 画像ファイル（File | null | undefined）
 * @returns Promise<Result<ImageFile, ImageFileError>> - 生成結果
 */
export const newImageFile = async (
  file: File | null | undefined
): Promise<Result<ImageFile, ImageFileError>> => {
  // null/undefinedチェック
  if (!file) {
    return { ok: false, error: createError(IMAGE_FILE_ERROR_CODE.EMPTY_FILE) };
  }

  // MIMEタイプチェック
  if (!isAllowedMimeType(file.type)) {
    return {
      ok: false,
      error: createError(IMAGE_FILE_ERROR_CODE.INVALID_MIME_TYPE),
    };
  }

  // ファイルサイズチェック
  if (file.size > MAX_FILE_SIZE_BYTES) {
    return {
      ok: false,
      error: createError(IMAGE_FILE_ERROR_CODE.FILE_TOO_LARGE),
    };
  }

  try {
    // Base64とDataURLを並列で読み込み
    const [base64, dataUrl] = await Promise.all([
      readFileAsBase64(file),
      readFileAsDataUrl(file),
    ]);

    const imageFile: ImageFile = {
      base64,
      mimeType: file.type as AllowedMimeType,
      fileName: file.name,
      fileSize: file.size,
      dataUrl,
      equals: (other: ImageFile) =>
        base64 === other.base64 &&
        file.type === other.mimeType &&
        file.name === other.fileName &&
        file.size === other.fileSize,
    };

    return { ok: true, value: imageFile };
  } catch {
    return { ok: false, error: createError(IMAGE_FILE_ERROR_CODE.READ_ERROR) };
  }
};

/**
 * ファイルサイズを人間が読みやすい形式にフォーマット
 * @param bytes - バイト数
 * @returns フォーマットされた文字列（例: "1.5 MB"）
 */
export const formatFileSize = (bytes: number): string => {
  if (bytes < 1024) return `${bytes} B`;
  if (bytes < 1024 * 1024) return `${(bytes / 1024).toFixed(1)} KB`;
  return `${(bytes / (1024 * 1024)).toFixed(1)} MB`;
};

/**
 * 最大ファイルサイズを取得
 * @returns 最大ファイルサイズ（バイト）
 */
export const getMaxFileSize = (): number => MAX_FILE_SIZE_BYTES;

/**
 * 許可されるMIMEタイプ一覧を取得
 * @returns 許可されるMIMEタイプの配列
 */
export const getAllowedMimeTypes = (): readonly AllowedMimeType[] =>
  ALLOWED_MIME_TYPES;
