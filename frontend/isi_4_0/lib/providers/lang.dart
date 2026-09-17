import 'package:flutter/material.dart';
import 'package:isi_4_0/repository/db_api.dart';
import 'package:isi_4_0/utils/language.dart';

class LanguageProvider extends ChangeNotifier {
  AuthService authService = AuthService();
  String _currentLanguage = "pt";
  String get currentLanguage => _currentLanguage;

  LanguageProvider() {
    initLanguage();
  }

  Future<void> initLanguage() async {
    _currentLanguage = await authService.getLang();
    notifyListeners();
  }

  void updateCurrentLanguage(String newLanguage) async {
    _currentLanguage = newLanguage;
    await authService.saveLang(newLanguage);
    notifyListeners();
  }

  Map<String, dynamic> getDataLanguage(String language) {
    Map<String, dynamic> translator = Language().getJsonData(language);
    return translator;
  }
}
