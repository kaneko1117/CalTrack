---
name: frontend-design
description: Create distinctive, production-grade frontend interfaces with high design quality. Use when building UI components, designing layouts, or when asked to "make it look good", "improve the design", or "create a unique UI".
metadata:
  author: CalTrack
  version: "1.0.0"
---

# Frontend Design Guidelines

汎用的なAI生成デザインを避け、独特で記憶に残るUIを作成するためのガイドライン。

## 技術スタック

- React + TypeScript + Vite
- shadcn/ui
- Tailwind CSS

---

## 1. Design Thinking Process

UIを作成する前に、以下の4つの質問に答えること。

### 1.1 Purpose（目的）

このインターフェースが解決する問題は何か？

```
// 悪い例: 目的が曖昧
「ユーザーがデータを見れるようにする」

// 良い例: 明確な目的
「ユーザーが今日の摂取カロリーを一目で把握し、
 目標との差分を直感的に理解できるようにする」
```

### 1.2 Tone（トーン・美的方向性）

どのような美的方向性を目指すか？極端な方向性を選ぶこと。

| トーン | 特徴 | 使用例 |
|--------|------|--------|
| Brutally Minimal | 余白多め、要素最小限、タイポグラフィ重視 | ダッシュボード、設定画面 |
| Maximalist Chaos | 情報密度高、レイヤー多重、視覚的刺激 | ランディングページ、プロモーション |
| Retro-Futuristic | レトロ要素 + 未来的要素の融合 | ブランディング重視の画面 |
| Warm & Organic | 丸み、暖色、手書き風要素 | オンボーディング、ヘルス系 |
| Dark & Precise | ダークテーマ、シャープなエッジ、高コントラスト | データ可視化、プロ向け |

### 1.3 Constraints（技術的制約）

考慮すべき技術的要件をリストアップする。

- レスポンシブ対応の範囲（モバイルファースト？）
- アクセシビリティ要件（WCAG レベル）
- パフォーマンス制約（アニメーション可否）
- ブラウザサポート範囲

### 1.4 Differentiation（差別化要素）

何がこのUIを忘れられないものにするか？

```
// 悪い例: 差別化なし
「きれいなカードレイアウト」

// 良い例: 明確な差別化
「カロリーを"貯金"に見立てた金庫モチーフのプログレスバー、
 食事記録時の硬貨投入アニメーション」
```

---

## 2. Typography Guidelines

### 2.1 避けるべきフォント

以下のフォントは汎用的すぎるため使用を避ける。

```css
/* 禁止フォント */
font-family: Inter;        /* 使い古された */
font-family: Roboto;       /* 汎用的すぎる */
font-family: Arial;        /* デフォルト感 */
font-family: Helvetica;    /* 無個性 */
font-family: Open Sans;    /* 見慣れすぎ */
font-family: system-ui;    /* 差別化不可 */
```

### 2.2 推奨フォントの選び方

目的に応じたフォント選択指針。

| 用途 | 推奨カテゴリ | 例 |
|------|-------------|-----|
| 見出し | Display / Decorative | Playfair Display, Space Grotesk, Clash Display |
| 本文 | Readable Sans | DM Sans, Plus Jakarta Sans, Outfit |
| 数値 | Monospace / Tabular | JetBrains Mono, IBM Plex Mono, Fira Code |
| アクセント | Variable / Experimental | Instrument Sans, Satoshi, General Sans |

### 2.3 実装例

```tsx
// tailwind.config.ts でのフォント設定
import { fontFamily } from "tailwindcss/defaultTheme";

export default {
  theme: {
    extend: {
      fontFamily: {
        // 見出し用: 個性的なDisplay font
        display: ["Space Grotesk", ...fontFamily.sans],
        // 本文用: 読みやすいSans
        sans: ["DM Sans", ...fontFamily.sans],
        // 数値用: Tabular figures対応のMono
        mono: ["JetBrains Mono", ...fontFamily.mono],
      },
    },
  },
};
```

```tsx
// コンポーネントでの使用
<h1 className="font-display text-4xl font-bold tracking-tight">
  今日のカロリー
</h1>
<p className="font-sans text-base text-muted-foreground">
  目標達成まであと少し
</p>
<span className="font-mono text-2xl tabular-nums">
  1,234 kcal
</span>
```

