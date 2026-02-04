import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { ImageInput, type ImageInputProps } from "./ImageInput";
import type { ImageFile } from "@/domain/valueObjects/imageFile";

// newImageFileをモック
vi.mock("@/domain/valueObjects/imageFile", async () => {
  const actual = await vi.importActual<
    typeof import("@/domain/valueObjects/imageFile")
  >("@/domain/valueObjects/imageFile");
  return {
    ...actual,
    newImageFile: vi.fn(),
  };
});

import { newImageFile } from "@/domain/valueObjects/imageFile";

const mockNewImageFile = vi.mocked(newImageFile);

// テスト用のモックImageFile
const createMockImageFile = (overrides?: Partial<ImageFile>): ImageFile => ({
  base64: "dGVzdA==",
  mimeType: "image/jpeg",
  fileName: "test.jpg",
  fileSize: 1024,
  dataUrl: "data:image/jpeg;base64,dGVzdA==",
  equals: () => true,
  ...overrides,
});

// テスト用のモックFile
const createMockFile = (
  name = "test.jpg",
  type = "image/jpeg",
  size = 1024
): File => {
  const file = new File(["test"], name, { type });
  Object.defineProperty(file, "size", { value: size });
  return file;
};

describe("ImageInput", () => {
  const defaultProps: ImageInputProps = {
    onImageSelect: vi.fn(),
  };

  beforeEach(() => {
    vi.clearAllMocks();
  });

  describe("初期表示", () => {
    it("画像選択ボタンが表示されること", () => {
      render(<ImageInput {...defaultProps} />);
      expect(
        screen.getByRole("button", { name: "画像を選択" })
      ).toBeInTheDocument();
    });

    it("許可されるファイル形式の説明が表示されること", () => {
      render(<ImageInput {...defaultProps} />);
      expect(screen.getByText(/JPEG, PNG, WebP/)).toBeInTheDocument();
      expect(screen.getByText(/10.0 MB/)).toBeInTheDocument();
    });

    it("非表示のinput要素が存在すること", () => {
      const { container } = render(<ImageInput {...defaultProps} />);
      const input = container.querySelector('input[type="file"]');
      expect(input).toBeInTheDocument();
      expect(input).toHaveClass("hidden");
    });

    it("inputのaccept属性が正しく設定されていること", () => {
      const { container } = render(<ImageInput {...defaultProps} />);
      const input = container.querySelector('input[type="file"]');
      expect(input).toHaveAttribute(
        "accept",
        "image/jpeg,image/png,image/webp"
      );
    });
  });

  describe("ボタンクリック", () => {
    it("画像選択ボタンをクリックするとinputがクリックされること", async () => {
      const { container } = render(<ImageInput {...defaultProps} />);
      const input = container.querySelector(
        'input[type="file"]'
      ) as HTMLInputElement;
      const clickSpy = vi.spyOn(input, "click");

      const button = screen.getByRole("button", { name: "画像を選択" });
      await userEvent.click(button);

      expect(clickSpy).toHaveBeenCalled();
    });
  });

  describe("ファイル選択", () => {
    it("有効なファイルを選択するとonImageSelectが呼ばれること", async () => {
      const mockImageFile = createMockImageFile();
      mockNewImageFile.mockResolvedValueOnce({
        ok: true,
        value: mockImageFile,
      });

      const onImageSelect = vi.fn();
      const { container } = render(
        <ImageInput {...defaultProps} onImageSelect={onImageSelect} />
      );

      const input = container.querySelector(
        'input[type="file"]'
      ) as HTMLInputElement;
      const file = createMockFile();

      fireEvent.change(input, { target: { files: [file] } });

      await waitFor(() => {
        expect(onImageSelect).toHaveBeenCalledWith(mockImageFile);
      });
    });

    it("ファイル選択後にinputの値がリセットされること", async () => {
      const mockImageFile = createMockImageFile();
      mockNewImageFile.mockResolvedValueOnce({
        ok: true,
        value: mockImageFile,
      });

      const { container } = render(<ImageInput {...defaultProps} />);
      const input = container.querySelector(
        'input[type="file"]'
      ) as HTMLInputElement;
      const file = createMockFile();

      fireEvent.change(input, { target: { files: [file] } });

      await waitFor(() => {
        expect(input.value).toBe("");
      });
    });

    it("ファイルが選択されなかった場合は何も起きないこと", async () => {
      const onImageSelect = vi.fn();
      const { container } = render(
        <ImageInput {...defaultProps} onImageSelect={onImageSelect} />
      );

      const input = container.querySelector(
        'input[type="file"]'
      ) as HTMLInputElement;

      fireEvent.change(input, { target: { files: [] } });

      await waitFor(() => {
        expect(onImageSelect).not.toHaveBeenCalled();
      });
    });
  });

  describe("バリデーションエラー", () => {
    it("newImageFileがエラーを返した場合、バリデーションエラーが表示されること", async () => {
      const errorMessage = "対応していない画像形式です";
      mockNewImageFile.mockResolvedValueOnce({
        ok: false,
        error: { code: "INVALID_MIME_TYPE", message: errorMessage },
      });

      const onImageSelect = vi.fn();
      const { container } = render(
        <ImageInput {...defaultProps} onImageSelect={onImageSelect} />
      );

      const input = container.querySelector(
        'input[type="file"]'
      ) as HTMLInputElement;
      const file = createMockFile("test.gif", "image/gif");

      fireEvent.change(input, { target: { files: [file] } });

      await waitFor(() => {
        expect(screen.getByText(errorMessage)).toBeInTheDocument();
      });

      expect(onImageSelect).not.toHaveBeenCalled();
    });

    it("バリデーションエラーはエラーアイコンと共に表示されること", async () => {
      const errorMessage = "ファイルサイズが大きすぎます";
      mockNewImageFile.mockResolvedValueOnce({
        ok: false,
        error: { code: "FILE_TOO_LARGE", message: errorMessage },
      });

      const { container } = render(<ImageInput {...defaultProps} />);
      const input = container.querySelector(
        'input[type="file"]'
      ) as HTMLInputElement;
      const file = createMockFile("large.jpg", "image/jpeg", 20 * 1024 * 1024);

      fireEvent.change(input, { target: { files: [file] } });

      await waitFor(() => {
        const errorElement = screen.getByText(errorMessage);
        expect(errorElement).toBeInTheDocument();
        // エラーメッセージがtext-destructiveクラスを持つ親要素内にあること
        expect(errorElement.parentElement).toHaveClass("text-destructive");
      });
    });

    it("新しいファイル選択時にバリデーションエラーがクリアされること", async () => {
      const errorMessage = "対応していない画像形式です";
      mockNewImageFile
        .mockResolvedValueOnce({
          ok: false,
          error: { code: "INVALID_MIME_TYPE", message: errorMessage },
        })
        .mockResolvedValueOnce({
          ok: true,
          value: createMockImageFile(),
        });

      const { container } = render(<ImageInput {...defaultProps} />);
      const input = container.querySelector(
        'input[type="file"]'
      ) as HTMLInputElement;

      // 最初にエラーを発生させる
      fireEvent.change(input, {
        target: { files: [createMockFile("test.gif", "image/gif")] },
      });

      await waitFor(() => {
        expect(screen.getByText(errorMessage)).toBeInTheDocument();
      });

      // 有効なファイルを選択
      fireEvent.change(input, {
        target: { files: [createMockFile()] },
      });

      await waitFor(() => {
        expect(screen.queryByText(errorMessage)).not.toBeInTheDocument();
      });
    });
  });

  describe("外部エラー表示", () => {
    it("error propが渡された場合、エラーメッセージが表示されること", () => {
      const errorMessage = "画像解析に失敗しました";
      render(<ImageInput {...defaultProps} error={errorMessage} />);
      expect(screen.getByText(errorMessage)).toBeInTheDocument();
    });

    it("error propとバリデーションエラーの両方がある場合、error propが優先されること", async () => {
      const externalError = "外部エラー";
      const validationError = "バリデーションエラー";

      mockNewImageFile.mockResolvedValueOnce({
        ok: false,
        error: { code: "INVALID_MIME_TYPE", message: validationError },
      });

      const { container } = render(
        <ImageInput {...defaultProps} error={externalError} />
      );
      const input = container.querySelector(
        'input[type="file"]'
      ) as HTMLInputElement;

      fireEvent.change(input, {
        target: { files: [createMockFile("test.gif", "image/gif")] },
      });

      await waitFor(() => {
        // validationErrorも内部状態としてセットされるが、表示されるのはexternalErrorのみ
        // (displayError = error || validationError なので)
        expect(screen.getByText(externalError)).toBeInTheDocument();
      });
    });
  });

  describe("解析中状態", () => {
    it("isAnalyzingがtrueの場合、ローディング表示になること", () => {
      render(<ImageInput {...defaultProps} isAnalyzing={true} />);

      // ボタンが非表示になる
      expect(
        screen.queryByRole("button", { name: "画像を選択" })
      ).not.toBeInTheDocument();

      // ローディングアニメーションが表示される
      const { container } = render(
        <ImageInput {...defaultProps} isAnalyzing={true} />
      );
      expect(container.querySelector(".animate-spin")).toBeInTheDocument();
    });

    it("isAnalyzingがtrueの場合、inputがdisabledになること", () => {
      const { container } = render(
        <ImageInput {...defaultProps} isAnalyzing={true} />
      );
      const input = container.querySelector('input[type="file"]');
      expect(input).toBeDisabled();
    });
  });

  describe("プレビュー表示", () => {
    it("previewUrlがある場合、プレビュー画像が表示されること", () => {
      const previewUrl = "data:image/jpeg;base64,test";
      render(<ImageInput {...defaultProps} previewUrl={previewUrl} />);

      const img = screen.getByAltText("選択された画像のプレビュー");
      expect(img).toBeInTheDocument();
      expect(img).toHaveAttribute("src", previewUrl);
    });

    it("プレビュー表示時、「別の画像を選択」ボタンが表示されること", () => {
      render(
        <ImageInput
          {...defaultProps}
          previewUrl="data:image/jpeg;base64,test"
        />
      );

      expect(
        screen.getByRole("button", { name: "別の画像を選択" })
      ).toBeInTheDocument();
      expect(
        screen.queryByRole("button", { name: "画像を選択" })
      ).not.toBeInTheDocument();
    });

    it("プレビュー表示時、ファイル形式の説明が非表示になること", () => {
      render(
        <ImageInput
          {...defaultProps}
          previewUrl="data:image/jpeg;base64,test"
        />
      );

      expect(screen.queryByText(/JPEG, PNG, WebP/)).not.toBeInTheDocument();
    });

    it("プレビュー表示時、クリアボタンが表示されること", () => {
      render(
        <ImageInput
          {...defaultProps}
          previewUrl="data:image/jpeg;base64,test"
          onClear={vi.fn()}
        />
      );

      expect(
        screen.getByRole("button", { name: "画像をクリア" })
      ).toBeInTheDocument();
    });

    it("onClearがない場合、クリアボタンが表示されないこと", () => {
      render(
        <ImageInput
          {...defaultProps}
          previewUrl="data:image/jpeg;base64,test"
        />
      );

      expect(
        screen.queryByRole("button", { name: "画像をクリア" })
      ).not.toBeInTheDocument();
    });
  });

  describe("クリア機能", () => {
    it("クリアボタンをクリックするとonClearが呼ばれること", async () => {
      const onClear = vi.fn();
      render(
        <ImageInput
          {...defaultProps}
          previewUrl="data:image/jpeg;base64,test"
          onClear={onClear}
        />
      );

      const clearButton = screen.getByRole("button", { name: "画像をクリア" });
      await userEvent.click(clearButton);

      expect(onClear).toHaveBeenCalled();
    });

    it("disabledの場合、クリアボタンが表示されないこと", () => {
      render(
        <ImageInput
          {...defaultProps}
          previewUrl="data:image/jpeg;base64,test"
          onClear={vi.fn()}
          disabled={true}
        />
      );

      expect(
        screen.queryByRole("button", { name: "画像をクリア" })
      ).not.toBeInTheDocument();
    });

    it("クリア時にバリデーションエラーもクリアされること", async () => {
      const errorMessage = "対応していない画像形式です";
      mockNewImageFile.mockResolvedValueOnce({
        ok: false,
        error: { code: "INVALID_MIME_TYPE", message: errorMessage },
      });

      const onClear = vi.fn();
      const { container, rerender } = render(
        <ImageInput {...defaultProps} onClear={onClear} />
      );

      const input = container.querySelector(
        'input[type="file"]'
      ) as HTMLInputElement;
      fireEvent.change(input, {
        target: { files: [createMockFile("test.gif", "image/gif")] },
      });

      await waitFor(() => {
        expect(screen.getByText(errorMessage)).toBeInTheDocument();
      });

      // プレビューを追加してクリアボタンを表示
      rerender(
        <ImageInput
          {...defaultProps}
          previewUrl="data:image/jpeg;base64,test"
          onClear={onClear}
        />
      );

      const clearButton = screen.getByRole("button", { name: "画像をクリア" });
      await userEvent.click(clearButton);

      expect(screen.queryByText(errorMessage)).not.toBeInTheDocument();
    });
  });

  describe("disabled状態", () => {
    it("disabledがtrueの場合、ボタンがdisabledになること", () => {
      render(<ImageInput {...defaultProps} disabled={true} />);

      const button = screen.getByRole("button", { name: "画像を選択" });
      expect(button).toBeDisabled();
    });

    it("disabledがtrueの場合、inputもdisabledになること", () => {
      const { container } = render(
        <ImageInput {...defaultProps} disabled={true} />
      );
      const input = container.querySelector('input[type="file"]');
      expect(input).toBeDisabled();
    });

    it("プレビュー表示中でdisabledの場合、「別の画像を選択」ボタンもdisabledになること", () => {
      render(
        <ImageInput
          {...defaultProps}
          previewUrl="data:image/jpeg;base64,test"
          disabled={true}
        />
      );

      const button = screen.getByRole("button", { name: "別の画像を選択" });
      expect(button).toBeDisabled();
    });
  });

  describe("解析中とプレビューの組み合わせ", () => {
    it("isAnalyzingがtrueでpreviewUrlがある場合、ローディング表示が優先されること", () => {
      render(
        <ImageInput
          {...defaultProps}
          isAnalyzing={true}
          previewUrl="data:image/jpeg;base64,test"
        />
      );

      // プレビュー画像は非表示
      expect(
        screen.queryByAltText("選択された画像のプレビュー")
      ).not.toBeInTheDocument();

      // ローディングアニメーションが表示される
      const { container } = render(
        <ImageInput
          {...defaultProps}
          isAnalyzing={true}
          previewUrl="data:image/jpeg;base64,test"
        />
      );
      expect(container.querySelector(".animate-spin")).toBeInTheDocument();
    });
  });
});
