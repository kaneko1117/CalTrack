/**
 * RecordForm コンポーネントのテスト
 */
import { describe, it, expect, vi, beforeEach } from "vitest";
import { render, screen, fireEvent, waitFor } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { RecordForm } from "./RecordForm";

// APIをモック
vi.mock("@/lib/api", () => ({
  post: vi.fn(),
}));

import { post } from "@/lib/api";

describe("RecordForm", () => {
  const mockOnSuccess = vi.fn();

  beforeEach(() => {
    vi.clearAllMocks();
    vi.mocked(post).mockResolvedValue({
      recordId: "test-id",
      eatenAt: "2024-01-01T12:00:00Z",
      totalCalories: 250,
    });
  });

  describe("初期状態", () => {
    it("初期状態で食事日時が設定されている", () => {
      render(<RecordForm />);

      const dateTimeInput = screen.getByLabelText("食事日時");
      expect(dateTimeInput).toBeInTheDocument();
      // datetime-localの値が設定されていることを確認
      expect(dateTimeInput).toHaveValue();
    });

    it("初期状態で1つの空の食品アイテムがある", () => {
      render(<RecordForm />);

      // 「食品 1」のラベルが表示されている
      expect(screen.getByText("食品 1")).toBeInTheDocument();

      // 食品名とカロリーの入力フィールドが1組ある
      const foodNameInputs = screen.getAllByLabelText("食品名");
      expect(foodNameInputs).toHaveLength(1);
      expect(foodNameInputs[0]).toHaveValue("");

      const calorieInputs = screen.getAllByLabelText("カロリー (kcal)");
      expect(calorieInputs).toHaveLength(1);
      expect(calorieInputs[0]).toHaveValue(null); // 0はnullとして表示
    });
  });

  describe("入力", () => {
    it("食品名を入力できる", async () => {
      const user = userEvent.setup();
      render(<RecordForm />);

      const foodNameInput = screen.getByLabelText("食品名");
      await user.type(foodNameInput, "白米");

      expect(foodNameInput).toHaveValue("白米");
    });

    it("カロリーを入力できる", async () => {
      const user = userEvent.setup();
      render(<RecordForm />);

      const calorieInput = screen.getByLabelText("カロリー (kcal)");
      await user.type(calorieInput, "250");

      expect(calorieInput).toHaveValue(250);
    });
  });

  describe("食品アイテムの追加と削除", () => {
    it("「食品を追加」ボタンでアイテムを追加できる", async () => {
      const user = userEvent.setup();
      render(<RecordForm />);

      // 初期は1つ
      expect(screen.getAllByLabelText("食品名")).toHaveLength(1);

      // 追加ボタンをクリック
      const addButton = screen.getByRole("button", { name: "食品を追加" });
      await user.click(addButton);

      // 2つになる
      expect(screen.getAllByLabelText("食品名")).toHaveLength(2);
      expect(screen.getByText("食品 1")).toBeInTheDocument();
      expect(screen.getByText("食品 2")).toBeInTheDocument();
    });

    it("削除ボタンでアイテムを削除できる", async () => {
      const user = userEvent.setup();
      render(<RecordForm />);

      // まず2つにする
      const addButton = screen.getByRole("button", { name: "食品を追加" });
      await user.click(addButton);
      expect(screen.getAllByLabelText("食品名")).toHaveLength(2);

      // 2番目を削除
      const deleteButton = screen.getByRole("button", { name: "食品 2 を削除" });
      await user.click(deleteButton);

      // 1つになる
      expect(screen.getAllByLabelText("食品名")).toHaveLength(1);
    });

    it("アイテムが1つの場合は削除ボタンが表示されない", () => {
      render(<RecordForm />);

      // 削除ボタンがない
      const deleteButtons = screen.queryAllByRole("button", { name: /を削除$/ });
      expect(deleteButtons).toHaveLength(0);
    });
  });

  describe("合計カロリー表示", () => {
    it("合計カロリーが表示される", () => {
      render(<RecordForm />);

      expect(screen.getByText(/合計:.*0 kcal/)).toBeInTheDocument();
    });

    it("カロリーを入力すると合計が更新される", async () => {
      const user = userEvent.setup();
      render(<RecordForm />);

      const calorieInput = screen.getByLabelText("カロリー (kcal)");
      await user.type(calorieInput, "250");

      expect(screen.getByText(/合計:.*250 kcal/)).toBeInTheDocument();
    });

    it("複数アイテムの合計が計算される", async () => {
      const user = userEvent.setup();
      render(<RecordForm />);

      // 1つ目に250kcal入力
      const firstCalorieInput = screen.getByLabelText("カロリー (kcal)");
      await user.type(firstCalorieInput, "250");

      // 2つ目を追加して150kcal入力
      await user.click(screen.getByRole("button", { name: "食品を追加" }));
      const calorieInputs = screen.getAllByLabelText("カロリー (kcal)");
      await user.type(calorieInputs[1], "150");

      // 合計400kcal
      expect(screen.getByText(/合計:.*400 kcal/)).toBeInTheDocument();
    });
  });

  describe("バリデーション", () => {
    it("バリデーションエラーが表示される（食品名が空）", async () => {
      const user = userEvent.setup();
      render(<RecordForm />);

      // 食品名を空のままカロリーを入力
      const calorieInput = screen.getByLabelText("カロリー (kcal)");
      await user.type(calorieInput, "250");

      // 食品名フィールドからフォーカスを外す（blurトリガー）
      const foodNameInput = screen.getByLabelText("食品名");
      await user.click(foodNameInput);
      await user.tab();

      await waitFor(() => {
        expect(screen.getByText("食品名を入力してください")).toBeInTheDocument();
      });
    });

    it("バリデーションエラーが表示される（カロリーが0以下）", async () => {
      const user = userEvent.setup();
      render(<RecordForm />);

      // 食品名を入力
      const foodNameInput = screen.getByLabelText("食品名");
      await user.type(foodNameInput, "白米");

      // カロリーに0を入力
      const calorieInput = screen.getByLabelText("カロリー (kcal)");
      await user.type(calorieInput, "0");
      await user.tab();

      await waitFor(() => {
        expect(screen.getByText("カロリーは1以上の整数で入力してください")).toBeInTheDocument();
      });
    });
  });

  describe("フォーム送信", () => {
    it("フォーム送信時にAPIが呼ばれる", async () => {
      const user = userEvent.setup();
      render(<RecordForm onSuccess={mockOnSuccess} />);

      // 食品名とカロリーを入力
      const foodNameInput = screen.getByLabelText("食品名");
      await user.type(foodNameInput, "白米");

      const calorieInput = screen.getByLabelText("カロリー (kcal)");
      await user.type(calorieInput, "250");

      // ボタンが有効化されることを確認
      await waitFor(() => {
        expect(screen.getByRole("button", { name: "記録する" })).not.toBeDisabled();
      });

      // 送信
      fireEvent.click(screen.getByRole("button", { name: "記録する" }));

      await waitFor(() => {
        expect(post).toHaveBeenCalledTimes(1);
        expect(post).toHaveBeenCalledWith(
          "/api/v1/records",
          expect.objectContaining({
            items: [{ name: "白米", calories: 250 }],
          })
        );
      });
    });

    it("送信成功時にonSuccessが呼ばれる", async () => {
      const user = userEvent.setup();
      render(<RecordForm onSuccess={mockOnSuccess} />);

      // 食品名とカロリーを入力
      const foodNameInput = screen.getByLabelText("食品名");
      await user.type(foodNameInput, "白米");

      const calorieInput = screen.getByLabelText("カロリー (kcal)");
      await user.type(calorieInput, "250");

      // ボタンが有効化されることを確認
      await waitFor(() => {
        expect(screen.getByRole("button", { name: "記録する" })).not.toBeDisabled();
      });

      // 送信
      fireEvent.click(screen.getByRole("button", { name: "記録する" }));

      await waitFor(() => {
        expect(mockOnSuccess).toHaveBeenCalledTimes(1);
      });
    });

    it("バリデーションエラーがある場合は送信ボタンが無効化される", async () => {
      render(<RecordForm />);

      // 初期状態（空のフォーム）では送信ボタンは無効
      // validateOnMountが非同期で実行されるためwaitForを使用
      await waitFor(() => {
        const submitButton = screen.getByRole("button", { name: "記録する" });
        expect(submitButton).toBeDisabled();
      });
    });
  });

  describe("レンダリング", () => {
    it("タイトルとフィールドが正しく表示される", () => {
      render(<RecordForm />);

      expect(screen.getByLabelText("食事日時")).toBeInTheDocument();
      expect(screen.getByText("食品リスト")).toBeInTheDocument();
      expect(screen.getByLabelText("食品名")).toBeInTheDocument();
      expect(screen.getByLabelText("カロリー (kcal)")).toBeInTheDocument();
      expect(screen.getByRole("button", { name: "食品を追加" })).toBeInTheDocument();
      expect(screen.getByRole("button", { name: "記録する" })).toBeInTheDocument();
    });
  });
});