### 2.4 タイポグラフィの階層

明確な視覚的階層を作ること。

```tsx
// 階層の例
const typographyScale = {
  // Display: 最重要の数値・見出し
  display: "text-5xl md:text-7xl font-display font-bold tracking-tighter",

  // Heading: セクション見出し
  h1: "text-3xl md:text-4xl font-display font-semibold tracking-tight",
  h2: "text-2xl md:text-3xl font-display font-semibold",
  h3: "text-xl md:text-2xl font-display font-medium",

  // Body: 本文
  body: "text-base font-sans leading-relaxed",
  small: "text-sm font-sans text-muted-foreground",

  // Data: 数値表示
  data: "font-mono tabular-nums",
};
```

---

## 3. Color & Theme Guidelines

### 3.1 避けるべきカラーパターン

```css
/* 禁止パターン */

/* 1. 紫グラデーション on 白背景（AI生成感が強い） */
background: linear-gradient(to right, #8b5cf6, #d946ef);

/* 2. 青から紫のグラデーション（使い古された） */
background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);

/* 3. パステル + 白の組み合わせ（無個性） */
background: #f8fafc;
color: #94a3b8;
```

### 3.2 カラーシステムの構築

CSS変数を使用した一貫したテーマを構築する。

```css
/* globals.css */
@layer base {
  :root {
    /*
     * CalTrack固有のカラーパレット
     * コンセプト: "健康的な活力" - 暖かみのあるグリーン基調
     */

    /* Primary: メインアクション、重要な要素 */
    --primary: 142 76% 36%;        /* 深みのあるグリーン */
    --primary-foreground: 0 0% 100%;

    /* Secondary: 補助的な要素 */
    --secondary: 45 93% 47%;        /* 活力のあるアンバー */
    --secondary-foreground: 0 0% 0%;

    /* Accent: 注目を引く要素 */
    --accent: 24 95% 53%;           /* 暖かいオレンジ */
    --accent-foreground: 0 0% 100%;

    /* Background: 背景色 */
    --background: 60 9% 98%;        /* オフホワイト（純白を避ける） */
    --foreground: 20 14% 4%;        /* ほぼ黒（純黒を避ける） */

    /* Muted: 控えめな要素 */
    --muted: 60 5% 96%;
    --muted-foreground: 25 5% 45%;

    /* Card: カード背景 */
    --card: 0 0% 100%;
    --card-foreground: 20 14% 4%;

    /* Border: 境界線 */
    --border: 20 6% 90%;

    /* Semantic: 意味のある色 */
    --success: 142 76% 36%;
    --warning: 45 93% 47%;
    --error: 0 84% 60%;
    --info: 199 89% 48%;
  }

  .dark {
    /* ダークテーマ: 純黒を避けた深みのある色 */
    --background: 20 14% 4%;
    --foreground: 60 9% 98%;
    --card: 24 10% 10%;
    --card-foreground: 60 9% 98%;
    --muted: 12 6% 15%;
    --muted-foreground: 24 5% 64%;
    --border: 12 6% 15%;
  }
}
```

### 3.3 カラーの使い方

```tsx
// 意図を持ったカラー使用
const CalorieDisplay = ({ current, goal }: Props) => {
  const percentage = (current / goal) * 100;

  // 状態に応じた色の変化
  const getStatusColor = () => {
    if (percentage < 50) return "text-success";      // まだ余裕
    if (percentage < 80) return "text-warning";      // 注意
    if (percentage < 100) return "text-accent";      // もうすぐ
    return "text-error";                              // 超過
  };

  return (
    <div className="relative">
      {/* グラデーションは控えめに、意図を持って使用 */}
      <div className="absolute inset-0 bg-gradient-to-b from-primary/5 to-transparent" />
      <span className={cn("font-mono text-5xl tabular-nums", getStatusColor())}>
        {current.toLocaleString()}
      </span>
    </div>
  );
};
```

---

## 4. Motion & Animation Guidelines

### 4.1 基本原則

