class MaskDetector {
  static String detectMaskType(String text) {
    if (text.contains(RegExp(r'^\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}$'))) {
      return 'IP';
    }

    bool match = RegExp(r'^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$').hasMatch(text);
    if (match) {
      return 'Email';
    }

    return 'Unknown';
  }
}
