import 'package:isi_4_0/utils/languages/app_en.dart';
import 'package:isi_4_0/utils/languages/app_es.dart';
import 'package:isi_4_0/utils/languages/app_pt.dart';

class Language {
  // Singleton
  static final Language instance = Language._internal();
  factory Language() {
    return instance;
  }
  Language._internal();

  Map<String, dynamic> jsonData = {};

  Map<String, dynamic> getJsonData(String lang) {
    switch (lang) {
      case 'es':
        jsonData = translateJsonES;
        break;
      case 'pt':
        jsonData = translateJsonPT;
        break;
      case 'en':
      default:
        jsonData = translateJsonEN;
        break;
    }

    return jsonData;
  }
}