- **目的のあるアニメーション**: 装飾ではなく、フィードバックや状態変化を伝える
- **CSS-firstアプローチ**: 可能な限りCSS transitionを使用
- **パフォーマンス**: transform と opacity のみアニメーション

### 4.2 CSS Transitionの活用

```css
/* globals.css */
@layer base {
  :root {
    /* 一貫したタイミング関数 */
    --ease-out-expo: cubic-bezier(0.16, 1, 0.3, 1);
    --ease-in-out-expo: cubic-bezier(0.87, 0, 0.13, 1);
    --ease-spring: cubic-bezier(0.34, 1.56, 0.64, 1);

    /* デュレーション */
    --duration-fast: 150ms;
    --duration-normal: 250ms;
    --duration-slow: 400ms;
  }
}
```

```tsx
// CSS Transitionの実装例
const Button = ({ children, ...props }: ButtonProps) => (
  <button
    className={cn(
      "relative overflow-hidden",
      "transition-all duration-[--duration-normal] ease-[--ease-out-expo]",
      // ホバー時: 軽い浮き上がり
      "hover:-translate-y-0.5 hover:shadow-lg",
      // アクティブ時: 押し込み
      "active:translate-y-0 active:shadow-sm",
      // フォーカス時: リング表示
      "focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary"
    )}
    {...props}
  >
    {children}
  </button>
);
```

### 4.3 Framer Motionの活用（必要な場合）

```tsx
import { motion, AnimatePresence } from "framer-motion";

// リスト項目のスタッガードアニメーション
const MealList = ({ meals }: { meals: Meal[] }) => (
  <motion.ul className="space-y-3">
    <AnimatePresence mode="popLayout">
      {meals.map((meal, index) => (
        <motion.li
          key={meal.id}
          layout
          initial={{ opacity: 0, y: 20 }}
          animate={{ opacity: 1, y: 0 }}
          exit={{ opacity: 0, x: -20 }}
          transition={{
            duration: 0.25,
            delay: index * 0.05,
            ease: [0.16, 1, 0.3, 1],
          }}
        >
          <MealCard meal={meal} />
        </motion.li>
      ))}
    </AnimatePresence>
  </motion.ul>
);

// 数値のカウントアップアニメーション
const AnimatedNumber = ({ value }: { value: number }) => {
  const springValue = useSpring(value, {
    stiffness: 100,
    damping: 30,
  });

  const display = useTransform(springValue, (v) =>
    Math.round(v).toLocaleString()
  );

  return (
    <motion.span className="font-mono tabular-nums">
      {display}
    </motion.span>
  );
};
```

### 4.4 マイクロインタラクション

```tsx
// 食事記録追加時のフィードバック
const AddMealButton = () => {
  const [isAdded, setIsAdded] = useState(false);

  return (
    <motion.button
      whileTap={{ scale: 0.95 }}
      onClick={() => {
        setIsAdded(true);
        // 成功フィードバック
        setTimeout(() => setIsAdded(false), 1500);
      }}
      className="relative"
    >
      <AnimatePresence mode="wait">
        {isAdded ? (
          <motion.span
            key="success"
            initial={{ opacity: 0, scale: 0.8 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.8 }}
            className="flex items-center gap-2 text-success"
          >
            <Check className="h-4 w-4" />
            追加しました
          </motion.span>
        ) : (
          <motion.span
            key="default"
            initial={{ opacity: 0, scale: 0.8 }}
            animate={{ opacity: 1, scale: 1 }}
            exit={{ opacity: 0, scale: 0.8 }}
            className="flex items-center gap-2"
          >
            <Plus className="h-4 w-4" />
            食事を追加
          </motion.span>
        )}
      </AnimatePresence>
    </motion.button>
  );
};
```

---

## 5. Layout & Spatial Composition

### 5.1 予期しないレイアウト

グリッドに縛られない、印象的なレイアウトを作る。

