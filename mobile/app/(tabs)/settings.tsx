import { View, Text } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

export default function SettingsScreen() {
  return (
    <SafeAreaView className="flex-1 bg-background" edges={["bottom"]}>
      <View className="flex-1 items-center justify-center px-6">
        <Text className="text-xl font-bold text-foreground mb-2">設定</Text>
        <Text className="text-base text-muted-foreground text-center">
          設定画面（準備中）
        </Text>
      </View>
    </SafeAreaView>
  );
}
