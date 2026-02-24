import { View, Text } from "react-native";
import { SafeAreaView } from "react-native-safe-area-context";

export default function HomeScreen() {
  return (
    <SafeAreaView className="flex-1 bg-background" edges={["bottom"]}>
      <View className="flex-1 items-center justify-center px-6">
        <Text className="text-3xl font-bold text-foreground mb-2">
          CalTrack
        </Text>
        <Text className="text-base text-muted-foreground text-center">
          カロリー管理アプリ
        </Text>
      </View>
    </SafeAreaView>
  );
}