```tsx
// 非対称レイアウトの例
const HeroSection = () => (
  <section className="relative min-h-[80vh] overflow-hidden">
    {/* 背景: 非対称の装飾 */}
    <div className="absolute -top-1/4 -right-1/4 h-[600px] w-[600px] rounded-full bg-primary/10 blur-3xl" />
    <div className="absolute -bottom-1/4 -left-1/4 h-[400px] w-[400px] rounded-full bg-secondary/10 blur-3xl" />

    {/* コンテンツ: オフセットされた配置 */}
    <div className="container relative grid min-h-[80vh] grid-cols-12 items-center gap-8">
      {/* 左側: テキスト（5列分、左に寄せる） */}
      <div className="col-span-12 md:col-span-5 md:col-start-1">
        <h1 className="font-display text-5xl font-bold tracking-tighter md:text-7xl">
          食べる。
          <br />
          <span className="text-primary">記録する。</span>
          <br />
          健康になる。
        </h1>
      </div>

      {/* 右側: ビジュアル（6列分、重なりを作る） */}
      <div className="col-span-12 md:col-span-6 md:col-start-6 md:-mt-20">
        <CalorieVisualization />
      </div>
    </div>
  </section>
);
```

### 5.2 オーバーラップとレイヤー

```tsx
// カードのオーバーラップ配置
const StatsOverview = () => (
  <div className="relative h-[400px]">
    {/* ベースカード */}
    <div className="absolute left-0 top-0 w-3/4 rounded-2xl bg-card p-6 shadow-lg">
      <h3 className="font-display text-lg font-semibold">今週の記録</h3>
      <WeeklyChart />
    </div>

    {/* オーバーラップするカード */}
    <div className="absolute right-0 top-20 w-1/2 rounded-2xl bg-primary p-6 text-primary-foreground shadow-xl">
      <span className="text-sm opacity-80">本日の達成率</span>
      <p className="font-mono text-4xl font-bold tabular-nums">87%</p>
    </div>

    {/* さらに上にオーバーラップ */}
    <div className="absolute bottom-0 left-1/4 rounded-full bg-secondary px-4 py-2 shadow-lg">
      <span className="text-sm font-medium">3日連続達成中！</span>
    </div>
  </div>
);
```

### 5.3 余白の戦略的使用

```tsx
// 余白で呼吸感を作る
const SectionLayout = ({ children, title }: LayoutProps) => (
  <section className="py-20 md:py-32">
    {/* 見出し: 大きな余白で存在感を出す */}
    <header className="mb-16 md:mb-24">
      <h2 className="font-display text-3xl font-bold tracking-tight md:text-5xl">
        {title}
      </h2>
    </header>

    {/* コンテンツ: 適度な余白 */}
    <div className="space-y-8">
      {children}
    </div>
  </section>
);
```

---

## 6. Backgrounds & Visual Details

### 6.1 雰囲気を作る背景

```tsx
// グレイン（ノイズ）テクスチャ
const GrainOverlay = () => (
  <div
    className="pointer-events-none fixed inset-0 z-50 opacity-[0.015]"
    style={{
      backgroundImage: `url("data:image/svg+xml,%3Csvg viewBox='0 0 256 256' xmlns='http://www.w3.org/2000/svg'%3E%3Cfilter id='noise'%3E%3CfeTurbulence type='fractalNoise' baseFrequency='0.9' numOctaves='4' stitchTiles='stitch'/%3E%3C/filter%3E%3Crect width='100%25' height='100%25' filter='url(%23noise)'/%3E%3C/svg%3E")`,
    }}
  />
);

// グラデーションメッシュ背景
const MeshGradient = () => (
  <div className="fixed inset-0 -z-10">
    <div className="absolute inset-0 bg-background" />
    <div className="absolute left-1/4 top-0 h-[500px] w-[500px] rounded-full bg-primary/20 blur-[100px]" />
    <div className="absolute right-1/4 top-1/3 h-[400px] w-[400px] rounded-full bg-secondary/20 blur-[100px]" />
    <div className="absolute bottom-0 left-1/2 h-[600px] w-[600px] -translate-x-1/2 rounded-full bg-accent/10 blur-[120px]" />
  </div>
);
```

### 6.2 装飾的なディテール

