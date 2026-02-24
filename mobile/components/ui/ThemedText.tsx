import { Text, type TextProps } from "react-native";

type ThemedTextVariant = "default" | "title" | "subtitle" | "caption";

type ThemedTextProps = TextProps & {
  variant?: ThemedTextVariant;
};

const variantClasses: Record<ThemedTextVariant, string> = {
  default: "text-base text-foreground",
  title: "text-2xl font-bold text-foreground",
  subtitle: "text-lg font-semibold text-foreground",
  caption: "text-sm text-muted-foreground",
};

export function ThemedText({
  variant = "default",
  className = "",
  ...props
}: ThemedTextProps) {
  return (
    <Text className={`${variantClasses[variant]} ${className}`} {...props} />
  );
}
