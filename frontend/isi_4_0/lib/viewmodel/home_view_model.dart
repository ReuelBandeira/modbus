// ignore_for_file: public_member_api_docs, sort_constructors_first
import 'package:flutter/material.dart';
import 'package:isi_4_0/repository/db_api.dart';

class HomeViewModel extends ChangeNotifier {
  AuthService authService = AuthService.instance;

  Future<void> loginView(BuildContext context) async {
    Navigator.pushReplacementNamed(context, '/');
  }

  Future<int> updateLanguage(String id, String language) async {
    int value = await authService.updateProfile(id, "", "", "", language);
    notifyListeners();
    return value;
  }

  Future<void> updatePageViewController(
      PageController pageController, int position) async {
    notifyListeners();
  }

  bool ticTac() {
    bool timeout = false;
    Future.delayed(const Duration(seconds: 10)).then((value) {
      timeout = true;
    });
    return timeout;
  }
}