```tsx
// ドット背景パターン
const DotPattern = () => (
  <div
    className="absolute inset-0 -z-10 opacity-[0.4]"
    style={{
      backgroundImage: `radial-gradient(circle, hsl(var(--foreground)) 1px, transparent 1px)`,
      backgroundSize: '24px 24px',
    }}
  />
);

// グリッド背景パターン
const GridPattern = () => (
  <div
    className="absolute inset-0 -z-10"
    style={{
      backgroundImage: `
        linear-gradient(to right, hsl(var(--border)) 1px, transparent 1px),
        linear-gradient(to bottom, hsl(var(--border)) 1px, transparent 1px)
      `,
      backgroundSize: '40px 40px',
    }}
  />
);
```

### 6.3 カードとコンテナのスタイリング

```tsx
// 深みのあるカードスタイル
const Card = ({ children, variant = "default" }: CardProps) => {
  const variants = {
    default: cn(
      "rounded-2xl bg-card p-6",
      "border border-border/50",
      "shadow-sm shadow-foreground/5"
    ),
    elevated: cn(
      "rounded-2xl bg-card p-6",
      "shadow-xl shadow-foreground/10",
      "ring-1 ring-border/10"
    ),
    glass: cn(
      "rounded-2xl p-6",
      "bg-card/80 backdrop-blur-xl",
      "border border-border/30"
    ),
  };

  return <div className={variants[variant]}>{children}</div>;
};
```

---

## 7. Anti-Patterns（避けるべきパターン）

### 7.1 Typography Anti-Patterns

```tsx
// 悪い例
<h1 className="font-sans text-2xl">見出し</h1>  // 階層が不明確
<p className="text-gray-400">本文</p>            // コントラスト不足

// 良い例
<h1 className="font-display text-4xl font-bold tracking-tight">見出し</h1>
<p className="text-muted-foreground">本文</p>
```

### 7.2 Color Anti-Patterns

```tsx
// 悪い例: AI生成感の強いグラデーション
<div className="bg-gradient-to-r from-purple-500 to-pink-500" />

// 悪い例: 意味のない色使い
<span className="text-blue-500">エラー</span>  // 青はエラーを示さない

// 良い例: 意図を持った色使い
<span className="text-error">エラーが発生しました</span>
```

### 7.3 Layout Anti-Patterns

```tsx
// 悪い例: すべてが中央寄せ
<div className="flex flex-col items-center justify-center text-center">
  <h1>見出し</h1>
  <p>本文</p>
  <button>ボタン</button>
</div>

// 良い例: 意図を持った配置
<div className="flex flex-col items-start">
  <h1 className="max-w-2xl">見出し</h1>
  <p className="mt-4 max-w-lg text-muted-foreground">本文</p>
  <button className="mt-8">ボタン</button>
</div>
```

### 7.4 Animation Anti-Patterns

```tsx
// 悪い例: 過剰なアニメーション
<motion.div
  animate={{ rotate: 360, scale: [1, 1.5, 1] }}
  transition={{ duration: 2, repeat: Infinity }}
/>

// 悪い例: 目的のないアニメーション
<div className="animate-bounce" />  // なぜバウンス？

// 良い例: 目的のあるアニメーション
<motion.div
  initial={{ opacity: 0, y: 10 }}
  animate={{ opacity: 1, y: 0 }}
  transition={{ duration: 0.25 }}
/>
```

---

## 8. Checklist

UIを作成する際は以下をチェックすること。

### Design Thinking
- [ ] Purpose（目的）が明確か
- [ ] Tone（美的方向性）が決まっているか
- [ ] Constraints（制約）を把握しているか
- [ ] Differentiation（差別化要素）があるか

### Typography
- [ ] 汎用フォントを避けているか
- [ ] 視覚的階層が明確か
- [ ] 数値にはtabular-numsを使用しているか

### Color
- [ ] AI生成感のあるグラデーションを避けているか
- [ ] CSS変数で一貫したテーマを使用しているか
- [ ] 純白・純黒を避けているか

### Motion
- [ ] アニメーションに目的があるか
- [ ] CSS transitionを優先しているか
- [ ] パフォーマンスを考慮しているか（transform, opacityのみ）

### Layout
- [ ] 予測可能すぎるレイアウトを避けているか
- [ ] 余白を戦略的に使用しているか
- [ ] オーバーラップやレイヤーを活用しているか

### Visual Details
- [ ] 背景に深みがあるか
- [ ] 装飾的なディテールが適切か
- [ ] カードスタイルに個性があるか
