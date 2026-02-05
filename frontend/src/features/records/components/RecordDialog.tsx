/**
 * RecordDialog - カロリー記録ダイアログ
 * 「記録する」ボタンを押すとダイアログが開き、RecordFormを表示
 */
import { useState } from "react";
import {
  Dialog,
  DialogContent,
  DialogHeader,
  DialogTitle,
  DialogTrigger,
} from "@/components/ui/dialog";
import { Button } from "@/components/ui/button";
import { RecordForm } from "./RecordForm";

/** RecordDialogコンポーネントのProps */
export type RecordDialogProps = {
  /** 記録作成成功時のコールバック */
  onSuccess?: () => void;
};

/**
 * PlusIcon - 追加ボタン用アイコン
 */
function PlusIcon({ className }: { className?: string }) {
  return (
    <svg
      xmlns="http://www.w3.org/2000/svg"
      viewBox="0 0 24 24"
      fill="none"
      stroke="currentColor"
      strokeWidth="2"
      strokeLinecap="round"
      strokeLinejoin="round"
      className={className}
      aria-hidden="true"
    >
      <line x1="12" y1="5" x2="12" y2="19" />
      <line x1="5" y1="12" x2="19" y2="12" />
    </svg>
  );
}

/**
 * RecordDialog - カロリー記録ダイアログ
 */
export function RecordDialog({ onSuccess }: RecordDialogProps) {
  const [open, setOpen] = useState(false);

  const handleSuccess = () => {
    setOpen(false);
    onSuccess?.();
  };

  return (
    <Dialog open={open} onOpenChange={setOpen}>
      <DialogTrigger asChild>
        <Button className="gap-2">
          <PlusIcon className="w-4 h-4" />
          記録する
        </Button>
      </DialogTrigger>
      <DialogContent
        className="max-w-lg max-h-[90vh] overflow-y-auto"
        onOpenAutoFocus={(e) => e.preventDefault()}
      >
        <DialogHeader>
          <DialogTitle>カロリー記録</DialogTitle>
        </DialogHeader>
        <RecordForm onSuccess={handleSuccess} />
      </DialogContent>
    </Dialog>
  );
}
