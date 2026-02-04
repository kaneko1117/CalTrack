import { describe, it, expect, vi, beforeEach } from 'vitest';
import {
  newImageFile,
  formatFileSize,
  getMaxFileSize,
  getAllowedMimeTypes,
  IMAGE_FILE_ERROR_CODE,
  IMAGE_FILE_ERROR_MESSAGE,
} from './imageFile';

// FileReaderモック用の設定
let mockResult: string | null = null;
let mockShouldError = false;

// FileReaderのモッククラス
class MockFileReader {
  result: string | null = null;
  error: Error | null = null;
  onload: (() => void) | null = null;
  onerror: (() => void) | null = null;

  readAsDataURL() {
    // 非同期でコールバックを発火
    setTimeout(() => {
      if (mockShouldError) {
        this.error = new Error('Read failed');
        this.onerror?.();
      } else {
        this.result = mockResult;
        this.onload?.();
      }
    }, 0);
  }
}

vi.stubGlobal('FileReader', MockFileReader);

// テスト用ヘルパー
const createMockFile = (
  name: string,
  size: number,
  type: string
): File => {
  const file = new File([''], name, { type });
  Object.defineProperty(file, 'size', { value: size });
  return file;
};

describe('ImageFile VO', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockResult = null;
    mockShouldError = false;
  });

  describe('newImageFile', () => {
    describe('バリデーション', () => {
      it('nullの場合はEMPTY_FILEエラーを返す', async () => {
        const result = await newImageFile(null);

        expect(result.ok).toBe(false);
        if (!result.ok) {
          expect(result.error.code).toBe(IMAGE_FILE_ERROR_CODE.EMPTY_FILE);
          expect(result.error.message).toBe(
            IMAGE_FILE_ERROR_MESSAGE.EMPTY_FILE
          );
        }
      });

      it('undefinedの場合はEMPTY_FILEエラーを返す', async () => {
        const result = await newImageFile(undefined);

        expect(result.ok).toBe(false);
        if (!result.ok) {
          expect(result.error.code).toBe(IMAGE_FILE_ERROR_CODE.EMPTY_FILE);
        }
      });

      it('許可されていないMIMEタイプの場合はINVALID_MIME_TYPEエラーを返す', async () => {
        const file = createMockFile('test.gif', 1024, 'image/gif');
        const result = await newImageFile(file);

        expect(result.ok).toBe(false);
        if (!result.ok) {
          expect(result.error.code).toBe(
            IMAGE_FILE_ERROR_CODE.INVALID_MIME_TYPE
          );
          expect(result.error.message).toBe(
            IMAGE_FILE_ERROR_MESSAGE.INVALID_MIME_TYPE
          );
        }
      });

      it('ファイルサイズが10MBを超える場合はFILE_TOO_LARGEエラーを返す', async () => {
        const file = createMockFile(
          'large.jpg',
          11 * 1024 * 1024,
          'image/jpeg'
        );
        const result = await newImageFile(file);

        expect(result.ok).toBe(false);
        if (!result.ok) {
          expect(result.error.code).toBe(IMAGE_FILE_ERROR_CODE.FILE_TOO_LARGE);
          expect(result.error.message).toBe(
            IMAGE_FILE_ERROR_MESSAGE.FILE_TOO_LARGE
          );
        }
      });
    });

    describe('正常系', () => {
      it('有効なJPEGファイルからImageFileを生成できる', async () => {
        const file = createMockFile('photo.jpg', 1024, 'image/jpeg');
        mockResult = 'data:image/jpeg;base64,dGVzdA==';

        const result = await newImageFile(file);

        expect(result.ok).toBe(true);
        if (result.ok) {
          expect(result.value.mimeType).toBe('image/jpeg');
          expect(result.value.fileName).toBe('photo.jpg');
          expect(result.value.fileSize).toBe(1024);
          expect(result.value.base64).toBe('dGVzdA==');
          expect(result.value.dataUrl).toBe('data:image/jpeg;base64,dGVzdA==');
        }
      });

      it('有効なPNGファイルからImageFileを生成できる', async () => {
        const file = createMockFile('image.png', 2048, 'image/png');
        mockResult = 'data:image/png;base64,cG5n';

        const result = await newImageFile(file);

        expect(result.ok).toBe(true);
        if (result.ok) {
          expect(result.value.mimeType).toBe('image/png');
        }
      });

      it('有効なWebPファイルからImageFileを生成できる', async () => {
        const file = createMockFile('image.webp', 512, 'image/webp');
        mockResult = 'data:image/webp;base64,d2VicA==';

        const result = await newImageFile(file);

        expect(result.ok).toBe(true);
        if (result.ok) {
          expect(result.value.mimeType).toBe('image/webp');
        }
      });
    });

    describe('equals', () => {
      it('同じ内容のImageFile同士はequalsがtrueを返す', async () => {
        const file = createMockFile('test.jpg', 1024, 'image/jpeg');
        mockResult = 'data:image/jpeg;base64,dGVzdA==';

        const result1 = await newImageFile(file);
        const result2 = await newImageFile(file);

        if (result1.ok && result2.ok) {
          expect(result1.value.equals(result2.value)).toBe(true);
        }
      });
    });

    describe('エラーハンドリング', () => {
      it('FileReader読み込みエラー時はREAD_ERRORを返す', async () => {
        const file = createMockFile('error.jpg', 1024, 'image/jpeg');
        mockShouldError = true;

        const result = await newImageFile(file);

        expect(result.ok).toBe(false);
        if (!result.ok) {
          expect(result.error.code).toBe(IMAGE_FILE_ERROR_CODE.READ_ERROR);
        }
      });
    });
  });

  describe('formatFileSize', () => {
    it('1024バイト未満はB単位で表示する', () => {
      expect(formatFileSize(500)).toBe('500 B');
      expect(formatFileSize(0)).toBe('0 B');
    });

    it('1KB以上1MB未満はKB単位で表示する', () => {
      expect(formatFileSize(1024)).toBe('1.0 KB');
      expect(formatFileSize(1536)).toBe('1.5 KB');
      expect(formatFileSize(102400)).toBe('100.0 KB');
    });

    it('1MB以上はMB単位で表示する', () => {
      expect(formatFileSize(1024 * 1024)).toBe('1.0 MB');
      expect(formatFileSize(5 * 1024 * 1024)).toBe('5.0 MB');
      expect(formatFileSize(10.5 * 1024 * 1024)).toBe('10.5 MB');
    });
  });

  describe('getMaxFileSize', () => {
    it('最大ファイルサイズ10MBを返す', () => {
      expect(getMaxFileSize()).toBe(10 * 1024 * 1024);
    });
  });

  describe('getAllowedMimeTypes', () => {
    it('許可されるMIMEタイプ一覧を返す', () => {
      const types = getAllowedMimeTypes();
      expect(types).toContain('image/jpeg');
      expect(types).toContain('image/png');
      expect(types).toContain('image/webp');
      expect(types.length).toBe(3);
    });
  });
});
