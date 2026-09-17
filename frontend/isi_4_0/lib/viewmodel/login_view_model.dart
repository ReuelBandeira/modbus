import 'package:flutter/material.dart';
import 'package:isi_4_0/repository/db_api.dart';
import 'package:isi_4_0/utils/enum_login.dart';

class LoginViewModel extends ChangeNotifier {
  CurrentPage _currentPage = CurrentPage.signIn;
  CurrentPage get currentPage => _currentPage;
  AuthService authService = AuthService.instance;

  Future<void> updateCurrentPage(CurrentPage page) async {
    _currentPage = page;
    notifyListeners();
  }

  Future<void> signInView(BuildContext context) async {
    Navigator.pushReplacementNamed(context, '/');
  }

  Future<void> homeView(BuildContext context) async {
    Navigator.pushReplacementNamed(context, '/home');
  }

  Future<void> recoverView(BuildContext context) async {
    Navigator.pushReplacementNamed(context, '/recover');
  }

  Future<void> registerView(BuildContext context) async {
    Navigator.pushReplacementNamed(context, '/register');
  }
}
